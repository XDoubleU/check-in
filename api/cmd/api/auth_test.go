package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/httptools"
	"github.com/xdoubleu/essentia/pkg/test"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestSignInUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "testpassword",
		RememberMe: true,
	}
	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.Equal(t, 2, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken")

	assert.Equal(t, fixtureData.DefaultUser.ID, rsData.ID)
	assert.Equal(t, fixtureData.DefaultUser.Username, rsData.Username)
	assert.Equal(t, fixtureData.DefaultUser.Role, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
}

func TestSignInUserNoRefresh(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "testpassword",
		RememberMe: false,
	}
	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.Equal(t, 1, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")

	assert.Equal(t, fixtureData.DefaultUser.ID, rsData.ID)
	assert.Equal(t, fixtureData.DefaultUser.Username, rsData.Username)
	assert.Equal(t, fixtureData.DefaultUser.Role, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
}

func TestSignInAdmin(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Admin",
		Password:   "testpassword",
		RememberMe: true,
	}
	tReq.SetReqData(data)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.Equal(t, 1, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")

	assert.Equal(t, fixtureData.AdminUser.ID, rsData.ID)
	assert.Equal(t, fixtureData.AdminUser.Username, rsData.Username)
	assert.Equal(t, fixtureData.AdminUser.Role, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
}

func TestSignInInexistentUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "inexistentuser",
		Password:   "testpassword",
		RememberMe: true,
	}
	tReq.SetReqData(data)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
	assert.Equal(t, "Invalid Credentials", rsData.Message)
}

func TestSignInWrongPassword(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "wrongpassword",
		RememberMe: true,
	}
	tReq.SetReqData(data)

	var rsData httptools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
	assert.Equal(t, "Invalid Credentials", rsData.Message)
}

func TestSignInFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")
	tReq.SetReqData(dtos.SignInDto{
		Username:   "",
		Password:   "",
		RememberMe: true,
	})

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes.SetExpectedBody(
		httptools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"username": "must be provided",
			"password": "must be provided",
		}),
	)

	vt := test.CreateMatrixTester()
	vt.AddTestCase(tReq, tRes)
	vt.Do(t)
}

func TestSignOut(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/signout")

	tReq.AddCookie(tokens.DefaultAccessToken)
	tReq.AddCookie(tokens.DefaultRefreshToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken=;")
}

func TestSignOutNoRefresh(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/signout")

	tReq.AddCookie(tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, 1, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
}

func TestSignOutNotLoggedIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/signout")

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
}

func TestRefresh(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/refresh")

	tReq.AddCookie(tokens.DefaultRefreshToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken")
}

func TestRefreshReusedToken(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/refresh")

	tReq.AddCookie(tokens.DefaultRefreshToken)

	rs1 := tReq.Do(t, nil)
	rs2 := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs1.StatusCode)
	assert.Equal(t, http.StatusUnauthorized, rs2.StatusCode)
}

func TestRefreshInvalidToken(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/refresh")

	tReq.AddCookie(tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
}
