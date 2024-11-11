package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	errortools "github.com/XDoubleU/essentia/pkg/errors"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestGetSortedSchoolsOK(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	// Should always stay at the bottom
	testEnv.createCheckIns(fixtures.DefaultLocation, 1, 10)

	school := testEnv.createSchools(1)[0]
	testEnv.createCheckIns(fixtures.DefaultLocation, school.ID, 10)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/checkins/schools",
	)
	tReq.AddCookie(fixtures.Tokens.DefaultAccessToken)

	rs := tReq.Do(t)

	var rsData []models.School
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

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

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized, nil, nil))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(fixtures.Tokens.ManagerAccessToken)
	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(fixtures.Tokens.AdminAccessToken)
	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	mt.Do(t)
}

func TestCreateCheckIn(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	school := testEnv.createSchools(1)[0]

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")
	//nolint:exhaustruct //other fields are optional
	tReq.SetBody(dtos.CreateCheckInDto{
		SchoolID: school.ID,
	})
	tReq.AddCookie(fixtures.Tokens.DefaultAccessToken)

	rs := tReq.Do(t)

	var rsData dtos.CheckInDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	loc, _ := time.LoadLocation("Europe/Brussels")

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, school.Name, rsData.SchoolName)
	assert.Equal(t, fixtures.DefaultLocation.ID, rsData.LocationID)
	assert.Equal(t, fixtures.DefaultLocation.Capacity, rsData.Capacity)
	assert.Equal(
		t,
		testApp.getTimeNow().In(loc).Format(constants.DateFormat),
		rsData.CreatedAt.Time.Format(constants.DateFormat),
	)
}

func TestCreateCheckInAndere(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	//nolint:exhaustruct //other fields are optional
	data := dtos.CreateCheckInDto{
		SchoolID: 1,
	}

	tReq.SetBody(data)

	tReq.AddCookie(fixtures.Tokens.DefaultAccessToken)

	rs := tReq.Do(t)

	var rsData dtos.CheckInDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	loc, _ := time.LoadLocation("Europe/Brussels")

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, "Andere", rsData.SchoolName)
	assert.Equal(t, fixtures.DefaultLocation.ID, rsData.LocationID)
	assert.Equal(t, fixtures.DefaultLocation.Capacity, rsData.Capacity)
	assert.Equal(
		t,
		testApp.getTimeNow().In(loc).Format(constants.DateFormat),
		rsData.CreatedAt.Time.Format(constants.DateFormat),
	)
}

func TestCreateCheckInAboveCap(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	//nolint:exhaustruct //other fields are optional
	data := dtos.CreateCheckInDto{
		SchoolID: 1,
	}

	tReq.SetBody(data)

	tReq.AddCookie(fixtures.Tokens.DefaultAccessToken)

	var rs *http.Response

	for i := 0; i <= int(fixtures.DefaultLocation.Capacity); i++ {
		rs = tReq.Do(t)
	}

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "location has no available spots", rsData.Message)
}

func TestCreateCheckInSchoolNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	//nolint:exhaustruct //other fields are optional
	data := dtos.CreateCheckInDto{
		SchoolID: 8000,
	}

	tReq.SetBody(data)

	tReq.AddCookie(fixtures.Tokens.DefaultAccessToken)

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

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
	tReq.AddCookie(fixtures.Tokens.DefaultAccessToken)
	//nolint:exhaustruct //other fields are optional
	tReq.SetBody(dtos.CreateCheckInDto{
		SchoolID: 0,
	})

	mt := test.CreateMatrixTester()

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity, nil,
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"schoolId": "must be greater than 0",
		}))

	mt.AddTestCase(tReq, tRes)

	mt.Do(t)
}

func TestCreateCheckInAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized, nil, nil))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(fixtures.Tokens.ManagerAccessToken)
	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(fixtures.Tokens.AdminAccessToken)
	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	mt.Do(t)
}
