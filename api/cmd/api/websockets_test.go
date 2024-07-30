package main

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	tWeb.SetParallelOperation(func(t *testing.T, _ *httptest.Server) {
		school, err := testEnv.services.Schools.GetByID(context.Background(), int64(1))
		require.Nil(t, err)

		_, err = testEnv.services.CheckIns.Create(
			context.Background(),
			testEnv.Fixtures.DefaultLocation,
			school,
		)
		require.Nil(t, err)
	})

	var locationStatesInitial []dtos.LocationStateDto
	var locationStatesFinal dtos.LocationStateDto
	err := tWeb.Do(t, &locationStatesInitial, &locationStatesFinal)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(locationStatesInitial))
	assert.Equal(
		t,
		locationStatesFinal.Capacity-1,
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
	tWeb.SetParallelOperation(func(t *testing.T, _ *httptest.Server) {
		newCap := int64(10)
		err := testEnv.services.Locations.Update(
			context.Background(),
			testEnv.Fixtures.DefaultLocation,
			testEnv.Fixtures.DefaultUser,
			//nolint:exhaustruct //other fields are optional
			dtos.UpdateLocationDto{
				Capacity: &newCap,
			},
		)
		require.Nil(t, err)
	})

	var locationStatesInitial []dtos.LocationStateDto
	var locationStatesFinal dtos.LocationStateDto
	err := tWeb.Do(t, &locationStatesInitial, &locationStatesFinal)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(locationStatesInitial))
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

	tWeb.SetParallelOperation(func(t *testing.T, _ *httptest.Server) {
		school, err := testEnv.services.Schools.GetByID(context.Background(), int64(1))
		require.Nil(t, err)

		_, err = testEnv.services.CheckIns.Create(
			context.Background(),
			testEnv.Fixtures.DefaultLocation,
			school,
		)
		require.Nil(t, err)
	})

	var locationState dtos.LocationStateDto
	err := tWeb.Do(t, nil, &locationState)

	assert.Nil(t, err)
	assert.Equal(
		t,
		locationState.NormalizedName,
		testEnv.Fixtures.DefaultLocation.NormalizedName,
	)
	assert.Equal(t, locationState.Capacity-1, locationState.Available)
}

func TestSingleLocationWebSocketCapUpdate(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tWeb := test.CreateWebSocketTester(testApp.routes())

	tWeb.SetInitialMessage(dtos.SubscribeMessageDto{
		Subject:        "single-location",
		NormalizedName: testEnv.Fixtures.DefaultLocation.NormalizedName,
	})

	tWeb.SetParallelOperation(func(t *testing.T, _ *httptest.Server) {
		newCap := int64(10)
		err := testEnv.services.Locations.Update(
			context.Background(),
			testEnv.Fixtures.DefaultLocation,
			testEnv.Fixtures.DefaultUser,
			//nolint:exhaustruct //other fields are optional
			dtos.UpdateLocationDto{
				Capacity: &newCap,
			},
		)
		require.Nil(t, err)
	})

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
