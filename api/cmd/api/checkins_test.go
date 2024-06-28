package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestGetSortedSchoolsOK(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	defaultLocation := *fixtureData.DefaultLocation

	for i := 0; i < 10; i++ {
		_, _ = testApp.repositories.CheckIns.Create(
			context.Background(),
			&defaultLocation,
			&models.School{ID: 1}, // Should always stay at the bottom
		)
	}

	for i := 0; i < 10; i++ {
		_, _ = testApp.repositories.CheckIns.Create(
			context.Background(),
			&defaultLocation,
			fixtureData.Schools[0],
		)
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/checkins/schools",
	)
	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData []models.School
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, fixtureData.Schools[0].ID, rsData[0].ID)
	assert.EqualValues(t, 1, rsData[len(rsData)-1].ID)
}

func TestGetSortedSchoolsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/checkins/schools",
	)

	mt := test.CreateMatrixTester(tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.ManagerAccessToken, http.StatusForbidden)
	mt.AddTestCaseCookieStatusCode(tokens.AdminAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestCreateCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: fixtureData.Schools[0].ID,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData dtos.CheckInDto
	rs := tReq.Do(t, &rsData)

	loc, _ := time.LoadLocation("Europe/Brussels")

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, fixtureData.Schools[0].Name, rsData.SchoolName)
	assert.Equal(t, fixtureData.DefaultLocation.ID, rsData.LocationID)
	assert.Equal(t, fixtureData.DefaultLocation.Capacity, rsData.Capacity)
	assert.Equal(
		t,
		time.Now().In(loc).Format(constants.DateFormat),
		rsData.CreatedAt.Time.Format(constants.DateFormat),
	)
}

func TestCreateCheckInAndere(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 1,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData dtos.CheckInDto
	rs := tReq.Do(t, &rsData)

	loc, _ := time.LoadLocation("Europe/Brussels")

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, "Andere", rsData.SchoolName)
	assert.Equal(t, fixtureData.DefaultLocation.ID, rsData.LocationID)
	assert.Equal(t, fixtureData.DefaultLocation.Capacity, rsData.Capacity)
	assert.Equal(
		t,
		time.Now().In(loc).Format(constants.DateFormat),
		rsData.CreatedAt.Time.Format(constants.DateFormat),
	)
}

func TestCreateCheckInAboveCap(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 1,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(tokens.DefaultAccessToken)

	var rs *http.Response
	var rsData httptools.ErrorDto

	for i := 0; i < int(fixtureData.DefaultLocation.Capacity)+1; i++ {
		rs = tReq.Do(t, &rsData)
	}

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "location has no available spots", rsData.Message)
}

func TestCreateCheckInSchoolNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 8000,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("school with schoolId '%d' doesn't exist", data.SchoolID),
		rsData.Message.(map[string]interface{})["schoolId"].(string),
	)
}

func TestCreateCheckInFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")
	tReq.AddCookie(tokens.DefaultAccessToken)

	mt := test.CreateMatrixTester(tReq)

	reqData := dtos.CreateCheckInDto{
		SchoolID: 0,
	}

	mt.AddTestCaseErrorMessage(reqData, map[string]interface{}{
		"schoolId": "must be greater than 0",
	})

	mt.Do(t)
}

func TestCreateCheckInAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/checkins")

	mt := test.CreateMatrixTester(tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.ManagerAccessToken, http.StatusForbidden)
	mt.AddTestCaseCookieStatusCode(tokens.AdminAccessToken, http.StatusForbidden)

	mt.Do(t)
}
