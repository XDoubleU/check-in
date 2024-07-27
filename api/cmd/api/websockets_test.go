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
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tWeb := test.CreateWebSocketTester(testApp.routes())

	//nolint:exhaustruct // other fields are optional
	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject: "all-locations",
	})

	tWeb.SetParallelOperation(createCheckIn)

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
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tWeb := test.CreateWebSocketTester(testApp.routes())
	//nolint:exhaustruct // other fields are optional
	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject: "all-locations",
	})
	tWeb.SetParallelOperation(updateCapacity)

	var locationStatesInitial []dtos.LocationStateDto
	var locationStatesFinal dtos.LocationStateDto
	err := tWeb.Do(t, &locationStatesInitial, &locationStatesFinal)

	assert.Nil(t, err)
	assert.Equal(t, 21, len(locationStatesInitial))
	assert.EqualValues(t, 10, locationStatesFinal.Capacity)
}

func TestSingleLocationWebSocketCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tWeb := test.CreateWebSocketTester(testApp.routes())

	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject:        "single-location",
		NormalizedName: fixtureData.DefaultLocation.NormalizedName,
	})

	tWeb.SetParallelOperation(createCheckIn)

	var locationState dtos.LocationStateDto
	err := tWeb.Do(t, nil, &locationState)

	assert.Nil(t, err)
	assert.Equal(
		t,
		locationState.NormalizedName,
		fixtureData.DefaultLocation.NormalizedName,
	)
	assert.Equal(t, locationState.Capacity-6, locationState.Available)
}

func TestSingleLocationWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer testEnv.TeardownSingle()

	tWeb := test.CreateWebSocketTester(testApp.routes())

	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject:        "single-location",
		NormalizedName: fixtureData.DefaultLocation.NormalizedName,
	})

	tWeb.SetParallelOperation(updateCapacity)

	var locationState dtos.LocationStateDto
	err := tWeb.Do(t, nil, &locationState)

	assert.Nil(t, err)
	assert.Equal(
		t,
		locationState.NormalizedName,
		fixtureData.DefaultLocation.NormalizedName,
	)
	assert.EqualValues(t, 10, locationState.Capacity)
}

func createCheckIn(t *testing.T, ts *httptest.Server) {
	data := dtos.CreateCheckInDto{
		SchoolID: fixtureData.Schools[0].ID,
	}

	tReq := test.CreateRequestTester(nil, http.MethodPost, "/checkins")
	tReq.SetTestServer(ts)
	tReq.SetReqData(data)
	tReq.AddCookie(tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
}

func updateCapacity(t *testing.T, ts *httptest.Server) {
	var capacity int64 = 10
	//nolint:exhaustruct // other fields are optional
	data := dtos.UpdateLocationDto{
		Capacity: &capacity,
	}

	tReq := test.CreateRequestTester(
		nil,
		http.MethodPatch,
		"/locations/%s",
		fixtureData.DefaultLocation.ID,
	)
	tReq.SetTestServer(ts)
	tReq.SetReqData(data)
	tReq.AddCookie(tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
