package main

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"
)

func TestGetInfoLoggedInUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq1 := test.CreateTestRequest(t, ts, http.MethodGet, "/current-user")
	tReq1.AddCookie(tokens.AdminAccessToken)

	tReq2 := test.CreateTestRequest(t, ts, http.MethodGet, "/current-user")
	tReq2.AddCookie(tokens.ManagerAccessToken)

	tReq3 := test.CreateTestRequest(t, ts, http.MethodGet, "/current-user")
	tReq3.AddCookie(tokens.DefaultAccessToken)

	var rs1Data, rs2Data, rs3Data models.User
	rs1 := tReq1.Do(&rs1Data)
	rs2 := tReq2.Do(&rs2Data)
	rs3 := tReq3.Do(&rs3Data)

	assert.Equal(t, rs1.StatusCode, http.StatusOK)
	assert.Equal(t, rs1Data.ID, fixtureData.AdminUser.ID)
	assert.Equal(t, rs1Data.Username, fixtureData.AdminUser.Username)
	assert.Equal(t, rs1Data.Role, fixtureData.AdminUser.Role)
	assert.Equal(t, len(rs1Data.PasswordHash), 0)
	assert.Nil(t, rs1Data.Location)

	assert.Equal(t, rs2.StatusCode, http.StatusOK)
	assert.Equal(t, rs2Data.ID, fixtureData.ManagerUser.ID)
	assert.Equal(t, rs2Data.Username, fixtureData.ManagerUser.Username)
	assert.Equal(t, rs2Data.Role, fixtureData.ManagerUser.Role)
	assert.Equal(t, len(rs2Data.PasswordHash), 0)
	assert.Nil(t, rs2Data.Location)

	assert.Equal(t, rs3.StatusCode, http.StatusOK)
	assert.Equal(t, rs3Data.ID, fixtureData.DefaultUser.ID)
	assert.Equal(t, rs3Data.Username, fixtureData.DefaultUser.Username)
	assert.Equal(t, rs3Data.Role, fixtureData.DefaultUser.Role)
	assert.Equal(t, len(rs3Data.PasswordHash), 0)
	assert.Equal(t, rs3Data.Location.ID, fixtureData.DefaultLocation.ID)
	assert.Equal(t, rs3Data.Location.Name, fixtureData.DefaultLocation.Name)
	assert.Equal(
		t,
		rs3Data.Location.NormalizedName,
		fixtureData.DefaultLocation.NormalizedName,
	)
	assert.Equal(t, rs3Data.Location.Available, fixtureData.DefaultLocation.Available)
	assert.Equal(t, rs3Data.Location.Capacity, fixtureData.DefaultLocation.Capacity)
	assert.Equal(
		t,
		rs3Data.Location.YesterdayFullAt,
		fixtureData.DefaultLocation.YesterdayFullAt,
	)
	assert.Equal(t, rs3Data.Location.UserID, fixtureData.DefaultLocation.UserID)
}

func TestGetInfoLoggedInUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/current-user")
	rs := tReq.Do(nil)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
}

func TestGetUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/users/"+fixtureData.DefaultUsers[0].ID)
		tReq.AddCookie(user)

		var rsData models.User
		rs := tReq.Do(&rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.DefaultUsers[0].ID)
		assert.Equal(t, rsData.Username, fixtureData.DefaultUsers[0].Username)
		assert.Equal(t, rsData.Role, fixtureData.DefaultUsers[0].Role)
		assert.Equal(t, len(rsData.PasswordHash), 0)
		assert.Equal(t, rsData.Location.ID, fixtureData.Locations[0].ID)
		assert.Equal(t, rsData.Location.Name, fixtureData.Locations[0].Name)
		assert.Equal(
			t,
			rsData.Location.NormalizedName,
			fixtureData.Locations[0].NormalizedName,
		)
		assert.Equal(t, rsData.Location.Available, fixtureData.Locations[0].Available)
		assert.Equal(t, rsData.Location.Capacity, fixtureData.Locations[0].Capacity)
		assert.Equal(
			t,
			rsData.Location.YesterdayFullAt,
			fixtureData.Locations[0].YesterdayFullAt,
		)
		assert.Equal(t, rsData.Location.UserID, fixtureData.Locations[0].UserID)
	}
}

func TestGetUserNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/users/"+id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
	)
}

func TestGetUserNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/users/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq1 := test.CreateTestRequest(t, ts, http.MethodGet, "/users/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10))

	tReq2 := test.CreateTestRequest(t, ts, http.MethodGet, "/users/"+strconv.FormatInt(fixtureData.Schools[0].ID, 10))
	tReq2.AddCookie(tokens.DefaultAccessToken)

	rs1 := tReq1.Do(nil)
	rs2 := tReq2.Do(nil)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestGetPaginatedManagerUsersDefaultPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	var rsData dtos.PaginatedUsersDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.EqualValues(t, rsData.Pagination.Current, 1)
	assert.EqualValues(
		t,
		rsData.Pagination.Total,
		math.Ceil(float64(fixtureData.AmountOfManagerUsers)/4),
	)
	assert.Equal(t, len(rsData.Data), 4)

	assert.Equal(t, rsData.Data[0].ID, fixtureData.ManagerUser.ID)
	assert.Equal(t, rsData.Data[0].Username, fixtureData.ManagerUser.Username)
	assert.Equal(t, len(rsData.Data[0].PasswordHash), 0)
	assert.Equal(t, rsData.Data[0].Role, fixtureData.ManagerUser.Role)
	assert.Nil(t, rsData.Data[0].Location)
}

func TestGetPaginatedManagerUsersSpecificPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	var rsData dtos.PaginatedUsersDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.EqualValues(t, rsData.Pagination.Current, 2)
	assert.EqualValues(
		t,
		rsData.Pagination.Total,
		math.Ceil(float64(fixtureData.AmountOfManagerUsers)/4),
	)
	assert.Equal(t, len(rsData.Data), 4)

	assert.Equal(t, rsData.Data[0].ID, fixtureData.ManagerUsers[3].ID)
	assert.Equal(t, rsData.Data[0].Username, fixtureData.ManagerUsers[3].Username)
	assert.Equal(t, len(rsData.Data[0].PasswordHash), 0)
	assert.Equal(t, rsData.Data[0].Role, fixtureData.ManagerUsers[3].Role)
	assert.Nil(t, rsData.Data[0].Location)
}

func TestGetPaginatedManagerUsersPageZero(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "0",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message, "invalid page query param")
}

func TestGetPaginatedManagerUsersAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq1 := test.CreateTestRequest(t, ts, http.MethodGet, "/users")

	tReq2 := test.CreateTestRequest(t, ts, http.MethodGet, "/users")
	tReq2.AddCookie(tokens.DefaultAccessToken)

	tReq3 := test.CreateTestRequest(t, ts, http.MethodGet, "/users")
	tReq3.AddCookie(tokens.ManagerAccessToken)

	rs1 := tReq1.Do(nil)
	rs2 := tReq2.Do(nil)
	rs3 := tReq3.Do(nil)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}

func TestCreateManagerUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	data := dtos.CreateUserDto{
		Username: "test",
		Password: "testpassword",
	}
	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusCreated)
	assert.Nil(t, uuid.Validate(rsData.ID))
	assert.Equal(t, rsData.Username, "test")
	assert.Equal(t, rsData.Role, models.ManagerRole)
	assert.Equal(t, len(rsData.PasswordHash), 0)
	assert.Nil(t, rsData.Location)
}

func TestCreateManagerUserUserNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	data := dtos.CreateUserDto{
		Username: "TestManagerUser0",
		Password: "testpassword",
	}
	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["username"].(string),
		fmt.Sprintf("user with username '%s' already exists", data.Username),
	)
}

func TestCreateManagerUserFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/users")
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(dtos.CreateUserDto{
		Username: "",
		Password: "",
	})

	vt := test.CreateValidatorTester(t)
	vt.AddTestCase(tReq, map[string]interface{}{
		"username": "must be provided",
		"password": "must be provided",
	})

	vt.Do()
}

func TestCreateManagerUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq1 := test.CreateTestRequest(t, ts, http.MethodPost, "/users")

	tReq2 := test.CreateTestRequest(t, ts, http.MethodPost, "/users")
	tReq2.AddCookie(tokens.DefaultAccessToken)

	tReq3 := test.CreateTestRequest(t, ts, http.MethodPost, "/users")
	tReq3.AddCookie(tokens.ManagerAccessToken)

	rs1 := tReq1.Do(nil)
	rs2 := tReq2.Do(nil)
	rs3 := tReq3.Do(nil)

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

	tReq := test.CreateTestRequest(t, ts, http.MethodPatch, "/users/"+fixtureData.ManagerUsers[0].ID)
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rsData.ID, fixtureData.ManagerUsers[0].ID)
	assert.Equal(t, rsData.Username, "test")
	assert.Equal(t, rsData.Role, models.ManagerRole)
	assert.Equal(t, len(rsData.PasswordHash), 0)
	assert.Nil(t, rsData.Location)
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

	tReq := test.CreateTestRequest(t, ts, http.MethodPatch, "/users/"+fixtureData.ManagerUsers[0].ID)
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["username"].(string),
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

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(t, ts, http.MethodPatch, "/users/"+id.String())
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
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

	tReq := test.CreateTestRequest(t, ts, http.MethodPatch, "/users/8000")
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

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

	tReq := test.CreateTestRequest(t, ts, http.MethodPatch, "/users/"+fixtureData.ManagerUsers[0].ID)
	tReq.AddCookie(tokens.AdminAccessToken)

	tReq.SetReqData(data)

	vt := test.CreateValidatorTester(t)
	vt.AddTestCase(tReq, map[string]interface{}{
		"username": "must be provided",
		"password": "must be provided",
	})

	vt.Do()
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
	req2.AddCookie(tokens.DefaultAccessToken)

	req3, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		nil,
	)
	req3.AddCookie(tokens.ManagerAccessToken)

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

	tReq := test.CreateTestRequest(t, ts, http.MethodDelete, "/users/"+fixtureData.ManagerUsers[0].ID)
	tReq.AddCookie(tokens.AdminAccessToken)

	var rsData models.User
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rsData.ID, fixtureData.ManagerUsers[0].ID)
	assert.Equal(t, rsData.Username, fixtureData.ManagerUsers[0].Username)
	assert.Equal(t, rsData.Role, fixtureData.ManagerUsers[0].Role)
	assert.Equal(t, len(rsData.PasswordHash), 0)
	assert.Nil(t, rsData.Location)
}

func TestDeleteManagerUserNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	id, _ := uuid.NewUUID()
	tReq := test.CreateTestRequest(t, ts, http.MethodDelete, "/users/"+id.String())
	tReq.AddCookie(tokens.AdminAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("user with id '%s' doesn't exist", id.String()),
	)
}

func TestDeleteManagerUserNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodDelete, "/users/8000")
	tReq.AddCookie(tokens.AdminAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

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
	req2.AddCookie(tokens.DefaultAccessToken)

	req3, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/users/"+fixtureData.ManagerUsers[0].ID,
		nil,
	)
	req3.AddCookie(tokens.ManagerAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)
	rs3, _ := ts.Client().Do(req3)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}
