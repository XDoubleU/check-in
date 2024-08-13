package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
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

	rs := tReq.Do(t)

	var rsData models.User
	httptools.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.Equal(t, 2, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken")

	assert.Equal(t, fixtures.DefaultUser.ID, rsData.ID)
	assert.Equal(t, fixtures.DefaultUser.Username, rsData.Username)
	assert.Equal(t, fixtures.DefaultUser.Role, rsData.Role)
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

	rs := tReq.Do(t)

	var rsData models.User
	httptools.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.Equal(t, 1, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")

	assert.Equal(t, fixtures.DefaultUser.ID, rsData.ID)
	assert.Equal(t, fixtures.DefaultUser.Username, rsData.Username)
	assert.Equal(t, fixtures.DefaultUser.Role, rsData.Role)
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

	rs := tReq.Do(t)

	var rsData models.User
	httptools.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.Equal(t, 1, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")

	assert.Equal(t, fixtures.AdminUser.ID, rsData.ID)
	assert.Equal(t, fixtures.AdminUser.Username, rsData.Username)
	assert.Equal(t, fixtures.AdminUser.Role, rsData.Role)
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

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	httptools.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
	assert.Equal(t, "invalid credentials", rsData.Message)
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

	rs := tReq.Do(t)

	var rsData errortools.ErrorDto
	httptools.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
	assert.Equal(t, "invalid credentials", rsData.Message)
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

	tRes := test.NewCaseResponse(http.StatusUnprocessableEntity, nil,
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"username": "must be provided",
			"password": "must be provided",
		}))

	vt := test.CreateMatrixTester()
	vt.AddTestCase(tReq, tRes)
	vt.Do(t)
}

func TestSignOut(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/signout")

	tReq.AddCookie(fixtures.Tokens.DefaultAccessToken)
	tReq.AddCookie(fixtures.Tokens.DefaultRefreshToken)

	rs := tReq.Do(t)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken=;")
}

func TestSignOutNoRefresh(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/signout")

	tReq.AddCookie(fixtures.Tokens.DefaultAccessToken)

	rs := tReq.Do(t)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, 1, len(rs.Header.Values("set-cookie")))
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken=;")
}

func TestSignOutNotLoggedIn(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/signout")

	rs := tReq.Do(t)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
}

func TestRefresh(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/refresh")

	tReq.AddCookie(fixtures.Tokens.DefaultRefreshToken)

	rs := tReq.Do(t)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Contains(t, rs.Header.Values("set-cookie")[0], "accessToken")
	assert.Contains(t, rs.Header.Values("set-cookie")[1], "refreshToken")
}

func TestRefreshReusedToken(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/refresh")

	tReq.AddCookie(fixtures.Tokens.DefaultRefreshToken)

	rs1 := tReq.Do(t)
	rs2 := tReq.Do(t)

	assert.Equal(t, http.StatusOK, rs1.StatusCode)
	assert.Equal(t, http.StatusUnauthorized, rs2.StatusCode)
}

func TestRefreshInvalidToken(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/auth/refresh")

	tReq.AddCookie(fixtures.Tokens.DefaultAccessToken)

	rs := tReq.Do(t)

	assert.Equal(t, http.StatusUnauthorized, rs.StatusCode)
}
