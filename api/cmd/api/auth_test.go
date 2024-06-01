package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"

	"github.com/stretchr/testify/assert"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/XDoubleU/essentia/pkg/test"
)

func TestSignInUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "testpassword",
		RememberMe: true,
	}
	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(&rsData)

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

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "testpassword",
		RememberMe: false,
	}
	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(&rsData)

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

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Admin",
		Password:   "testpassword",
		RememberMe: true,
	}
	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(&rsData)

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

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "inexistentuser",
		Password:   "testpassword",
		RememberMe: true,
	}
	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rsData.Message, "Invalid Credentials")
}

func TestSignInWrongPassword(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "wrongpassword",
		RememberMe: true,
	}
	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(&rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rsData.Message, "Invalid Credentials")
}

func TestSignInFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodPost, "/auth/signin")

	tReq.SetReqData(dtos.SignInDto{
		Username:   "",
		Password:   "",
		RememberMe: true,
	})

	errorMessage := map[string]interface{}{
		"username": "must be provided",
		"password": "must be provided",
	}

	vt := test.CreateValidatorTester(t)
	vt.AddTestCase(tReq, errorMessage)
	vt.Do()
}

func TestSignOut(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/auth/signout")

	tReq.AddCookie(tokens.DefaultAccessToken)
	tReq.AddCookie(tokens.DefaultRefreshToken)

	rs := tReq.Do(nil)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken=;")
}

func TestSignOutNoRefresh(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/auth/signout")

	tReq.AddCookie(tokens.DefaultAccessToken)

	rs := tReq.Do(nil)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, len(rs.Header.Values("set-cookie")), 1)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
}

func TestSignOutNotLoggedIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/auth/signout")

	rs := tReq.Do(nil)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
}

func TestRefresh(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/auth/refresh")

	tReq.AddCookie(tokens.DefaultRefreshToken)

	rs := tReq.Do(nil)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken")
}

func TestRefreshReusedToken(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/auth/refresh")

	tReq.AddCookie(tokens.DefaultRefreshToken)

	rs1 := tReq.Do(nil)
	rs2 := tReq.Do(nil)

	assert.Equal(t, rs1.StatusCode, http.StatusOK)
	assert.Equal(t, rs2.StatusCode, http.StatusUnauthorized)
}

func TestRefreshInvalidToken(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq := test.CreateTestRequest(t, ts, http.MethodGet, "/auth/refresh")

	tReq.AddCookie(tokens.DefaultAccessToken)

	rs := tReq.Do(nil)

	assert.Equal(t, rs.StatusCode, http.StatusUnauthorized)
}
