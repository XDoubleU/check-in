package main

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/test"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestGetInfoLoggedInUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq1 := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/current-user")
	tReq1.AddCookie(testEnv.Tokens.AdminAccessToken)

	tReq2 := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/current-user")
	tReq2.AddCookie(testEnv.Tokens.ManagerAccessToken)

	tReq3 := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/current-user")
	tReq3.AddCookie(testEnv.Tokens.DefaultAccessToken)

	var rs1Data, rs2Data, rs3Data models.User
	rs1 := tReq1.Do(t, &rs1Data)
	rs2 := tReq2.Do(t, &rs2Data)
	rs3 := tReq3.Do(t, &rs3Data)

	assert.Equal(t, http.StatusOK, rs1.StatusCode)
	assert.Equal(t, testEnv.Fixtures.AdminUser.ID, rs1Data.ID)
	assert.Equal(t, testEnv.Fixtures.AdminUser.Username, rs1Data.Username)
	assert.Equal(t, testEnv.Fixtures.AdminUser.Role, rs1Data.Role)
	assert.Equal(t, 0, len(rs1Data.PasswordHash))
	assert.Nil(t, rs1Data.Location)

	assert.Equal(t, http.StatusOK, rs2.StatusCode)
	assert.Equal(t, testEnv.Fixtures.ManagerUser.ID, rs2Data.ID)
	assert.Equal(t, testEnv.Fixtures.ManagerUser.Username, rs2Data.Username)
	assert.Equal(t, testEnv.Fixtures.ManagerUser.Role, rs2Data.Role)
	assert.Equal(t, 0, len(rs2Data.PasswordHash))
	assert.Nil(t, rs2Data.Location)

	assert.Equal(t, http.StatusOK, rs3.StatusCode)
	assert.Equal(t, testEnv.Fixtures.DefaultUser.ID, rs3Data.ID)
	assert.Equal(t, testEnv.Fixtures.DefaultUser.Username, rs3Data.Username)
	assert.Equal(t, testEnv.Fixtures.DefaultUser.Role, rs3Data.Role)
	assert.Equal(t, 0, len(rs3Data.PasswordHash))
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.ID, rs3Data.Location.ID)
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.Name, rs3Data.Location.Name)
	assert.Equal(
		t,
		testEnv.Fixtures.DefaultLocation.NormalizedName,
		rs3Data.Location.NormalizedName,
	)
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.Available, rs3Data.Location.Available)
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.Capacity, rs3Data.Location.Capacity)
	assert.Equal(
		t,
		testEnv.Fixtures.DefaultLocation.YesterdayFullAt,
		rs3Data.Location.YesterdayFullAt,
	)
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.UserID, rs3Data.Location.UserID)
}

func TestGetInfoLoggedInUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/current-user")

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReq, test.NewCaseResponse(http.StatusUnauthorized))

	mt.Do(t)
}

func TestGetUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	users := []*http.Cookie{
		testEnv.Tokens.AdminAccessToken,
		testEnv.Tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/users/%s",
			testEnv.Fixtures.DefaultUsers[0].ID,
		)
		tReq.AddCookie(user)

		var rsData models.User
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, testEnv.Fixtures.DefaultUsers[0].ID, rsData.ID)
		assert.Equal(t, testEnv.Fixtures.DefaultUsers[0].Username, rsData.Username)
		assert.Equal(t, testEnv.Fixtures.DefaultUsers[0].Role, rsData.Role)
		assert.Equal(t, 0, len(rsData.PasswordHash))
		assert.Equal(t, testEnv.Fixtures.Locations[0].ID, rsData.Location.ID)
		assert.Equal(t, testEnv.Fixtures.Locations[0].Name, rsData.Location.Name)
		assert.Equal(
			t,
			testEnv.Fixtures.Locations[0].NormalizedName,
			rsData.Location.NormalizedName,
		)
		assert.Equal(t, testEnv.Fixtures.Locations[0].Available, rsData.Location.Available)
		assert.Equal(t, testEnv.Fixtures.Locations[0].Capacity, rsData.Location.Capacity)
		assert.Equal(
			t,
			testEnv.Fixtures.Locations[0].YesterdayFullAt,
			rsData.Location.YesterdayFullAt,
		)
		assert.Equal(t, testEnv.Fixtures.Locations[0].UserID, rsData.Location.UserID)
	}
}

func TestGetUserNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	id, _ := uuid.NewUUID()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/users/%s",
		id.String(),
	)
	tReq.AddCookie(testEnv.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetUserNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/users/8000")
	tReq.AddCookie(testEnv.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestGetUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/users/%d",
		testEnv.Fixtures.Schools[0].ID,
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestGetPaginatedManagerUsersDefaultPage(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/users")
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	var rsData dtos.PaginatedUsersDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 1, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(testEnv.Fixtures.AmountOfManagerUsers)/4),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 4, len(rsData.Data))

	assert.Equal(t, testEnv.Fixtures.ManagerUser.ID, rsData.Data[0].ID)
	assert.Equal(t, testEnv.Fixtures.ManagerUser.Username, rsData.Data[0].Username)
	assert.Equal(t, 0, len(rsData.Data[0].PasswordHash))
	assert.Equal(t, testEnv.Fixtures.ManagerUser.Role, rsData.Data[0].Role)
	assert.Nil(t, rsData.Data[0].Location)
}

func TestGetPaginatedManagerUsersSpecificPage(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/users")
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	var rsData dtos.PaginatedUsersDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 2, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(testEnv.Fixtures.AmountOfManagerUsers)/4),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 4, len(rsData.Data))

	assert.Equal(t, testEnv.Fixtures.ManagerUsers[3].ID, rsData.Data[0].ID)
	assert.Equal(t, testEnv.Fixtures.ManagerUsers[3].Username, rsData.Data[0].Username)
	assert.Equal(t, 0, len(rsData.Data[0].PasswordHash))
	assert.Equal(t, testEnv.Fixtures.ManagerUsers[3].Role, rsData.Data[0].Role)
	assert.Nil(t, rsData.Data[0].Location)
}

func TestGetPaginatedManagerUsersPageFull(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/users")
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	test.PaginatedEndpointTester(
		t,
		tReq,
		"page",
		int(math.Ceil(float64(testEnv.Fixtures.AmountOfManagerUsers)/4)),
	)
}

func TestGetPaginatedManagerUsersAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/users")

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(testEnv.Tokens.ManagerAccessToken)

	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestCreateManagerUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/users")
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	data := dtos.CreateUserDto{
		Username: "test",
		Password: "testpassword",
	}
	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Nil(t, uuid.Validate(rsData.ID))
	assert.Equal(t, "test", rsData.Username)
	assert.Equal(t, models.ManagerRole, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
	assert.Nil(t, rsData.Location)
}

func TestCreateManagerUserUserNameExists(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/users")
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	data := dtos.CreateUserDto{
		Username: "TestManagerUser0",
		Password: "testpassword",
	}
	tReq.SetReqData(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with username '%s' already exists", data.Username),
		rsData.Message.(map[string]interface{})["username"].(string),
	)
}

func TestCreateManagerUserFailValidation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/users")
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)
	tReq.SetReqData(dtos.CreateUserDto{
		Username: "",
		Password: "",
	})

	mt := test.CreateMatrixTester()

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes.SetExpectedBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"username": "must be provided",
			"password": "must be provided",
		}),
	)

	mt.AddTestCase(tReq, tRes)

	mt.Do(t)
}

func TestCreateManagerUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/users")

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(testEnv.Tokens.ManagerAccessToken)

	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestUpdateManagerUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	username, password := "test", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/users/%s",
		testEnv.Fixtures.ManagerUsers[0].ID,
	)
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, testEnv.Fixtures.ManagerUsers[0].ID, rsData.ID)
	assert.Equal(t, "test", rsData.Username)
	assert.Equal(t, models.ManagerRole, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
	assert.Nil(t, rsData.Location)
}

func TestUpdateManagerUserUserNameExists(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	username, password := "TestManagerUser1", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/users/%s",
		testEnv.Fixtures.ManagerUsers[0].ID,
	)
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with username '%s' already exists", *data.Username),
		rsData.Message.(map[string]interface{})["username"].(string),
	)
}

func TestUpdateManagerUserNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	username, password := "test", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	id, _ := uuid.NewUUID()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/users/%s",
		id.String(),
	)
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateManagerUserNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	username, password := "test", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPatch, "/users/8000")
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestUpdateManagerUserFailValidation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	username, password := "", ""

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/users/%s",
		testEnv.Fixtures.ManagerUsers[0].ID,
	)
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)
	tReq.SetReqData(dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	})

	mt := test.CreateMatrixTester()

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes.SetExpectedBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"username": "must be provided",
			"password": "must be provided",
		}),
	)

	mt.AddTestCase(tReq, tRes)

	mt.Do(t)
}

func TestUpdateManagerUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/users/%s",
		testEnv.Fixtures.ManagerUsers[0].ID,
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(testEnv.Tokens.ManagerAccessToken)

	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestDeleteManagerUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/users/%s",
		testEnv.Fixtures.ManagerUsers[0].ID,
	)
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, testEnv.Fixtures.ManagerUsers[0].ID, rsData.ID)
	assert.Equal(t, testEnv.Fixtures.ManagerUsers[0].Username, rsData.Username)
	assert.Equal(t, testEnv.Fixtures.ManagerUsers[0].Role, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
	assert.Nil(t, rsData.Location)
}

func TestDeleteManagerUserNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	id, _ := uuid.NewUUID()
	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/users/%s",
		id.String(),
	)
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteManagerUserNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodDelete, "/users/8000")
	tReq.AddCookie(testEnv.Tokens.AdminAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestDeleteManagerUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/users/%s",
		testEnv.Fixtures.ManagerUsers[0].ID,
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(testEnv.Tokens.ManagerAccessToken)

	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}
