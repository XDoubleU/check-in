package main

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestGetInfoLoggedInUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq1 := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/current-user")
	tReq1.AddCookie(tokens.AdminAccessToken)

	tReq2 := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/current-user")
	tReq2.AddCookie(tokens.ManagerAccessToken)

	tReq3 := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/current-user")
	tReq3.AddCookie(tokens.DefaultAccessToken)

	var rs1Data, rs2Data, rs3Data models.User
	rs1 := tReq1.Do(t, &rs1Data)
	rs2 := tReq2.Do(t, &rs2Data)
	rs3 := tReq3.Do(t, &rs3Data)

	assert.Equal(t, http.StatusOK, rs1.StatusCode)
	assert.Equal(t, fixtureData.AdminUser.ID, rs1Data.ID)
	assert.Equal(t, fixtureData.AdminUser.Username, rs1Data.Username)
	assert.Equal(t, fixtureData.AdminUser.Role, rs1Data.Role)
	assert.Equal(t, 0, len(rs1Data.PasswordHash))
	assert.Nil(t, rs1Data.Location)

	assert.Equal(t, http.StatusOK, rs2.StatusCode)
	assert.Equal(t, fixtureData.ManagerUser.ID, rs2Data.ID)
	assert.Equal(t, fixtureData.ManagerUser.Username, rs2Data.Username)
	assert.Equal(t, fixtureData.ManagerUser.Role, rs2Data.Role)
	assert.Equal(t, 0, len(rs2Data.PasswordHash))
	assert.Nil(t, rs2Data.Location)

	assert.Equal(t, http.StatusOK, rs3.StatusCode)
	assert.Equal(t, fixtureData.DefaultUser.ID, rs3Data.ID)
	assert.Equal(t, fixtureData.DefaultUser.Username, rs3Data.Username)
	assert.Equal(t, fixtureData.DefaultUser.Role, rs3Data.Role)
	assert.Equal(t, 0, len(rs3Data.PasswordHash))
	assert.Equal(t, fixtureData.DefaultLocation.ID, rs3Data.Location.ID)
	assert.Equal(t, fixtureData.DefaultLocation.Name, rs3Data.Location.Name)
	assert.Equal(
		t,
		fixtureData.DefaultLocation.NormalizedName,
		rs3Data.Location.NormalizedName,
	)
	assert.Equal(t, fixtureData.DefaultLocation.Available, rs3Data.Location.Available)
	assert.Equal(t, fixtureData.DefaultLocation.Capacity, rs3Data.Location.Capacity)
	assert.Equal(
		t,
		fixtureData.DefaultLocation.YesterdayFullAt,
		rs3Data.Location.YesterdayFullAt,
	)
	assert.Equal(t, fixtureData.DefaultLocation.UserID, rs3Data.Location.UserID)
}

func TestGetInfoLoggedInUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/current-user")

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)

	mt.Do(t)
}

func TestGetUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/users/%s", fixtureData.DefaultUsers[0].ID)
		tReq.AddCookie(user)

		var rsData models.User
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, fixtureData.DefaultUsers[0].ID, rsData.ID)
		assert.Equal(t, fixtureData.DefaultUsers[0].Username, rsData.Username)
		assert.Equal(t, fixtureData.DefaultUsers[0].Role, rsData.Role)
		assert.Equal(t, 0, len(rsData.PasswordHash))
		assert.Equal(t, fixtureData.Locations[0].ID, rsData.Location.ID)
		assert.Equal(t, fixtureData.Locations[0].Name, rsData.Location.Name)
		assert.Equal(
			t,
			fixtureData.Locations[0].NormalizedName,
			rsData.Location.NormalizedName,
		)
		assert.Equal(t, fixtureData.Locations[0].Available, rsData.Location.Available)
		assert.Equal(t, fixtureData.Locations[0].Capacity, rsData.Location.Capacity)
		assert.Equal(
			t,
			fixtureData.Locations[0].YesterdayFullAt,
			rsData.Location.YesterdayFullAt,
		)
		assert.Equal(t, fixtureData.Locations[0].UserID, rsData.Location.UserID)
	}
}

func TestGetUserNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/users/%s", id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetUserNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/users/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/users/%d", fixtureData.Schools[0].ID)

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestGetPaginatedManagerUsersDefaultPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	var rsData dtos.PaginatedUsersDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 1, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(fixtureData.AmountOfManagerUsers)/4),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 4, len(rsData.Data))

	assert.Equal(t, fixtureData.ManagerUser.ID, rsData.Data[0].ID)
	assert.Equal(t, fixtureData.ManagerUser.Username, rsData.Data[0].Username)
	assert.Equal(t, 0, len(rsData.Data[0].PasswordHash))
	assert.Equal(t, fixtureData.ManagerUser.Role, rsData.Data[0].Role)
	assert.Nil(t, rsData.Data[0].Location)
}

func TestGetPaginatedManagerUsersSpecificPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	var rsData dtos.PaginatedUsersDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 2, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(fixtureData.AmountOfManagerUsers)/4),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 4, len(rsData.Data))

	assert.Equal(t, fixtureData.ManagerUsers[3].ID, rsData.Data[0].ID)
	assert.Equal(t, fixtureData.ManagerUsers[3].Username, rsData.Data[0].Username)
	assert.Equal(t, 0, len(rsData.Data[0].PasswordHash))
	assert.Equal(t, fixtureData.ManagerUsers[3].Role, rsData.Data[0].Role)
	assert.Nil(t, rsData.Data[0].Location)
}

func TestGetPaginatedManagerUsersPageFull(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	test.TestPaginatedEndpoint(t, tReq, "page", int(math.Ceil(float64(fixtureData.AmountOfManagerUsers)/4)))
}

func TestGetPaginatedManagerUsersAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/users")

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)
	mt.AddTestCaseCookieStatusCode(tokens.ManagerAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestCreateManagerUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

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
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	data := dtos.CreateUserDto{
		Username: "TestManagerUser0",
		Password: "testpassword",
	}
	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with username '%s' already exists", data.Username),
		rsData.Message.(map[string]interface{})["username"].(string),
	)
}

func TestCreateManagerUserFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	data := dtos.CreateUserDto{
		Username: "",
		Password: "",
	}

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseErrorMessage(data, map[string]interface{}{
		"username": "must be provided",
		"password": "must be provided",
	})

	mt.Do(t)
}

func TestCreateManagerUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/users")

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)
	mt.AddTestCaseCookieStatusCode(tokens.ManagerAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestUpdateManagerUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	username, password := "test", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/users/%s", fixtureData.ManagerUsers[0].ID)
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, fixtureData.ManagerUsers[0].ID, rsData.ID)
	assert.Equal(t, "test", rsData.Username)
	assert.Equal(t, models.ManagerRole, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
	assert.Nil(t, rsData.Location)
}

func TestUpdateManagerUserUserNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	username, password := "TestManagerUser1", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/users/%s", fixtureData.ManagerUsers[0].ID)
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with username '%s' already exists", *data.Username),
		rsData.Message.(map[string]interface{})["username"].(string),
	)
}

func TestUpdateManagerUserNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	username, password := "test", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/users/%s", id.String())
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateManagerUserNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	username, password := "test", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/users/8000")
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestUpdateManagerUserFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	username, password := "", ""
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/users/%s", fixtureData.ManagerUsers[0].ID)
	tReq.AddCookie(tokens.AdminAccessToken)

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseErrorMessage(data, map[string]interface{}{
		"username": "must be provided",
		"password": "must be provided",
	})

	mt.Do(t)
}

func TestUpdateManagerUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/users/%s", fixtureData.ManagerUsers[0].ID)

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)
	mt.AddTestCaseCookieStatusCode(tokens.ManagerAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestDeleteManagerUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/users/%s", fixtureData.ManagerUsers[0].ID)
	tReq.AddCookie(tokens.AdminAccessToken)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, fixtureData.ManagerUsers[0].ID, rsData.ID)
	assert.Equal(t, fixtureData.ManagerUsers[0].Username, rsData.Username)
	assert.Equal(t, fixtureData.ManagerUsers[0].Role, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
	assert.Nil(t, rsData.Location)
}

func TestDeleteManagerUserNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	id, _ := uuid.NewUUID()
	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/users/%s", id.String())
	tReq.AddCookie(tokens.AdminAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteManagerUserNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/users/8000")
	tReq.AddCookie(tokens.AdminAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestDeleteManagerUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/users/%s", fixtureData.ManagerUsers[0].ID)

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)
	mt.AddTestCaseCookieStatusCode(tokens.ManagerAccessToken, http.StatusForbidden)

	mt.Do(t)
}
