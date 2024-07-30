package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/test"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestSignInUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "testpassword",
		RememberMe: true,
	}
	tReq.SetBody(data)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.Equal(t, 2, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken")

	assert.Equal(t, testEnv.Fixtures.DefaultUser.ID, rsData.ID)
	assert.Equal(t, testEnv.Fixtures.DefaultUser.Username, rsData.Username)
	assert.Equal(t, testEnv.Fixtures.DefaultUser.Role, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
}

func TestSignInUserNoRefresh(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "testpassword",
		RememberMe: false,
	}
	tReq.SetBody(data)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.Equal(t, 1, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")

	assert.Equal(t, testEnv.Fixtures.DefaultUser.ID, rsData.ID)
	assert.Equal(t, testEnv.Fixtures.DefaultUser.Username, rsData.Username)
	assert.Equal(t, testEnv.Fixtures.DefaultUser.Role, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
}

func TestSignInAdmin(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Admin",
		Password:   "testpassword",
		RememberMe: true,
	}
	tReq.SetBody(data)

	var rsData models.User
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.Equal(t, 1, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")

	assert.Equal(t, testEnv.Fixtures.AdminUser.ID, rsData.ID)
	assert.Equal(t, testEnv.Fixtures.AdminUser.Username, rsData.Username)
	assert.Equal(t, testEnv.Fixtures.AdminUser.Role, rsData.Role)
	assert.Equal(t, 0, len(rsData.PasswordHash))
}

func TestSignInInexistentUser(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "inexistentuser",
		Password:   "testpassword",
		RememberMe: true,
	}
	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
	assert.Equal(t, "Invalid Credentials", rsData.Message)
}

func TestSignInWrongPassword(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")

	data := dtos.SignInDto{
		Username:   "Default",
		Password:   "wrongpassword",
		RememberMe: true,
	}
	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
	assert.Equal(t, "Invalid Credentials", rsData.Message)
}

func TestSignInFailValidation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/auth/signin")
	tReq.SetBody(dtos.SignInDto{
		Username:   "",
		Password:   "",
		RememberMe: true,
	})

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes.SetBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"username": "must be provided",
			"password": "must be provided",
		}),
	)

	vt := test.CreateMatrixTester()
	vt.AddTestCase(tReq, tRes)
	vt.Do(t)
}

func TestSignOut(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/signout")

	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultRefreshToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken=;")
}

func TestSignOutNoRefresh(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/signout")

	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, 1, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
}

func TestSignOutNotLoggedIn(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/signout")

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
}

func TestRefresh(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/refresh")

	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultRefreshToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken")
}

func TestRefreshReusedToken(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/refresh")

	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultRefreshToken)

	rs1 := tReq.Do(t, nil)
	rs2 := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs1.StatusCode)
	assert.Equal(t, http.StatusUnauthorized, rs2.StatusCode)
}

func TestRefreshInvalidToken(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/refresh")

	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
}
