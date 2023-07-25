package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"check-in/api/internal/assert"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"
)

func TestSignInUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "testpassword",
		RememberMe: true,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/auth/signin",
		bytes.NewReader(body),
	)
	rs, _ := ts.Client().Do(req)

	var rsData models.User
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.Equal(t, len(rs.Header.Values("set-cookie")), 2)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken")

	assert.Equal(t, rsData.ID, fixtureData.DefaultUser.ID)
	assert.Equal(t, rsData.Username, fixtureData.DefaultUser.Username)
	assert.Equal(t, rsData.Role, fixtureData.DefaultUser.Role)
	assert.Equal(t, len(rsData.PasswordHash), 0)
}

func TestSignInUserNoRefresh(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "testpassword",
		RememberMe: false,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/auth/signin",
		bytes.NewReader(body),
	)
	rs, _ := ts.Client().Do(req)

	var rsData models.User
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.Equal(t, len(rs.Header.Values("set-cookie")), 1)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")

	assert.Equal(t, rsData.ID, fixtureData.DefaultUser.ID)
	assert.Equal(t, rsData.Username, fixtureData.DefaultUser.Username)
	assert.Equal(t, rsData.Role, fixtureData.DefaultUser.Role)
	assert.Equal(t, len(rsData.PasswordHash), 0)
}

func TestSignInAdmin(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SignInDto{
		Username:   "Admin",
		Password:   "testpassword",
		RememberMe: true,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/auth/signin",
		bytes.NewReader(body),
	)
	rs, _ := ts.Client().Do(req)

	var rsData models.User
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.Equal(t, len(rs.Header.Values("set-cookie")), 1)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")

	assert.Equal(t, rsData.ID, fixtureData.AdminUser.ID)
	assert.Equal(t, rsData.Username, fixtureData.AdminUser.Username)
	assert.Equal(t, rsData.Role, fixtureData.AdminUser.Role)
	assert.Equal(t, len(rsData.PasswordHash), 0)
}

func TestSignInInexistentUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SignInDto{
		Username:   "inexistentuser",
		Password:   "testpassword",
		RememberMe: true,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/auth/signin",
		bytes.NewReader(body),
	)
	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rsData.Message, "Invalid Credentials")
}

func TestSignInWrongPassword(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "wrongpassword",
		RememberMe: true,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/auth/signin",
		bytes.NewReader(body),
	)
	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rsData.Message, "Invalid Credentials")
}

func TestSignInFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.SignInDto{
		Username:   "",
		Password:   "",
		RememberMe: true,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/auth/signin",
		bytes.NewReader(body),
	)
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

func TestSignOut(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/auth/signout", nil)
	req.AddCookie(&tokens.DefaultAccessToken)
	req.AddCookie(&tokens.DefaultRefreshToken)

	rs, _ := ts.Client().Do(req)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken=;")
}

func TestSignOutNoRefresh(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/auth/signout", nil)
	req.AddCookie(&tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, len(rs.Header.Values("set-cookie")), 1)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
}

func TestSignOutNotLoggedIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/auth/signout", nil)
	rs, _ := ts.Client().Do(req)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
}

func TestRefresh(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/auth/refresh", nil)
	req.AddCookie(&tokens.DefaultRefreshToken)

	rs, _ := ts.Client().Do(req)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken")
}

func TestRefreshReusedToken(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/auth/refresh", nil)
	req.AddCookie(&tokens.DefaultRefreshToken)

	rs1, _ := ts.Client().Do(req)
	rs2, _ := ts.Client().Do(req)

	assert.Equal(t, rs1.StatusCode, http.StatusOK)
	assert.Equal(t, rs2.StatusCode, http.StatusUnauthorized)
}

func TestRefreshInvalidToken(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/auth/refresh", nil)
	req.AddCookie(&tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
}
