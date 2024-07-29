package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/test"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestGetSortedSchoolsOK(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	// Should always stay at the bottom
	testEnv.createCheckIns(testEnv.Fixtures.DefaultLocation, 1, 10)

	school := testEnv.createSchools(1)[0]
	testEnv.createCheckIns(testEnv.Fixtures.DefaultLocation, school.ID, 10)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/checkins/schools",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	var rsData []models.School
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, school.ID, rsData[0].ID)
	assert.EqualValues(t, 1, rsData[len(rsData)-1].ID)
}

func TestGetSortedSchoolsAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/checkins/schools",
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)
	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(testEnv.Fixtures.Tokens.AdminAccessToken)
	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestCreateCheckIn(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	school := testEnv.createSchools(1)[0]

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")
	tReq.SetReqData(dtos.CreateCheckInDto{
		SchoolID: school.ID,
	})
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	var rsData dtos.CheckInDto
	rs := tReq.Do(t, &rsData)

	loc, _ := time.LoadLocation("Europe/Brussels")

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, school.Name, rsData.SchoolName)
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.ID, rsData.LocationID)
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.Capacity, rsData.Capacity)
	assert.Equal(
		t,
		time.Now().In(loc).Format(constants.DateFormat),
		rsData.CreatedAt.Time.Format(constants.DateFormat),
	)
}

func TestCreateCheckInAndere(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 1,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	var rsData dtos.CheckInDto
	rs := tReq.Do(t, &rsData)

	loc, _ := time.LoadLocation("Europe/Brussels")

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, "Andere", rsData.SchoolName)
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.ID, rsData.LocationID)
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.Capacity, rsData.Capacity)
	assert.Equal(
		t,
		time.Now().In(loc).Format(constants.DateFormat),
		rsData.CreatedAt.Time.Format(constants.DateFormat),
	)
}

func TestCreateCheckInAboveCap(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 1,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	var rs *http.Response
	var rsData errortools.ErrorDto

	for i := 0; i < int(testEnv.Fixtures.DefaultLocation.Capacity)+1; i++ {
		rs = tReq.Do(t, &rsData)
	}

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "location has no available spots", rsData.Message)
}

func TestCreateCheckInSchoolNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 8000,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("school with schoolId '%d' doesn't exist", data.SchoolID),
		rsData.Message.(map[string]interface{})["schoolId"].(string),
	)
}

func TestCreateCheckInFailValidation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)
	tReq.SetReqData(dtos.CreateCheckInDto{
		SchoolID: 0,
	})

	mt := test.CreateMatrixTester()

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes.SetExpectedBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"schoolId": "must be greater than 0",
		}),
	)

	mt.AddTestCase(tReq, tRes)

	mt.Do(t)
}

func TestCreateCheckInAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)
	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(testEnv.Fixtures.Tokens.AdminAccessToken)
	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}
