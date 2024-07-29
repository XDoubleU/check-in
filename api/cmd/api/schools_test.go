package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/test"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestGetPaginatedSchoolsDefaultPage(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	testEnv.createSchools(20)
	amount, err := testEnv.services.Schools.GetTotalCount(context.Background())
	require.Nil(t, err)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
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
			math.Ceil(float64(*amount)/4),
			rsData.Pagination.Total,
		)
		assert.Equal(t, 4, len(rsData.Data))

		assert.EqualValues(t, 1, rsData.Data[0].ID)
		assert.Equal(t, "Andere", rsData.Data[0].Name)
		assert.Equal(t, true, rsData.Data[0].ReadOnly)
	}
}

func TestGetPaginatedSchoolsSpecificPage(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	schools := testEnv.createSchools(20)
	amount, err := testEnv.services.Schools.GetTotalCount(context.Background())
	require.Nil(t, err)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/schools")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	var rsData dtos.PaginatedSchoolsDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 2, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(*amount)/4),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 4, len(rsData.Data))

	assert.Equal(t, schools[11].ID, rsData.Data[0].ID)
	assert.Equal(t, schools[11].Name, rsData.Data[0].Name)
	assert.Equal(t, schools[11].ReadOnly, rsData.Data[0].ReadOnly)
}

func TestGetPaginatedSchoolsFull(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	amount := 20
	testEnv.createSchools(amount)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/schools")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	test.PaginatedEndpointTester(
		t,
		tReq,
		"page",
		int(math.Ceil(float64(amount)/4)),
	)
}

func TestGetPaginatedSchoolsAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/schools")

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestCreateSchool(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
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
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	data := dtos.SchoolDto{
		Name: "Andere",
	}

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/schools")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("school with name '%s' already exists", data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestCreateSchoolFailValidation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/schools")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)
	tReq.SetReqData(dtos.SchoolDto{
		Name: "",
	})

	mt := test.CreateMatrixTester()

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes.SetExpectedBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"name": "must be provided",
		}),
	)

	mt.AddTestCase(tReq, tRes)

	mt.Do(t)
}

func TestCreateSchoolAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/schools")

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestUpdateSchool(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
	}

	for i, user := range users {
		school := testEnv.createSchools(1)[0]

		unique := fmt.Sprintf("test%d", i)

		data := dtos.SchoolDto{
			Name: unique,
		}

		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodPatch,
			"/schools/%d",
			school.ID,
		)
		tReq.AddCookie(user)

		tReq.SetReqData(data)

		var rsData models.School
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, school.ID, rsData.ID)
		assert.Equal(t, unique, rsData.Name)
		assert.Equal(t, false, rsData.ReadOnly)
	}
}

func TestUpdateSchoolNameExists(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	data := dtos.SchoolDto{
		Name: "Andere",
	}

	school := testEnv.createSchools(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/%d",
		school.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("school with name '%s' already exists", data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestUpdateSchoolReadOnly(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	data := dtos.SchoolDto{
		Name: "test",
	}

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPatch, "/schools/1")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"school with id '1' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateSchoolNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	data := dtos.SchoolDto{
		Name: "test",
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/8000",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"school with id '8000' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateSchoolNotInt(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	data := dtos.SchoolDto{
		Name: "test",
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/aaaa",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(
		t,
		"invalid URL param 'id' with value 'aaaa', should be an integer",
		rsData.Message,
	)
}

func TestUpdateSchoolFailValidation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	school := testEnv.createSchools(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/%d",
		school.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)
	tReq.SetReqData(dtos.SchoolDto{
		Name: "",
	})

	mt := test.CreateMatrixTester()

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes.SetExpectedBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"name": "must be provided",
		}),
	)

	mt.AddTestCase(tReq, tRes)

	mt.Do(t)
}

func TestUpdateSchoolAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	school := testEnv.createSchools(1)[0]

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/schools/%d",
		school.ID,
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestDeleteSchool(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
	}

	for _, user := range users {
		school := testEnv.createSchools(1)[0]

		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodDelete,
			"/schools/%d",
			school.ID,
		)
		tReq.AddCookie(user)

		var rsData models.School
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, school.ID, rsData.ID)
		assert.Equal(t, school.Name, rsData.Name)
		assert.Equal(t, false, rsData.ReadOnly)
	}
}

func TestDeleteSchoolReadOnly(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodDelete, "/schools/1")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"school with id '1' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteSchoolNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/schools/8000",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"school with id '8000' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteSchoolNotInt(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/schools/aaaa",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(
		t,
		"invalid URL param 'id' with value 'aaaa', should be an integer",
		rsData.Message.(string),
	)
}

func TestDeleteSchoolAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	school := testEnv.createSchools(1)[0]

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/schools/%d",
		school.ID,
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}
