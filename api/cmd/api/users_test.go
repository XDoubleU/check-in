package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"testing"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	errortools "github.com/XDoubleU/essentia/pkg/errors"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestGetInfoLoggedInUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq1 := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/current-user")
	tReq1.AddCookie(fixtures.Tokens.AdminAccessToken)

	tReq2 := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/current-user")
	tReq2.AddCookie(fixtures.Tokens.ManagerAccessToken)

	tReq3 := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/current-user")
	tReq3.AddCookie(fixtures.Tokens.DefaultAccessToken)

	var rs1Data, rs2Data, rs3Data models.User
	rs1 := tReq1.Do(t)
	rs2 := tReq2.Do(t)
	rs3 := tReq3.Do(t)

	err := httptools.ReadJSON(rs1.Body, &rs1Data)
	require.Nil(t, err)
	err = httptools.ReadJSON(rs2.Body, &rs2Data)
	require.Nil(t, err)
	err = httptools.ReadJSON(rs3.Body, &rs3Data)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs1.StatusCode)
	assert.Equal(t, fixtures.AdminUser.ID, rs1Data.ID)
	assert.Equal(t, fixtures.AdminUser.Username, rs1Data.Username)
	assert.Equal(t, fixtures.AdminUser.Role, rs1Data.Role)
	assert.Equal(t, 0, len(rs1Data.PasswordHash))
	assert.Nil(t, rs1Data.Location)

	assert.Equal(t, http.StatusOK, rs2.StatusCode)
	assert.Equal(t, fixtures.ManagerUser.ID, rs2Data.ID)
	assert.Equal(t, fixtures.ManagerUser.Username, rs2Data.Username)
	assert.Equal(t, fixtures.ManagerUser.Role, rs2Data.Role)
	assert.Equal(t, 0, len(rs2Data.PasswordHash))
	assert.Nil(t, rs2Data.Location)

	assert.Equal(t, http.StatusOK, rs3.StatusCode)
	assert.Equal(t, fixtures.DefaultUser.ID, rs3Data.ID)
	assert.Equal(t, fixtures.DefaultUser.Username, rs3Data.Username)
	assert.Equal(t, fixtures.DefaultUser.Role, rs3Data.Role)
	assert.Equal(t, 0, len(rs3Data.PasswordHash))
	assert.Equal(t, fixtures.DefaultLocation.ID, rs3Data.Location.ID)
	assert.Equal(t, fixtures.DefaultLocation.Name, rs3Data.Location.Name)
	assert.Equal(
		t,
		fixtures.DefaultLocation.NormalizedName,
		rs3Data.Location.NormalizedName,
	)
	assert.Equal(
		t,
		fixtures.DefaultLocation.Available,
		rs3Data.Location.Available,
	)
	assert.Equal(
		t,
		fixtures.DefaultLocation.Capacity,
		rs3Data.Location.Capacity,
	)
	assert.Equal(
		t,
		fixtures.DefaultLocation.YesterdayFullAt,
		rs3Data.Location.YesterdayFullAt,
	)
	assert.Equal(t, fixtures.DefaultLocation.UserID, rs3Data.Location.UserID)
}

func TestGetInfoLoggedInUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/current-user")

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReq, test.NewCaseResponse(http.StatusUnauthorized, nil, nil))

	mt.Do(t)
}

func TestGetUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	userID := testEnv.createLocations(1)[0].UserID
	defaultUser, _ := testApp.services.Locations.GetDefaultUserByUserID(
		context.Background(),
		userID,
	)

	users := []*http.Cookie{
		fixtures.Tokens.AdminAccessToken,
		fixtures.Tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/users/%s",
			userID,
		)
		tReq.AddCookie(user)

		rs := tReq.Do(t)

		var rsData models.User
		err := httptools.ReadJSON(rs.Body, &rsData)
		require.Nil(t, err)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, defaultUser.ID, rsData.ID)
		assert.Equal(t, defaultUser.Username, rsData.Username)
		assert.Equal(t, defaultUser.Role, rsData.Role)
		assert.Equal(t, 0, len(rsData.PasswordHash))
		assert.Equal(t, defaultUser.Location.ID, rsData.Location.ID)
		assert.Equal(t, defaultUser.Location.Name, rsData.Location.Name)
		assert.Equal(
			t,
			defaultUser.Location.NormalizedName,
			rsData.Location.NormalizedName,
		)
		assert.Equal(t, defaultUser.Location.Available, rsData.Location.Available)
		assert.Equal(t, defaultUser.Location.Capacity, rsData.Location.Capacity)
		assert.Equal(
			t,
			defaultUser.Location.YesterdayFullAt,
			rsData.Location.YesterdayFullAt,
		)
		assert.Equal(t, defaultUser.Location.UserID, rsData.Location.UserID)
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
	tReq.AddCookie(fixtures.Tokens.ManagerAccessToken)

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

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
	tReq.AddCookie(fixtures.Tokens.ManagerAccessToken)

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestGetUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	userID := testEnv.createLocations(1)[0].UserID

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/users/%s",
		userID,
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized, nil, nil))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	mt.Do(t)
}

func TestGetPaginatedManagerUsersDefaultPage(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	testEnv.createManagerUsers(20)

	amount, err := testApp.services.Users.GetTotalCount(context.Background())
	require.Nil(t, err)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/users")
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	rs := tReq.Do(t)

	var rsData dtos.PaginatedUsersDto
	err = httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 1, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(*amount)/4),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 4, len(rsData.Data))

	assert.Equal(t, fixtures.ManagerUser.ID, rsData.Data[0].ID)
	assert.Equal(t, fixtures.ManagerUser.Username, rsData.Data[0].Username)
	assert.Equal(t, 0, len(rsData.Data[0].PasswordHash))
	assert.Equal(t, fixtures.ManagerUser.Role, rsData.Data[0].Role)
	assert.Nil(t, rsData.Data[0].Location)
}

func TestGetPaginatedManagerUsersSpecificPage(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	users := testEnv.createManagerUsers(20)

	amount, err := testApp.services.Users.GetTotalCount(context.Background())
	require.Nil(t, err)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/users")
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	rs := tReq.Do(t)

	var rsData dtos.PaginatedUsersDto
	err = httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 2, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(*amount)/4),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 4, len(rsData.Data))

	assert.Equal(t, users[11].ID, rsData.Data[0].ID)
	assert.Equal(t, users[11].Username, rsData.Data[0].Username)
	assert.Equal(t, 0, len(rsData.Data[0].PasswordHash))
	assert.Equal(t, users[11].Role, rsData.Data[0].Role)
	assert.Nil(t, rsData.Data[0].Location)
}

func TestGetPaginatedManagerUsersPageFull(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	testEnv.createManagerUsers(20)

	amount, err := testApp.services.Users.GetTotalCount(context.Background())
	require.Nil(t, err)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/users")
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	test.PaginatedEndpointTester(
		t,
		tReq,
		"page",
		int(math.Ceil(float64(*amount)/4)),
	)
}

func TestGetPaginatedManagerUsersAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/users")

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized, nil, nil))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(fixtures.Tokens.ManagerAccessToken)

	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	mt.Do(t)
}

func TestCreateManagerUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/users")
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	//nolint:exhaustruct //other fields are optional
	data := dtos.CreateUserDto{
		Username: "test",
		Password: "testpassword",
	}
	tReq.SetBody(data)

	rs := tReq.Do(t)

	var rsData models.User
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

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
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	//nolint:exhaustruct //other fields are optional
	data := dtos.CreateUserDto{
		Username: fixtures.ManagerUser.Username,
		Password: "testpassword",
	}
	tReq.SetBody(data)

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

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
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)
	//nolint:exhaustruct //other fields are optional
	tReq.SetBody(dtos.CreateUserDto{
		Username: "",
		Password: "",
	})

	mt := test.CreateMatrixTester()

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity, nil,
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"username": "must be provided",
			"password": "must be provided",
		}))

	mt.AddTestCase(tReq, tRes)

	mt.Do(t)
}

func TestCreateManagerUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/users")

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized, nil, nil))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(fixtures.Tokens.ManagerAccessToken)

	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	mt.Do(t)
}

func TestUpdateManagerUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	user := testEnv.createManagerUsers(1)[0]

	username, password := "test", "testpassword"
	//nolint:exhaustruct //other fields are optional
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/users/%s",
		user.ID,
	)
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	tReq.SetBody(data)

	rs := tReq.Do(t)

	var rsData models.User
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, user.ID, rsData.ID)
	assert.Equal(t, "test", rsData.Username)
	assert.Equal(t, models.ManagerRole, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
	assert.Nil(t, rsData.Location)
}

func TestUpdateManagerUserUserNameExists(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	user := testEnv.createManagerUsers(1)[0]

	username, password := fixtures.ManagerUser.Username, "testpassword"
	//nolint:exhaustruct //other fields are optional
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/users/%s",
		user.ID,
	)
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	tReq.SetBody(data)

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

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
	//nolint:exhaustruct //other fields are optional
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
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	tReq.SetBody(data)

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

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
	//nolint:exhaustruct //other fields are optional
	data := dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPatch, "/users/8000")
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	tReq.SetBody(data)

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestUpdateManagerUserFailValidation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	user := testEnv.createManagerUsers(1)[0]

	username, password := "", ""

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/users/%s",
		user.ID,
	)
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)
	//nolint:exhaustruct //other fields are optional
	tReq.SetBody(dtos.UpdateUserDto{
		Username: &username,
		Password: &password,
	})

	mt := test.CreateMatrixTester()

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity, nil,
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"username": "must be provided",
			"password": "must be provided",
		}))

	mt.AddTestCase(tReq, tRes)

	mt.Do(t)
}

func TestUpdateManagerUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	user := testEnv.createManagerUsers(1)[0]

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/users/%s",
		user.ID,
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized, nil, nil))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(fixtures.Tokens.ManagerAccessToken)

	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	mt.Do(t)
}

func TestDeleteManagerUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	user := testEnv.createManagerUsers(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/users/%s",
		user.ID,
	)
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	rs := tReq.Do(t)

	var rsData models.User
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, user.ID, rsData.ID)
	assert.Equal(t, user.Username, rsData.Username)
	assert.Equal(t, user.Role, rsData.Role)
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
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

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
	tReq.AddCookie(fixtures.Tokens.AdminAccessToken)

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestDeleteManagerUserAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	user := testEnv.createManagerUsers(1)[0]

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/users/%s",
		user.ID,
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized, nil, nil))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(fixtures.Tokens.ManagerAccessToken)

	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	mt.Do(t)
}
