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

	"github.com/google/uuid"

	"check-in/api/internal/assert"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"
)

func TestGetInfoLoggedInUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodGet, ts.URL+"/current-user", nil)
	req1.AddCookie(&tokens.AdminAccessToken)

	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/current-user", nil)
	req2.AddCookie(&tokens.ManagerAccessToken)

	req3, _ := http.NewRequest(http.MethodGet, ts.URL+"/current-user", nil)
	req3.AddCookie(&tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)
	rs3, _ := ts.Client().Do(req3)

	var rs1Data, rs2Data, rs3Data models.User
	_ = helpers.ReadJSON(rs1.Body, &rs1Data)
	_ = helpers.ReadJSON(rs2.Body, &rs2Data)
	_ = helpers.ReadJSON(rs3.Body, &rs3Data)

	assert.Equal(t, rs1.StatusCode, http.StatusOK)
	assert.Equal(t, rs1Data.ID, fixtureData.AdminUser.ID)
	assert.Equal(t, rs1Data.Username, fixtureData.AdminUser.Username)
	assert.Equal(t, rs1Data.Role, fixtureData.AdminUser.Role)
	assert.Equal(t, len(rs1Data.PasswordHash), 0)

	assert.Equal(t, rs2.StatusCode, http.StatusOK)
	assert.Equal(t, rs2Data.ID, fixtureData.ManagerUser.ID)
	assert.Equal(t, rs2Data.Username, fixtureData.ManagerUser.Username)
	assert.Equal(t, rs2Data.Role, fixtureData.ManagerUser.Role)
	assert.Equal(t, len(rs2Data.PasswordHash), 0)

	assert.Equal(t, rs3.StatusCode, http.StatusOK)
	assert.Equal(t, rs3Data.ID, fixtureData.DefaultUser.ID)
	assert.Equal(t, rs3Data.Username, fixtureData.DefaultUser.Username)
	assert.Equal(t, rs3Data.Role, fixtureData.DefaultUser.Role)
	assert.Equal(t, len(rs3Data.PasswordHash), 0)
}

func TestGetInfoLoggedInUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/current-user", nil)

	rs, _ := ts.Client().Do(req)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
}

func TestGetUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []http.Cookie{tokens.AdminAccessToken, tokens.ManagerAccessToken}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/users/"+fixtureData.DefaultUsers[0].ID,
			nil,
		)
		req.AddCookie(&user)

		rs, _ := ts.Client().Do(req)

		var rsData models.User
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.DefaultUsers[0].ID)
		assert.Equal(t, rsData.Username, fixtureData.DefaultUsers[0].Username)
		assert.Equal(t, rsData.Role, fixtureData.DefaultUsers[0].Role)
		assert.Equal(t, len(rsData.PasswordHash), 0)
	}
}

func TestGetUserNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users/"+id.String(), nil)
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
	)
}

func TestGetUserNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users/8000", nil)
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/users/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10),
		nil,
	)

	req2, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/users/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10),
		nil,
	)
	req2.AddCookie(&tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestGetPaginatedManagerUsersDefaultPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", nil)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.PaginatedUsersDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rsData.Pagination.Current, 1)
	assert.Equal(
		t,
		rsData.Pagination.Total,
		int64(math.Ceil(float64(fixtureData.AmountOfUsers)/4)),
	)
	assert.Equal(t, len(rsData.Data), 4)
}

func TestGetPaginatedManagerUsersSpecificPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users?page=2", nil)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.PaginatedUsersDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rsData.Pagination.Current, 2)
	assert.Equal(
		t,
		rsData.Pagination.Total,
		int64(math.Ceil(float64(fixtureData.AmountOfUsers)/4)),
	)
	assert.Equal(t, len(rsData.Data), 4)
}

func TestGetPaginatedManagerUsersPageZero(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users?page=0", nil)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message, "invalid page query param")
}

func TestGetPaginatedManagerUsersAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", nil)

	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", nil)
	req2.AddCookie(&tokens.DefaultAccessToken)

	req3, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", nil)
	req3.AddCookie(&tokens.ManagerAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)
	rs3, _ := ts.Client().Do(req3)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}

func TestCreateManagerUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CreateUserDto{
		Username: "test",
		Password: "testpassword",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", bytes.NewReader(body))
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData models.User
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusCreated)
	assert.IsUUID(t, rsData.ID)
	assert.Equal(t, rsData.Username, "test")
	assert.Equal(t, rsData.Role, models.ManagerRole)
	assert.Equal(t, len(rsData.PasswordHash), 0)
}

func TestCreateManagerUserUserNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CreateUserDto{
		Username: "TestManagerUser0",
		Password: "testpassword",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", bytes.NewReader(body))
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("user with username '%s' already exists", data.Username),
	)
}

func TestCreateManagerUserFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CreateUserDto{
		Username: "",
		Password: "",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", bytes.NewReader(body))
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["username"],
		"must be provided",
	)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["password"],
		"must be provided",
	)
}

func TestCreateManagerUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", nil)

	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", nil)
	req2.AddCookie(&tokens.DefaultAccessToken)

	req3, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", nil)
	req3.AddCookie(&tokens.ManagerAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)
	rs3, _ := ts.Client().Do(req3)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}

func TestUpdateManagerUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	username, password := "test", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		bytes.NewReader(body),
	)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData models.User
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rsData.ID, fixtureData.ManagerUsers[0].ID)
	assert.Equal(t, rsData.Username, "test")
	assert.Equal(t, rsData.Role, models.ManagerRole)
	assert.Equal(t, len(rsData.PasswordHash), 0)
}

func TestUpdateManagerUserUserNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	username, password := "TestManagerUser1", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		bytes.NewReader(body),
	)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("user with username '%s' already exists", *data.Username),
	)
}

func TestUpdateManagerUserNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	username, password := "test", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/users/"+id.String(),
		bytes.NewReader(body),
	)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
	)
}

func TestUpdateManagerUserNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	username, password := "test", "testpassword"
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)

	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/users/8000",
		bytes.NewReader(body),
	)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestUpdateManagerUserFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	username, password := "", ""
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		bytes.NewReader(body),
	)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["username"],
		"must be provided",
	)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["password"],
		"must be provided",
	)
}

func TestUpdateManagerUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		nil,
	)

	req2, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		nil,
	)
	req2.AddCookie(&tokens.DefaultAccessToken)

	req3, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		nil,
	)
	req3.AddCookie(&tokens.ManagerAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)
	rs3, _ := ts.Client().Do(req3)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}

func TestDeleteManagerUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		nil,
	)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData models.User
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rsData.ID, fixtureData.ManagerUsers[0].ID)
	assert.Equal(t, rsData.Username, fixtureData.ManagerUsers[0].Username)
	assert.Equal(t, rsData.Role, fixtureData.ManagerUsers[0].Role)
	assert.Equal(t, len(rsData.PasswordHash), 0)
}

func TestDeleteManagerUserNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/users/"+id.String(), nil)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
	)
}

func TestDeleteManagerUserNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/users/8000", nil)
	req.AddCookie(&tokens.AdminAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestDeleteManagerUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		nil,
	)

	req2, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		nil,
	)
	req2.AddCookie(&tokens.DefaultAccessToken)

	req3, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		nil,
	)
	req3.AddCookie(&tokens.ManagerAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)
	rs3, _ := ts.Client().Do(req3)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}
