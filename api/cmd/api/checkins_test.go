package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"check-in/api/internal/assert"
	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"
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

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/checkins/schools", nil)
	req.AddCookie(tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData []models.School
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rsData[0].ID, fixtureData.Schools[0].ID)
	assert.Equal(t, rsData[len(rsData)-1].ID, 1)
}

func TestGetSortedSchoolsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodGet, ts.URL+"/checkins/schools", nil)

	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/checkins/schools", nil)
	req2.AddCookie(tokens.ManagerAccessToken)

	req3, _ := http.NewRequest(http.MethodGet, ts.URL+"/checkins/schools", nil)
	req3.AddCookie(tokens.AdminAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)
	rs3, _ := ts.Client().Do(req3)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}

func TestCreateCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CheckInDto{
		SchoolID: fixtureData.Schools[0].ID,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/checkins",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData models.CheckIn
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusCreated)
	assert.Equal(t, rsData.SchoolID, data.SchoolID)
	assert.Equal(t, rsData.LocationID, fixtureData.DefaultLocation.ID)
	assert.Equal(t, rsData.Capacity, fixtureData.DefaultLocation.Capacity)
	assert.Equal(
		t,
		rsData.CreatedAt.Format(constants.DateFormat),
		time.Now().Format(constants.DateFormat),
	)
}

func TestCreateCheckInAndere(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CheckInDto{
		SchoolID: 1,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/checkins",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData models.CheckIn
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusCreated)
	assert.Equal(t, rsData.SchoolID, data.SchoolID)
	assert.Equal(t, rsData.LocationID, fixtureData.DefaultLocation.ID)
	assert.Equal(t, rsData.Capacity, fixtureData.DefaultLocation.Capacity)
	assert.Equal(
		t,
		rsData.CreatedAt.Format(constants.DateFormat),
		time.Now().Format(constants.DateFormat),
	)
}

func TestCreateCheckInAboveCap(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CheckInDto{
		SchoolID: 1,
	}

	body, _ := json.Marshal(data)

	var rs *http.Response
	for i := 0; i < int(fixtureData.DefaultLocation.Capacity)+1; i++ {
		req, _ := http.NewRequest(
			http.MethodPost,
			ts.URL+"/checkins",
			bytes.NewReader(body),
		)
		req.AddCookie(tokens.DefaultAccessToken)

		rs, _ = ts.Client().Do(req)
	}

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message, "location has no available spots")
}

func TestCreateCheckInSchoolNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CheckInDto{
		SchoolID: 8000,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/checkins",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("school with id '%d' doesn't exist", data.SchoolID),
	)
}

func TestCreateCheckInFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CheckInDto{
		SchoolID: 0,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/checkins",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

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

	req1, _ := http.NewRequest(http.MethodPost, ts.URL+"/checkins", nil)

	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/checkins", nil)
	req2.AddCookie(tokens.ManagerAccessToken)

	req3, _ := http.NewRequest(http.MethodPost, ts.URL+"/checkins", nil)
	req3.AddCookie(tokens.AdminAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)
	rs3, _ := ts.Client().Do(req3)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}
