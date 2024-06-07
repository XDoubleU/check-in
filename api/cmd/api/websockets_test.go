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

	tWeb := test.CreateTestWebsocket(testApp.routes())

	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject: "all-locations",
	})

	tWeb.SetParallelOperation(createCheckIn)

	var locationUpdateEventsInitial []models.LocationUpdateEvent
	var locationUpdateEventsFinal []models.LocationUpdateEvent
	tWeb.Do(t, &locationUpdateEventsInitial, &locationUpdateEventsFinal)

	assert.Equal(t, len(locationUpdateEventsInitial), 21)
	assert.Equal(t, len(locationUpdateEventsFinal), 1)
	assert.Equal(
		t,
		locationUpdateEventsFinal[0].Available,
		locationUpdateEventsFinal[0].Capacity-6,
	)
}

func TestAllLocationsWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tWeb := test.CreateTestWebsocket(testApp.routes())
	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject: "all-locations",
	})
	tWeb.SetParallelOperation(updateCapacity)

	var locationUpdateEventsInitial []models.LocationUpdateEvent
	var locationUpdateEventsFinal []models.LocationUpdateEvent
	tWeb.Do(t, &locationUpdateEventsInitial, &locationUpdateEventsFinal)

	assert.Equal(t, len(locationUpdateEventsInitial), 21)
	assert.Equal(t, len(locationUpdateEventsFinal), 1)
	assert.EqualValues(t, locationUpdateEventsFinal[0].Capacity, 10)
}

func TestSingleLocationWebSocketCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tWeb := test.CreateTestWebsocket(testApp.routes())

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
	assert.Equal(t, locationUpdateEvent.Available, locationUpdateEvent.Capacity-6)
}

func TestSingleLocationWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tWeb := test.CreateTestWebsocket(testApp.routes())

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
	assert.EqualValues(t, locationUpdateEvent.Capacity, 10)
}

func createCheckIn(t *testing.T, ts *httptest.Server) {
	data := dtos.CreateCheckInDto{
		SchoolID: fixtureData.Schools[0].ID,
	}

	tReq := test.CreateTestRequest(nil, http.MethodPost, "/checkins")
	tReq.SetTestServer(ts)
	tReq.SetReqData(data)
	tReq.AddCookie(tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, rs.StatusCode, http.StatusCreated)
}

func updateCapacity(t *testing.T, ts *httptest.Server) {
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Capacity: &capacity,
	}

	tReq := test.CreateTestRequest(nil, http.MethodPatch, "/locations/%s", fixtureData.DefaultLocation.ID)
	tReq.SetTestServer(ts)
	tReq.SetReqData(data)
	tReq.AddCookie(tokens.DefaultAccessToken)

	rs := tReq.Do(t, nil)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
}
