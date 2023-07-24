package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nhooyr.io/websocket/wsjson"

	"check-in/api/internal/assert"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"
)

const timeout = time.Minute
const sleep = time.Second

func TestAllLocationsWebSocketCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewServer(testApp.routes())
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	ws := tests.DialWebsocket(wsURL, timeout)

	msg := dtos.SubscribeMessageDto{
		Subject: "all-locations",
	}

	err := wsjson.Write(context.Background(), ws, msg)
	if err != nil {
		panic(err)
	}

	var locationUpdateEvents []models.LocationUpdateEvent
	err = wsjson.Read(context.Background(), ws, &locationUpdateEvents)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(locationUpdateEvents), 21)

	go func() {
		time.Sleep(sleep)

		data := dtos.CheckInDto{
			SchoolID: fixtureData.Schools[0].ID,
		}

		body, _ := json.Marshal(data)
		req, _ := http.NewRequest(
			http.MethodPost,
			ts.URL+"/checkins",
			bytes.NewReader(body),
		)
		req.AddCookie(&tokens.DefaultAccessToken)

		rs, _ := ts.Client().Do(req)

		assert.Equal(t, rs.StatusCode, http.StatusCreated)
	}()

	err = wsjson.Read(context.Background(), ws, &locationUpdateEvents)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(locationUpdateEvents), 1)
	assert.Equal(
		t,
		locationUpdateEvents[0].Available,
		locationUpdateEvents[0].Capacity-6,
	)
}

func TestAllLocationsWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewServer(testApp.routes())
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	ws := tests.DialWebsocket(wsURL, timeout)

	msg := dtos.SubscribeMessageDto{
		Subject: "all-locations",
	}

	err := wsjson.Write(context.Background(), ws, msg)
	if err != nil {
		panic(err)
	}

	var locationUpdateEvents []models.LocationUpdateEvent
	err = wsjson.Read(context.Background(), ws, &locationUpdateEvents)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(locationUpdateEvents), 21)

	go func() {
		time.Sleep(sleep)

		var capacity int64 = 10
		data := dtos.UpdateLocationDto{
			Capacity: &capacity,
		}

		body, _ := json.Marshal(data)
		req, _ := http.NewRequest(http.MethodPatch,
			ts.URL+"/locations/"+fixtureData.DefaultLocation.ID, bytes.NewReader(body))
		req.AddCookie(&tokens.DefaultAccessToken)

		rs, _ := ts.Client().Do(req)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
	}()

	err = wsjson.Read(context.Background(), ws, &locationUpdateEvents)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(locationUpdateEvents), 1)
	assert.Equal(t, locationUpdateEvents[0].Capacity, 10)
}

func TestSingleLocationWebSocketCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewServer(testApp.routes())
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	ws := tests.DialWebsocket(wsURL, timeout)

	msg := dtos.SubscribeMessageDto{
		Subject:        "single-location",
		NormalizedName: fixtureData.DefaultLocation.NormalizedName,
	}

	err := wsjson.Write(context.Background(), ws, msg)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(sleep)

		data := dtos.CheckInDto{
			SchoolID: fixtureData.Schools[0].ID,
		}

		body, _ := json.Marshal(data)
		req, _ := http.NewRequest(http.MethodPost,
			ts.URL+"/checkins", bytes.NewReader(body))
		req.AddCookie(&tokens.DefaultAccessToken)

		rs, _ := ts.Client().Do(req)

		assert.Equal(t, rs.StatusCode, http.StatusCreated)
	}()

	var locationUpdateEvent models.LocationUpdateEvent
	err = wsjson.Read(context.Background(), ws, &locationUpdateEvent)
	if err != nil {
		panic(err)
	}

	assert.Equal(
		t,
		locationUpdateEvent.NormalizedName,
		fixtureData.DefaultLocation.NormalizedName,
	)
	assert.Equal(t, locationUpdateEvent.Available, locationUpdateEvent.Capacity-6)
}

func TestSingleLocationWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewServer(testApp.routes())
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	ws := tests.DialWebsocket(wsURL, timeout)

	msg := dtos.SubscribeMessageDto{
		Subject:        "single-location",
		NormalizedName: fixtureData.DefaultLocation.NormalizedName,
	}

	err := wsjson.Write(context.Background(), ws, msg)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(sleep)

		var capacity int64 = 10
		data := dtos.UpdateLocationDto{
			Capacity: &capacity,
		}

		body, _ := json.Marshal(data)
		req, _ := http.NewRequest(http.MethodPatch,
			ts.URL+"/locations/"+fixtureData.DefaultLocation.ID, bytes.NewReader(body))
		req.AddCookie(&tokens.DefaultAccessToken)

		rs, _ := ts.Client().Do(req)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
	}()

	var locationUpdateEvent models.LocationUpdateEvent
	err = wsjson.Read(context.Background(), ws, &locationUpdateEvent)
	if err != nil {
		panic(err)
	}

	assert.Equal(
		t,
		locationUpdateEvent.NormalizedName,
		fixtureData.DefaultLocation.NormalizedName,
	)
	assert.Equal(t, locationUpdateEvent.Capacity, 10)
}
