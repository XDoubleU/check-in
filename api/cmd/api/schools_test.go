package main

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestGetPaginatedSchoolsDefaultPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/schools")
		tReq.AddCookie(user)

		var rsData dtos.PaginatedSchoolsDto
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		assert.EqualValues(t, 1, rsData.Pagination.Current)
		assert.EqualValues(
			t,
			math.Ceil(float64(fixtureData.AmountOfSchools)/4),
			rsData.Pagination.Total,
		)
		assert.Equal(t, 4, len(rsData.Data))

		assert.EqualValues(t, 1, rsData.Data[0].ID)
		assert.Equal(t, "Andere", rsData.Data[0].Name)
		assert.Equal(t, true, rsData.Data[0].ReadOnly)
	}
}

func TestGetPaginatedSchoolsSpecificPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/schools")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	var rsData dtos.PaginatedSchoolsDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 2, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(fixtureData.AmountOfSchools)/4),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 4, len(rsData.Data))

	assert.Equal(t, fixtureData.Schools[11].ID, rsData.Data[0].ID)
	assert.Equal(t, fixtureData.Schools[11].Name, rsData.Data[0].Name)
	assert.Equal(t, fixtureData.Schools[11].ReadOnly, rsData.Data[0].ReadOnly)
}

func TestGetPaginatedSchoolsFull(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/schools")
	tReq.AddCookie(tokens.ManagerAccessToken)

	test.PaginatedEndpointTester(
		t,
		tReq,
		"page",
		int(math.Ceil(float64(fixtureData.AmountOfSchools)/4)),
	)
}

func TestGetPaginatedSchoolsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/schools")

	mt := test.CreateMatrixTester(tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestCreateSchool(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		unique := fmt.Sprintf("test%d", i)

		data := dtos.SchoolDto{
			Name: unique,
		}

		tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/schools")
		tReq.AddCookie(user)

		tReq.SetReqData(data)

		var rsData models.School
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusCreated, rs.StatusCode)
		assert.Equal(t, unique, rsData.Name)
		assert.Equal(t, false, rsData.ReadOnly)
	}
}

func TestCreateSchoolNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	data := dtos.SchoolDto{
		Name: "TestSchool0",
	}

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/schools")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("school with name '%s' already exists", data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestCreateSchoolFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/schools")
	tReq.AddCookie(tokens.ManagerAccessToken)

	data := dtos.SchoolDto{
		Name: "",
	}

	mt := test.CreateMatrixTester(tReq)
	mt.AddTestCaseErrorMessage(data, map[string]interface{}{
		"name": "must be provided",
	})

	mt.Do(t)
}

func TestCreateSchoolAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/schools")

	mt := test.CreateMatrixTester(tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestUpdateSchool(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		unique := fmt.Sprintf("test%d", i)

		data := dtos.SchoolDto{
			Name: unique,
		}

		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodPatch,
			"/schools/%d",
			fixtureData.Schools[0].ID,
		)
		tReq.AddCookie(user)

		tReq.SetReqData(data)

		var rsData models.School
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, fixtureData.Schools[0].ID, rsData.ID)
		assert.Equal(t, unique, rsData.Name)
		assert.Equal(t, false, rsData.ReadOnly)
	}
}

func TestUpdateSchoolNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	data := dtos.SchoolDto{
		Name: "TestSchool1",
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/%d",
		fixtureData.Schools[0].ID,
	)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("school with name '%s' already exists", data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestUpdateSchoolReadOnly(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	data := dtos.SchoolDto{
		Name: "test",
	}

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPatch, "/schools/1")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"school with id '1' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateSchoolNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	data := dtos.SchoolDto{
		Name: "test",
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/8000",
	)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"school with id '8000' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateSchoolNotInt(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	data := dtos.SchoolDto{
		Name: "test",
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/aaaa",
	)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(
		t,
		"invalid URL param 'id' with value 'aaaa', should be an integer",
		rsData.Message,
	)
}

func TestUpdateSchoolFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/%d",
		fixtureData.Schools[0].ID,
	)
	tReq.AddCookie(tokens.ManagerAccessToken)

	data := dtos.SchoolDto{
		Name: "",
	}

	mt := test.CreateMatrixTester(tReq)
	mt.AddTestCaseErrorMessage(data, map[string]interface{}{
		"name": "must be provided",
	})

	mt.Do(t)
}

func TestUpdateSchoolAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/%d",
		fixtureData.Schools[0].ID,
	)

	mt := test.CreateMatrixTester(tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestDeleteSchool(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodDelete,
			"/schools/%d",
			fixtureData.Schools[i].ID,
		)
		tReq.AddCookie(user)

		var rsData models.School
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, fixtureData.Schools[i].ID, rsData.ID)
		assert.Equal(t, fixtureData.Schools[i].Name, rsData.Name)
		assert.Equal(t, false, rsData.ReadOnly)
	}
}

func TestDeleteSchoolReadOnly(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodDelete, "/schools/1")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"school with id '1' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteSchoolNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/schools/8000",
	)
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"school with id '8000' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteSchoolNotInt(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/schools/aaaa",
	)
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(
		t,
		"invalid URL param 'id' with value 'aaaa', should be an integer",
		rsData.Message.(string),
	)
}

func TestDeleteSchoolAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/schools/%d",
		fixtureData.Schools[0].ID,
	)

	mt := test.CreateMatrixTester(tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}
