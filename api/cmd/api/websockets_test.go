package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/test"

	"check-in/api/internal/dtos"
)

func TestAllLocationsWebSocketCheckIn(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tWeb := test.CreateWebSocketTester(testApp.routes())

	//nolint:exhaustruct // other fields are optional
	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject: "all-locations",
	})

	tWeb.SetParallelOperation(func(t *testing.T, ts *httptest.Server) { createCheckIn(t, ts, testEnv) })

	var locationStatesInitial []dtos.LocationStateDto
	var locationStatesFinal dtos.LocationStateDto
	err := tWeb.Do(t, &locationStatesInitial, &locationStatesFinal)

	assert.Nil(t, err)
	assert.Equal(t, 21, len(locationStatesInitial))
	assert.Equal(
		t,
		locationStatesFinal.Capacity-6,
		locationStatesFinal.Available,
	)
}

func TestAllLocationsWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tWeb := test.CreateWebSocketTester(testApp.routes())
	//nolint:exhaustruct // other fields are optional
	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject: "all-locations",
	})
	tWeb.SetParallelOperation(func(t *testing.T, ts *httptest.Server) { updateCapacity(t, ts, testEnv) })

	var locationStatesInitial []dtos.LocationStateDto
	var locationStatesFinal dtos.LocationStateDto
	err := tWeb.Do(t, &locationStatesInitial, &locationStatesFinal)

	assert.Nil(t, err)
	assert.Equal(t, 21, len(locationStatesInitial))
	assert.EqualValues(t, 10, locationStatesFinal.Capacity)
}

func TestSingleLocationWebSocketCheckIn(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tWeb := test.CreateWebSocketTester(testApp.routes())

	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject:        "single-location",
		NormalizedName: testEnv.Fixtures.DefaultLocation.NormalizedName,
	})

	tWeb.SetParallelOperation(func(t *testing.T, ts *httptest.Server) { createCheckIn(t, ts, testEnv) })

	var locationState dtos.LocationStateDto
	err := tWeb.Do(t, nil, &locationState)

	assert.Nil(t, err)
	assert.Equal(
		t,
		locationState.NormalizedName,
		testEnv.Fixtures.DefaultLocation.NormalizedName,
	)
	assert.Equal(t, locationState.Capacity-6, locationState.Available)
}

func TestSingleLocationWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tWeb := test.CreateWebSocketTester(testApp.routes())

	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject:        "single-location",
		NormalizedName: testEnv.Fixtures.DefaultLocation.NormalizedName,
	})

	tWeb.SetParallelOperation(func(t *testing.T, ts *httptest.Server) { updateCapacity(t, ts, testEnv) })

	var locationState dtos.LocationStateDto
	err := tWeb.Do(t, nil, &locationState)

	assert.Nil(t, err)
	assert.Equal(
		t,
		locationState.NormalizedName,
		testEnv.Fixtures.DefaultLocation.NormalizedName,
	)
	assert.EqualValues(t, 10, locationState.Capacity)
}

func createCheckIn(t *testing.T, ts *httptest.Server, testEnv TestEnv) {
	data := dtos.CreateCheckInDto{
		SchoolID: testEnv.Fixtures.Schools[0].ID,
	}

	tReq := test.CreateRequestTester(nil, http.MethodPost, "/checkins")
	tReq.SetTestServer(ts)
	tReq.SetReqData(data)
	tReq.AddCookie(testEnv.Tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
}

func updateCapacity(t *testing.T, ts *httptest.Server, testEnv TestEnv) {
	var capacity int64 = 10
	//nolint:exhaustruct // other fields are optional
	data := dtos.UpdateLocationDto{
		Capacity: &capacity,
	}

	tReq := test.CreateRequestTester(
		nil,
		http.MethodPatch,
		"/locations/%s",
		testEnv.Fixtures.DefaultLocation.ID,
	)
	tReq.SetTestServer(ts)
	tReq.SetReqData(data)
	tReq.AddCookie(testEnv.Tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
