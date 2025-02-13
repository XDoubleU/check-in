package main

import (
	"net/http"
	"testing"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestGetState(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/state",
	)
	rs := tReq.Do(t)

	var rsData models.State
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, false, rsData.IsMaintenance)
	assert.Equal(t, true, rsData.IsDatabaseActive)
}

func TestUpdateState(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	data := dtos.StateDto{
		IsMaintenance: true,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/state",
	)
	tReq.AddCookie(testEnv.fixtures.Tokens.AdminAccessToken)

	tReq.SetData(data)

	rs := tReq.Do(t)

	var rsData models.State
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, true, rsData.IsMaintenance)
	assert.Equal(t, true, rsData.IsDatabaseActive)
}

func TestUpdateStateAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/state",
	)

	mt := test.CreateMatrixTester()

	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized, nil, nil))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	tReq3 := tReqBase.Copy()
	tReq3.AddCookie(testEnv.fixtures.Tokens.ManagerAccessToken)

	mt.AddTestCase(tReq3, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	mt.Do(t)
}
