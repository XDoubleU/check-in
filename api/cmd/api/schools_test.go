package main

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestGetPaginatedSchoolsDefaultPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/schools")
		tReq.AddCookie(user)

		var rsData dtos.PaginatedSchoolsDto
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		assert.EqualValues(t, rsData.Pagination.Current, 1)
		assert.EqualValues(
			t,
			rsData.Pagination.Total,
			math.Ceil(float64(fixtureData.AmountOfLocations)/4),
		)
		assert.Equal(t, len(rsData.Data), 4)

		assert.EqualValues(t, rsData.Data[0].ID, 1)
		assert.Equal(t, rsData.Data[0].Name, "Andere")
		assert.Equal(t, rsData.Data[0].ReadOnly, true)
	}
}

func TestGetPaginatedSchoolsSpecificPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/schools")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	var rsData dtos.PaginatedSchoolsDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.EqualValues(t, rsData.Pagination.Current, 2)
	assert.EqualValues(
		t,
		rsData.Pagination.Total,
		math.Ceil(float64(fixtureData.AmountOfLocations)/4),
	)
	assert.Equal(t, len(rsData.Data), 4)

	assert.Equal(t, rsData.Data[0].ID, fixtureData.Schools[11].ID)
	assert.Equal(t, rsData.Data[0].Name, fixtureData.Schools[11].Name)
	assert.Equal(t, rsData.Data[0].ReadOnly, fixtureData.Schools[11].ReadOnly)
}

func TestGetPaginatedSchoolsPageZero(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/schools")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "0",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message, "invalid query param 'page' with value '0', can't be '0'")
}

func TestGetPaginatedSchoolsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq1 := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/schools")

	tReq2 := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/schools")
	tReq2.AddCookie(tokens.DefaultAccessToken)

	rs1 := tReq1.Do(t, nil)
	rs2 := tReq2.Do(t, nil)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestCreateSchool(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

		tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/schools")
		tReq.AddCookie(user)

		tReq.SetReqData(data)

		var rsData models.School
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusCreated)
		assert.Equal(t, rsData.Name, unique)
		assert.Equal(t, rsData.ReadOnly, false)
	}
}

func TestCreateSchoolNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "TestSchool0",
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/schools")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["name"].(string),
		fmt.Sprintf("school with name '%s' already exists", data.Name),
	)
}

func TestCreateSchoolFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/schools")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(dtos.SchoolDto{
		Name: "",
	})

	vt := test.CreateValidatorTester(t)
	vt.AddTestCase(tReq, map[string]interface{}{
		"name": "must be provided",
	})

	vt.Do(t)
}

func TestCreateSchoolAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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
	defer test.TeardownSingle(testEnv)

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

		tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/schools/%d", fixtureData.Schools[0].ID)
		tReq.AddCookie(user)

		tReq.SetReqData(data)

		var rsData models.School
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.Schools[0].ID)
		assert.Equal(t, rsData.Name, unique)
		assert.Equal(t, rsData.ReadOnly, false)
	}
}

func TestUpdateSchoolNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "TestSchool1",
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/schools/%d", fixtureData.Schools[0].ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["name"].(string),
		fmt.Sprintf("school with name '%s' already exists", data.Name),
	)
}

func TestUpdateSchoolReadOnly(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "test",
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/schools/1")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		"school with id '1' doesn't exist",
	)
}

func TestUpdateSchoolNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "test",
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/schools/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		"school with id '8000' doesn't exist",
	)
}

func TestUpdateSchoolNotInt(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SchoolDto{
		Name: "test",
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/schools/aaaa")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message, "invalid URL param 'id' with value 'aaaa', should be an integer")
}

func TestUpdateSchoolFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/schools/%d", fixtureData.Schools[0].ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(dtos.SchoolDto{
		Name: "",
	})

	vt := test.CreateValidatorTester(t)
	vt.AddTestCase(tReq, map[string]interface{}{
		"name": "must be provided",
	})

	vt.Do(t)
}

func TestUpdateSchoolAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/schools/%d", fixtureData.Schools[i].ID)
		tReq.AddCookie(user)

		var rsData models.School
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.Schools[i].ID)
		assert.Equal(t, rsData.Name, fixtureData.Schools[i].Name)
		assert.Equal(t, rsData.ReadOnly, false)
	}
}

func TestDeleteSchoolReadOnly(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/schools/1")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		"school with id '1' doesn't exist",
	)
}

func TestDeleteSchoolNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/schools/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		"school with id '8000' doesn't exist",
	)
}

func TestDeleteSchoolNotInt(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/schools/aaaa")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "invalid URL param 'id' with value 'aaaa', should be an integer")
}

func TestDeleteSchoolAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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
