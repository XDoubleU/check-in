package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"check-in/api/internal/assert"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"
)

func TestGetPaginatedSchoolsDefaultPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/schools", nil)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData dtos.PaginatedSchoolsDto
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		assert.Equal(t, rsData.Pagination.Current, 1)
		assert.Equal(
			t,
			rsData.Pagination.Total,
			int64(math.Ceil(float64(fixtureData.AmountOfLocations)/4)),
		)
		assert.Equal(t, len(rsData.Data), 4)

		assert.Equal(t, rsData.Data[0].ID, 1)
		assert.Equal(t, rsData.Data[0].Name, "Andere")
		assert.Equal(t, rsData.Data[0].ReadOnly, true)
	}
}

func TestGetPaginatedSchoolsSpecificPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/schools?page=2", nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.PaginatedSchoolsDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.Equal(t, rsData.Pagination.Current, 2)
	assert.Equal(
		t,
		rsData.Pagination.Total,
		int64(math.Ceil(float64(fixtureData.AmountOfLocations)/4)),
	)
	assert.Equal(t, len(rsData.Data), 4)

	assert.Equal(t, rsData.Data[0].ID, fixtureData.Schools[11].ID)
	assert.Equal(t, rsData.Data[0].Name, fixtureData.Schools[11].Name)
	assert.Equal(t, rsData.Data[0].ReadOnly, fixtureData.Schools[11].ReadOnly)
}

func TestGetPaginatedSchoolsPageZero(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/schools?page=0", nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message, "invalid page query param")
}

func TestGetPaginatedSchoolsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodGet, ts.URL+"/schools", nil)

	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/schools", nil)
	req2.AddCookie(tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestCreateSchool(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		unique := fmt.Sprintf("test%d", i)

		data := dtos.SchoolDto{
			Name: unique,
		}

		body, _ := json.Marshal(data)
		req, _ := http.NewRequest(
			http.MethodPost,
			ts.URL+"/schools",
			bytes.NewReader(body),
		)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData models.School
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusCreated)
		assert.Equal(t, rsData.Name, unique)
		assert.Equal(t, rsData.ReadOnly, false)
	}
}

func TestCreateSchoolNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "TestSchool0",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/schools", bytes.NewReader(body))
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("school with name '%s' already exists", data.Name),
	)
}

func TestCreateSchoolFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/schools", bytes.NewReader(body))
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["name"],
		"must be provided",
	)
}

func TestCreateSchoolAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodPost, ts.URL+"/schools", nil)

	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/schools", nil)
	req2.AddCookie(tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestUpdateSchool(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		unique := fmt.Sprintf("test%d", i)

		data := dtos.SchoolDto{
			Name: unique,
		}

		body, _ := json.Marshal(data)
		req, _ := http.NewRequest(
			http.MethodPatch,
			ts.URL+"/schools/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10),
			bytes.NewReader(body),
		)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData models.School
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.Schools[0].ID)
		assert.Equal(t, rsData.Name, unique)
		assert.Equal(t, rsData.ReadOnly, false)
	}
}

func TestUpdateSchoolNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "TestSchool1",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/schools/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10),
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("school with name '%s' already exists", data.Name),
	)
}

func TestUpdateSchoolReadOnly(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "test",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/schools/1",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(t, rsData.Message, "school with id '1' doesn't exist")
}

func TestUpdateSchoolNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "test",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/schools/8000",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(t, rsData.Message, "school with id '8000' doesn't exist")
}

func TestUpdateSchoolNotInt(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "test",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/schools/aaaa",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message, "invalid id parameter")
}

func TestUpdateSchoolFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/schools/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10),
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["name"],
		"must be provided",
	)
}

func TestUpdateSchoolAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/schools/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10),
		nil,
	)

	req2, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/schools/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10),
		nil,
	)
	req2.AddCookie(tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestDeleteSchool(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		req, _ := http.NewRequest(
			http.MethodDelete,
			ts.URL+"/schools/"+strconv.FormatInt(fixtureData.Schools[i].ID, 10),
			nil,
		)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData models.School
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.Schools[i].ID)
		assert.Equal(t, rsData.Name, fixtureData.Schools[i].Name)
		assert.Equal(t, rsData.ReadOnly, false)
	}
}

func TestDeleteSchoolReadOnly(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/schools/1", nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(t, rsData.Message.(string), "school with id '1' doesn't exist")
}

func TestDeleteSchoolNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/schools/8000", nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(t, rsData.Message.(string), "school with id '8000' doesn't exist")
}

func TestDeleteSchoolNotInt(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/schools/aaaa", nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "invalid id parameter")
}

func TestDeleteSchoolAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/schools/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10),
		nil,
	)

	req2, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/schools/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10),
		nil,
	)
	req2.AddCookie(tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}
