package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestAllLocationsWebSocketCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tWeb := test.CreateWebsocketTester(testApp.routes())

	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject: "all-locations",
	})

	tWeb.SetParallelOperation(createCheckIn)

	var locationUpdateEventsInitial []models.LocationUpdateEvent
	var locationUpdateEventsFinal []models.LocationUpdateEvent
	tWeb.Do(t, &locationUpdateEventsInitial, &locationUpdateEventsFinal)

	assert.Equal(t, 21, len(locationUpdateEventsInitial))
	assert.Equal(t, 1, len(locationUpdateEventsFinal))
	assert.Equal(
		t,
		locationUpdateEventsFinal[0].Capacity-6,
		locationUpdateEventsFinal[0].Available,
	)
}

func TestAllLocationsWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tWeb := test.CreateWebsocketTester(testApp.routes())
	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject: "all-locations",
	})
	tWeb.SetParallelOperation(updateCapacity)

	var locationUpdateEventsInitial []models.LocationUpdateEvent
	var locationUpdateEventsFinal []models.LocationUpdateEvent
	tWeb.Do(t, &locationUpdateEventsInitial, &locationUpdateEventsFinal)

	assert.Equal(t, 21, len(locationUpdateEventsInitial))
	assert.Equal(t, 1, len(locationUpdateEventsFinal))
	assert.EqualValues(t, 10, locationUpdateEventsFinal[0].Capacity)
}

func TestSingleLocationWebSocketCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tWeb := test.CreateWebsocketTester(testApp.routes())

	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject:        "single-location",
		NormalizedName: fixtureData.DefaultLocation.NormalizedName,
	})

	tWeb.SetParallelOperation(createCheckIn)

	var locationUpdateEvent models.LocationUpdateEvent
	tWeb.Do(t, nil, &locationUpdateEvent)

	assert.Equal(
		t,
		locationUpdateEvent.NormalizedName,
		fixtureData.DefaultLocation.NormalizedName,
	)
	assert.Equal(t, locationUpdateEvent.Capacity-6, locationUpdateEvent.Available)
}

func TestSingleLocationWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tWeb := test.CreateWebsocketTester(testApp.routes())

	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject:        "single-location",
		NormalizedName: fixtureData.DefaultLocation.NormalizedName,
	})

	tWeb.SetParallelOperation(updateCapacity)

	var locationUpdateEvent models.LocationUpdateEvent
	tWeb.Do(t, nil, &locationUpdateEvent)

	assert.Equal(
		t,
		locationUpdateEvent.NormalizedName,
		fixtureData.DefaultLocation.NormalizedName,
	)
	assert.EqualValues(t, 10, locationUpdateEvent.Capacity)
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
