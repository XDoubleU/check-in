package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestGetSortedSchoolsOK(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	defaultLocation := *fixtureData.DefaultLocation

	for i := 0; i < 10; i++ {
		_, _ = testApp.services.CheckIns.Create(
			context.Background(),
			&defaultLocation,
			&models.School{ID: 1},
		)
	}

	for i := 0; i < 10; i++ {
		_, _ = testApp.services.CheckIns.Create(
			context.Background(),
			&defaultLocation,
			fixtureData.Schools[0],
		)
	}

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/checkins/schools")
	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData []models.School
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rsData[0].ID, fixtureData.Schools[0].ID)
	assert.Equal(t, rsData[len(rsData)-1].ID, int64(1))
}

func TestGetSortedSchoolsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq1 := test.CreateTestRequest(t, ts, http.MethodGet, "/checkins/schools")

	tReq2 := test.CreateTestRequest(t, ts, http.MethodGet, "/checkins/schools")
	tReq2.AddCookie(tokens.ManagerAccessToken)

	tReq3 := test.CreateTestRequest(t, ts, http.MethodGet, "/checkins/schools")
	tReq3.AddCookie(tokens.AdminAccessToken)

	rs1 := tReq1.Do(nil)
	rs2 := tReq2.Do(nil)
	rs3 := tReq3.Do(nil)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}

func TestCreateCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: fixtureData.Schools[0].ID,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData dtos.CheckInDto
	rs := tReq.Do(&rsData)

	loc, _ := time.LoadLocation("Europe/Brussels")

	assert.Equal(t, rs.StatusCode, http.StatusCreated)
	assert.Equal(t, rsData.SchoolName, fixtureData.Schools[0].Name)
	assert.Equal(t, rsData.LocationID, fixtureData.DefaultLocation.ID)
	assert.Equal(t, rsData.Capacity, fixtureData.DefaultLocation.Capacity)
	assert.Equal(
		t,
		rsData.CreatedAt.Time.Format(constants.DateFormat),
		time.Now().In(loc).Format(constants.DateFormat),
	)
}

func TestCreateCheckInAndere(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 1,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData dtos.CheckInDto
	rs := tReq.Do(&rsData)

	loc, _ := time.LoadLocation("Europe/Brussels")

	assert.Equal(t, rs.StatusCode, http.StatusCreated)
	assert.Equal(t, rsData.SchoolName, "Andere")
	assert.Equal(t, rsData.LocationID, fixtureData.DefaultLocation.ID)
	assert.Equal(t, rsData.Capacity, fixtureData.DefaultLocation.Capacity)
	assert.Equal(
		t,
		rsData.CreatedAt.Time.Format(constants.DateFormat),
		time.Now().In(loc).Format(constants.DateFormat),
	)
}

func TestCreateCheckInAboveCap(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 1,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(tokens.DefaultAccessToken)

	var rs *http.Response
	var rsData http_tools.ErrorDto

	for i := 0; i < int(fixtureData.DefaultLocation.Capacity)+1; i++ {
		rs = tReq.Do(&rsData)
	}

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message, "location has no available spots")
}

func TestCreateCheckInSchoolNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 8000,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["schoolId"].(string),
		fmt.Sprintf("school with id '%d' doesn't exist", data.SchoolID),
	)
}

func TestCreateCheckInFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/checkins")

	data := dtos.CreateCheckInDto{
		SchoolID: 0,
	}

	tReq.SetReqData(data)

	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["schoolId"],
		"must be greater than zero",
	)
}

func TestCreateCheckInAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq1 := test.CreateTestRequest(t, ts, http.MethodPost, "/checkins")

	tReq2 := test.CreateTestRequest(t, ts, http.MethodPost, "/checkins")
	tReq2.AddCookie(tokens.ManagerAccessToken)

	tReq3 := test.CreateTestRequest(t, ts, http.MethodPost, "/checkins")
	tReq3.AddCookie(tokens.AdminAccessToken)

	rs1 := tReq1.Do(nil)
	rs2 := tReq2.Do(nil)
	rs3 := tReq3.Do(nil)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}
