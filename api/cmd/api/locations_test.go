package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"check-in/api/internal/assert"
	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"
)

func TestYesterdayFullAt(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	loc, _ := time.LoadLocation("Europe/Brussels")

	now := time.Now().In(loc).AddDate(0, 0, -1)

	for i := 0; i < int(fixtureData.DefaultLocation.Capacity); i++ {
		query := `
			INSERT INTO check_ins 
			(location_id, school_id, capacity, created_at)
			VALUES ($1, $2, $3, $4)
		`

		_, err := testEnv.TestTx.Exec(
			context.Background(),
			query,
			fixtureData.DefaultLocation.ID,
			1,
			fixtureData.DefaultLocation.Capacity,
			now,
		)
		if err != nil {
			panic(err)
		}
	}

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/locations/"+fixtureData.DefaultLocation.ID,
		nil,
	)
	req.AddCookie(tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData models.Location
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.Equal(t, rsData.AvailableYesterday, 0)
	assert.Equal(t, rsData.CapacityYesterday, fixtureData.DefaultLocation.Capacity)
	assert.Equal(t, rsData.YesterdayFullAt.Valid, true)
	assert.Equal(t, rsData.YesterdayFullAt.Time.Day(), now.Day())
	assert.Equal(t, rsData.YesterdayFullAt.Time.Hour(), now.Hour())
}

func TestGetCheckInsLocationRangeRawSingle(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	startDate := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, utc)
	startDate = *helpers.StartOfDay(&startDate)

	endDate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, utc)
	endDate = *helpers.StartOfDay(&endDate)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/all-locations/checkins/range",
			nil,
		)
		req.AddCookie(user)

		query := req.URL.Query()
		query.Add("ids", fixtureData.DefaultLocation.ID)
		query.Add("startDate", startDate.Format(constants.DateFormat))
		query.Add("endDate", endDate.Format(constants.DateFormat))
		query.Add("returnType", "raw")

		req.URL.RawQuery = query.Encode()

		rs, _ := ts.Client().Do(req)

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		assert.Equal(t, rsData[startDate.Format(time.RFC3339)].Capacities.Len(), 0)

		value, present := rsData[startDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, value, 0)
		assert.Equal(t, present, true)

		capacity, _ := rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].Capacities.Get(
			fixtureData.DefaultLocation.ID,
		)
		assert.Equal(
			t,
			capacity,
			fixtureData.DefaultLocation.Capacity,
		)

		value, present = rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].
			Schools.Get(
			"Andere",
		)
		assert.Equal(t, value, 5)
		assert.Equal(t, present, true)

		assert.Equal(t, rsData[endDate.Format(time.RFC3339)].Capacities.Len(), 0)

		value, present = rsData[endDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, value, 0)
		assert.Equal(t, present, true)
	}
}

func TestGetCheckInsLocationRangeRawMultiple(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	startDate := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, utc)
	startDate = *helpers.StartOfDay(&startDate)

	endDate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, utc)
	endDate = *helpers.StartOfDay(&endDate)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/all-locations/checkins/range",
			nil,
		)
		req.AddCookie(user)

		query := req.URL.Query()
		query.Add(
			"ids",
			fmt.Sprintf(
				"%s,%s",
				fixtureData.DefaultLocation.ID,
				fixtureData.Locations[0].ID,
			),
		)
		query.Add("startDate", startDate.Format(constants.DateFormat))
		query.Add("endDate", endDate.Format(constants.DateFormat))
		query.Add("returnType", "raw")

		req.URL.RawQuery = query.Encode()

		rs, _ := ts.Client().Do(req)

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		assert.Equal(t, rsData[startDate.Format(time.RFC3339)].Capacities.Len(), 0)

		value, present := rsData[startDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, value, 0)
		assert.Equal(t, present, true)

		capacity0, _ := rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].
			Capacities.Get(
			fixtureData.DefaultLocation.ID,
		)
		capacity1, _ := rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].
			Capacities.Get(
			fixtureData.Locations[0].ID,
		)
		assert.Equal(
			t,
			capacity0,
			fixtureData.DefaultLocation.Capacity,
		)

		assert.Equal(
			t,
			capacity1,
			fixtureData.Locations[0].Capacity,
		)

		value, present = rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].
			Schools.Get(
			"Andere",
		)
		assert.Equal(t, value, 10)
		assert.Equal(t, present, true)

		assert.Equal(t, rsData[endDate.Format(time.RFC3339)].Capacities.Len(), 0)

		value, present = rsData[endDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, value, 0)
		assert.Equal(t, present, true)
	}
}

func TestGetCheckInsLocationRangeCSV(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	startDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 2).Format(constants.DateFormat)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/all-locations/checkins/range",
			nil,
		)
		req.AddCookie(user)

		query := req.URL.Query()
		query.Add("ids", fixtureData.DefaultLocation.ID)
		query.Add("startDate", startDate)
		query.Add("endDate", endDate)
		query.Add("returnType", "csv")

		req.URL.RawQuery = query.Encode()

		rs, _ := ts.Client().Do(req)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rs.Header.Get("content-type"), "text/csv")
	}
}

func TestGetCheckInsLocationRangeNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	startDate := time.Now().Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/range",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	query := req.URL.Query()
	query.Add("ids", id.String())
	query.Add("startDate", startDate)
	query.Add("endDate", endDate)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
	)
}

func TestGetCheckInsLocationRangeNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	startDate := time.Now().Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/range",
		nil,
	)
	req.AddCookie(tokens.DefaultAccessToken)

	query := req.URL.Query()
	query.Add("ids", fixtureData.Locations[0].ID)
	query.Add("startDate", startDate)
	query.Add("endDate", endDate)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
	)
}

func TestGetCheckInsLocationRangeStartDateMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/range",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	query := req.URL.Query()
	query.Add("ids", fixtureData.Locations[0].ID)
	query.Add("endDate", endDate)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "missing startDate param in query")
}

func TestGetCheckInsLocationRangeEndDateMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	startDate := time.Now().Format(constants.DateFormat)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/range",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	query := req.URL.Query()
	query.Add("ids", fixtureData.Locations[0].ID)
	query.Add("startDate", startDate)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "missing endDate param in query")
}

func TestGetCheckInsLocationRangeReturnTypeMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	startDate := time.Now().Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/range",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	query := req.URL.Query()
	query.Add("ids", fixtureData.Locations[0].ID)
	query.Add("startDate", startDate)
	query.Add("endDate", endDate)

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "missing returnType param in query")
}

func TestGetCheckInsLocationRangeNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	startDate := time.Now().Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/range",
		nil,
	)
	req.AddCookie(tokens.DefaultAccessToken)

	query := req.URL.Query()
	query.Add("ids", "8000")
	query.Add("startDate", startDate)
	query.Add("endDate", endDate)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetCheckInsLocationRangeAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/range",
		nil,
	)

	rs1, _ := ts.Client().Do(req1)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
}

func TestGetCheckInsLocationDayRawSingle(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utc)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/all-locations/checkins/day",
			nil,
		)
		req.AddCookie(user)

		query := req.URL.Query()
		query.Add("ids", fixtureData.DefaultLocation.ID)
		query.Add("date", date.Format(constants.DateFormat))
		query.Add("returnType", "raw")

		req.URL.RawQuery = query.Encode()

		rs, _ := ts.Client().Do(req)

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		var checkInDate string
		for k := range rsData {
			checkInDate = k
			break
		}

		capacity, _ := rsData[checkInDate].Capacities.Get(
			fixtureData.DefaultLocation.ID,
		)
		assert.InRange(t, capacity, 21, 25)

		value, present := rsData[checkInDate].Schools.Get("Andere")
		assert.Equal(t, value, 5)
		assert.Equal(t, present, true)
	}
}

func TestGetCheckInsLocationDayRawMultiple(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utc)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/all-locations/checkins/day",
			nil,
		)
		req.AddCookie(user)

		query := req.URL.Query()
		query.Add(
			"ids",
			fmt.Sprintf(
				"%s,%s",
				fixtureData.DefaultLocation.ID,
				fixtureData.Locations[0].ID,
			),
		)
		query.Add("date", date.Format(constants.DateFormat))
		query.Add("returnType", "raw")

		req.URL.RawQuery = query.Encode()

		rs, _ := ts.Client().Do(req)

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		var checkInDate string
		for k := range rsData {
			checkInDate = k
			break
		}

		capacity0, _ := rsData[checkInDate].Capacities.Get(
			fixtureData.DefaultLocation.ID,
		)
		capacity1, _ := rsData[checkInDate].Capacities.Get(fixtureData.Locations[0].ID)
		assert.InRange(t, capacity0, 21, 25)

		assert.Equal(
			t,
			capacity1,
			fixtureData.Locations[0].Capacity,
		)

		value, present := rsData[checkInDate].Schools.Get("Andere")
		assert.Equal(t, value, 10)
		assert.Equal(t, present, true)
	}
}

func TestGetCheckInsLocationDayCSV(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	date := time.Now().Format(constants.DateFormat)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/all-locations/checkins/day",
			nil,
		)
		req.AddCookie(user)

		query := req.URL.Query()
		query.Add("ids", fixtureData.DefaultLocation.ID)
		query.Add("date", date)
		query.Add("returnType", "csv")

		req.URL.RawQuery = query.Encode()

		rs, _ := ts.Client().Do(req)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rs.Header.Get("content-type"), "text/csv")
	}
}

func TestGetCheckInsLocationDayNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	date := time.Now().Format(constants.DateFormat)

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/day",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	query := req.URL.Query()
	query.Add("ids", id.String())
	query.Add("date", date)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
	)
}

func TestGetCheckInsLocationDayNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	date := time.Now().Format(constants.DateFormat)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/day",
		nil,
	)
	req.AddCookie(tokens.DefaultAccessToken)

	query := req.URL.Query()
	query.Add("ids", fixtureData.Locations[0].ID)
	query.Add("date", date)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
	)
}

func TestGetCheckInsLocationDateMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/day",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	query := req.URL.Query()
	query.Add("ids", fixtureData.Locations[0].ID)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "missing date param in query")
}

func TestGetCheckInsLocationReturnTypeMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	date := time.Now().Format(constants.DateFormat)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/day",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	query := req.URL.Query()
	query.Add("ids", fixtureData.Locations[0].ID)
	query.Add("date", date)

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "missing returnType param in query")
}

func TestGetCheckInsLocationDayNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	date := time.Now().Format(constants.DateFormat)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/day",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	query := req.URL.Query()
	query.Add("ids", "8000")
	query.Add("date", date)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetCheckInsLocationDayAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/all-locations/checkins/day",
		nil,
	)

	rs1, _ := ts.Client().Do(req1)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
}

func TestGetAllCheckInsToday(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/locations/"+fixtureData.DefaultLocation.ID+"/checkins",
			nil,
		)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData []dtos.CheckInDto
		_ = helpers.ReadJSON(rs.Body, &rsData)

		loc, _ := time.LoadLocation("Europe/Brussels")

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, len(rsData), 5)
		assert.Equal(t, rsData[0].LocationID, fixtureData.DefaultLocation.ID)
		assert.Equal(t, rsData[0].SchoolName, "Andere")
		assert.Equal(
			t,
			rsData[0].CreatedAt.Time.Format(constants.DateFormat),
			time.Now().In(loc).Format(constants.DateFormat),
		)
	}
}

func TestGetAllCheckInsTodayNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/locations/"+id.String()+"/checkins",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
	)
}

func TestGetAllCheckInsTodayNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins",
		nil,
	)
	req.AddCookie(tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
	)
}

func TestGetAllCheckInsTodayNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/locations/8000/checkins",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetCheckInsTodayAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins",
		nil,
	)

	rs1, _ := ts.Client().Do(req1)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
}

func TestDeleteCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		req, _ := http.NewRequest(
			http.MethodDelete,
			ts.URL+"/locations/"+fixtureData.DefaultLocation.ID+"/checkins/"+strconv.FormatInt(
				fixtureData.CheckIns[i].ID,
				10,
			),
			nil,
		)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData dtos.CheckInDto
		_ = helpers.ReadJSON(rs.Body, &rsData)

		loc, _ := time.LoadLocation("Europe/Brussels")

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.CheckIns[i].ID)
		assert.Equal(t, rsData.LocationID, fixtureData.DefaultLocation.ID)
		assert.Equal(t, rsData.SchoolName, "Andere")
		assert.Equal(
			t,
			rsData.CreatedAt.Time.Format(constants.DateFormat),
			time.Now().In(loc).Format(constants.DateFormat),
		)
	}
}

func TestDeleteCheckInLocationNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/locations/"+id.String()+"/checkins/1",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
	)
}

func TestDeleteCheckInCheckInNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/locations/"+fixtureData.DefaultLocation.ID+"/checkins/8000",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		"checkIn with id '8000' doesn't exist",
	)
}

func TestDeleteCheckInNotToday(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	query := `
			INSERT INTO check_ins 
			(location_id, school_id, capacity, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`

	var checkIn models.CheckIn
	err := testEnv.TestTx.QueryRow(
		context.Background(),
		query,
		fixtureData.DefaultLocation.ID,
		1,
		fixtureData.DefaultLocation.Capacity,
		time.Now().AddDate(0, 0, -1),
	).Scan(&checkIn.ID)
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/locations/"+fixtureData.DefaultLocation.ID+"/checkins/"+strconv.FormatInt(
			checkIn.ID,
			10,
		),
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(
		t,
		rsData.Message,
		"checkIn didn't occur today and thus can't be deleted",
	)
}

func TestDeleteCheckInNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/locations/8000/checkins/1",
		nil,
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestDeleteCheckInAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/1",
		nil,
	)

	req2, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/1",
		nil,
	)

	req2.AddCookie(tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestGetPaginatedLocationsDefaultPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations", nil)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData dtos.PaginatedLocationsDto
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		assert.Equal(t, rsData.Pagination.Current, 1)
		assert.Equal(
			t,
			rsData.Pagination.Total,
			int64(math.Ceil(float64(fixtureData.AmountOfLocations)/3)),
		)
		assert.Equal(t, len(rsData.Data), 3)

		assert.Equal(t, rsData.Data[0].ID, fixtureData.DefaultLocation.ID)
		assert.Equal(t, rsData.Data[0].Name, fixtureData.DefaultLocation.Name)
		assert.Equal(
			t,
			rsData.Data[0].NormalizedName,
			fixtureData.DefaultLocation.NormalizedName,
		)
		assert.Equal(t, rsData.Data[0].Available, fixtureData.DefaultLocation.Available)
		assert.Equal(t, rsData.Data[0].Capacity, fixtureData.DefaultLocation.Capacity)
		assert.NotEqual(
			t,
			rsData.Data[0].AvailableYesterday,
			0,
		)
		assert.NotEqual(
			t,
			rsData.Data[0].CapacityYesterday,
			0,
		)
		assert.Equal(
			t,
			rsData.Data[0].YesterdayFullAt,
			fixtureData.DefaultLocation.YesterdayFullAt,
		)
		assert.Equal(t, rsData.Data[0].TimeZone, fixtureData.DefaultLocation.TimeZone)
		assert.Equal(t, rsData.Data[0].UserID, fixtureData.DefaultLocation.UserID)
	}
}

func TestGetPaginatedLocationsSpecificPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations?page=2", nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.PaginatedLocationsDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.Equal(t, rsData.Pagination.Current, 2)
	assert.Equal(
		t,
		rsData.Pagination.Total,
		int64(math.Ceil(float64(fixtureData.AmountOfLocations)/3)),
	)
	assert.Equal(t, len(rsData.Data), 3)

	assert.Equal(t, rsData.Data[0].ID, fixtureData.Locations[10].ID)
	assert.Equal(t, rsData.Data[0].Name, fixtureData.Locations[10].Name)
	assert.Equal(
		t,
		rsData.Data[0].NormalizedName,
		fixtureData.Locations[10].NormalizedName,
	)
	assert.Equal(t, rsData.Data[0].Available, fixtureData.Locations[10].Available)
	assert.Equal(t, rsData.Data[0].Capacity, fixtureData.Locations[10].Capacity)
	assert.NotEqual(
		t,
		rsData.Data[0].AvailableYesterday,
		0,
	)
	assert.NotEqual(
		t,
		rsData.Data[0].CapacityYesterday,
		0,
	)
	assert.Equal(
		t,
		rsData.Data[0].YesterdayFullAt,
		fixtureData.Locations[10].YesterdayFullAt,
	)
	assert.Equal(t, rsData.Data[0].TimeZone, fixtureData.Locations[10].TimeZone)
	assert.Equal(t, rsData.Data[0].UserID, fixtureData.Locations[10].UserID)
}

func TestGetPaginatedLocationsPageZero(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations?page=0", nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message, "invalid page query param")
}

func TestGetPaginatedLocationsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations", nil)

	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations", nil)
	req2.AddCookie(tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestGetAllLocations(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/all-locations", nil)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData []models.Location
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		assert.Equal(t, len(rsData), 21)
		assert.Equal(t, rsData[0].ID, fixtureData.DefaultLocation.ID)
		assert.Equal(t, rsData[0].Name, fixtureData.DefaultLocation.Name)
		assert.Equal(
			t,
			rsData[0].NormalizedName,
			fixtureData.DefaultLocation.NormalizedName,
		)
		assert.Equal(t, rsData[0].Available, fixtureData.DefaultLocation.Available)
		assert.Equal(t, rsData[0].Capacity, fixtureData.DefaultLocation.Capacity)
		assert.NotEqual(t, rsData[0].AvailableYesterday, 0)
		assert.NotEqual(t, rsData[0].CapacityYesterday, 0)
		assert.Equal(
			t,
			rsData[0].YesterdayFullAt,
			fixtureData.DefaultLocation.YesterdayFullAt,
		)
		assert.Equal(t, rsData[0].TimeZone, fixtureData.DefaultLocation.TimeZone)
		assert.Equal(t, rsData[0].UserID, fixtureData.DefaultLocation.UserID)
	}
}

func TestGetAllLocationsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodGet, ts.URL+"/all-locations", nil)

	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/all-locations", nil)
	req2.AddCookie(tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestGetLocation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/locations/"+fixtureData.DefaultLocation.ID,
			nil,
		)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData models.Location
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.DefaultLocation.ID)
		assert.Equal(t, rsData.Name, fixtureData.DefaultLocation.Name)
		assert.Equal(
			t,
			rsData.NormalizedName,
			fixtureData.DefaultLocation.NormalizedName,
		)
		assert.Equal(t, rsData.Available, fixtureData.DefaultLocation.Available)
		assert.Equal(t, rsData.Capacity, fixtureData.DefaultLocation.Capacity)
		assert.NotEqual(t, rsData.AvailableYesterday, 0)
		assert.NotEqual(t, rsData.CapacityYesterday, 0)
		assert.Equal(
			t,
			rsData.YesterdayFullAt,
			fixtureData.DefaultLocation.YesterdayFullAt,
		)
		assert.Equal(t, rsData.TimeZone, fixtureData.DefaultLocation.TimeZone)
		assert.Equal(t, rsData.UserID, fixtureData.DefaultLocation.UserID)
	}
}

func TestGetLocationNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations/"+id.String(), nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
	)
}

func TestGetLocationNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		nil,
	)
	req.AddCookie(tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
	)
}

func TestGetLocationNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations/8000", nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetLocationAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		nil,
	)

	rs1, _ := ts.Client().Do(req1)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
}

func TestCreateLocation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		unique := fmt.Sprintf("test%d", i)

		data := dtos.CreateLocationDto{
			Name:     unique,
			Capacity: 10,
			Username: unique,
			Password: "testpassword",
			TimeZone: "Europe/Brussels",
		}

		body, _ := json.Marshal(data)

		req, _ := http.NewRequest(
			http.MethodPost,
			ts.URL+"/locations",
			bytes.NewReader(body),
		)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData models.Location
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusCreated)
		assert.IsUUID(t, rsData.ID)
		assert.Equal(t, rsData.Name, data.Name)
		assert.Equal(t, rsData.NormalizedName, data.Name)
		assert.Equal(t, rsData.Available, data.Capacity)
		assert.Equal(t, rsData.Capacity, data.Capacity)
		assert.Equal(t, rsData.AvailableYesterday, data.Capacity)
		assert.Equal(t, rsData.CapacityYesterday, data.Capacity)
		assert.Equal(t, rsData.YesterdayFullAt.Valid, false)
		assert.Equal(t, rsData.TimeZone, data.TimeZone)
		assert.IsUUID(t, rsData.UserID)
	}
}

func TestCreateLocationNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CreateLocationDto{
		Name:     "TestLocation0",
		Capacity: 10,
		Username: "test",
		Password: "testpassword",
		TimeZone: "Europe/Brussels",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/locations",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["name"].(string),
		fmt.Sprintf("location with name '%s' already exists", data.Name),
	)
}

func TestCreateLocationNormalizedNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CreateLocationDto{
		Name:     "$TestLocation0$",
		Capacity: 10,
		Username: "test",
		Password: "testpassword",
		TimeZone: "Europe/Brussels",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/locations",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["name"].(string),
		fmt.Sprintf("location with name '%s' already exists", data.Name),
	)
}

func TestCreateLocationUserNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CreateLocationDto{
		Name:     "test",
		Capacity: 10,
		Username: "TestDefaultUser1",
		Password: "testpassword",
		TimeZone: "Europe/Brussels",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/locations",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["username"].(string),
		fmt.Sprintf("user with username '%s' already exists", data.Username),
	)
}

func TestCreateLocationInvalidCapacity(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CreateLocationDto{
		Name:     "test",
		Capacity: -1,
		Username: "test",
		Password: "testpassword",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/locations",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["capacity"],
		"must be greater than zero",
	)
}

func TestCreateLocationInvalidTimeZone(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	data := dtos.CreateLocationDto{
		Name:     "test",
		Capacity: 10,
		Username: "test",
		Password: "testpassword",
		TimeZone: "wrong",
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/locations",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["timeZone"],
		"must be provided and must be a valid IANA value",
	)
}

func TestCreateLocationAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodPost, ts.URL+"/locations", nil)

	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/locations", nil)
	req2.AddCookie(tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}

func TestUpdateLocation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		unique := fmt.Sprintf("test%d", i)
		name, username, password := unique, unique, "testpassword"
		timeZone := "Europe/Brussels"
		var capacity int64 = 3
		data := dtos.UpdateLocationDto{
			Name:     &name,
			Capacity: &capacity,
			Username: &username,
			Password: &password,
			TimeZone: &timeZone,
		}

		body, _ := json.Marshal(data)

		req, _ := http.NewRequest(
			http.MethodPatch,
			ts.URL+"/locations/"+fixtureData.Locations[0].ID,
			bytes.NewReader(body),
		)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData models.Location
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.Locations[0].ID)
		assert.Equal(t, rsData.Name, *data.Name)
		assert.Equal(t, rsData.NormalizedName, *data.Name)
		assert.Equal(t, rsData.Available, 0)
		assert.Equal(t, rsData.Capacity, *data.Capacity)
		assert.Equal(
			t,
			rsData.AvailableYesterday,
			fixtureData.Locations[0].AvailableYesterday,
		)
		assert.Equal(
			t,
			rsData.CapacityYesterday,
			fixtureData.Locations[0].CapacityYesterday,
		)
		assert.Equal(t, rsData.YesterdayFullAt.Valid, false)
		assert.Equal(t, rsData.TimeZone, *data.TimeZone)
		assert.Equal(t, rsData.UserID, fixtureData.Locations[0].UserID)
	}
}

func TestUpdateLocationNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	name, username, password := "TestLocation1", "test", "testpassword"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["name"].(string),
		fmt.Sprintf("location with name '%s' already exists", *data.Name),
	)
}

func TestUpdateLocationNormalizedNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	name, username, password := "$TestLocation1$", "test", "testpassword"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["name"].(string),
		fmt.Sprintf("location with name '%s' already exists", *data.Name),
	)
}

func TestUpdateLocationUserNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	name, username, password := "test", "TestDefaultUser1", "testpassword"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["username"].(string),
		fmt.Sprintf("user with username '%s' already exists", *data.Username),
	)
}

func TestUpdateLocationInvalidCapacity(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	name, username, password := "test", "test", "testpassword"
	var capacity int64 = -1
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["capacity"],
		"must be greater than zero",
	)
}

func TestUpdateLocationInvalidTimeZone(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	name, username, password, timeZone := "test", "test", "testpassword", "wrong"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["timeZone"],
		"must be provided and must be a valid IANA value",
	)
}

func TestUpdateLocationNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	name, username, password := "test", "test", "password"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/locations/"+id.String(),
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
	)
}

func TestUpdateLocationNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	name, username, password := "test", "test", "testpassword"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
	)
}

func TestUpdateLocationNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	name, username, password := "test", "test", "password"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	body, _ := json.Marshal(data)

	req, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/locations/8000",
		bytes.NewReader(body),
	)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestUpdateLocationAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodPatch,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		nil,
	)

	rs1, _ := ts.Client().Do(req1)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
}

func TestDeleteLocation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		req, _ := http.NewRequest(
			http.MethodDelete,
			ts.URL+"/locations/"+fixtureData.Locations[i].ID,
			nil,
		)
		req.AddCookie(user)

		rs, _ := ts.Client().Do(req)

		var rsData models.Location
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.Locations[i].ID)
		assert.Equal(t, rsData.Name, fixtureData.Locations[i].Name)
		assert.Equal(t, rsData.NormalizedName, fixtureData.Locations[i].NormalizedName)
		assert.Equal(t, rsData.Available, fixtureData.Locations[i].Available)
		assert.Equal(t, rsData.Capacity, fixtureData.Locations[i].Capacity)
		assert.Equal(
			t,
			rsData.AvailableYesterday,
			fixtureData.Locations[i].AvailableYesterday,
		)
		assert.Equal(
			t,
			rsData.CapacityYesterday,
			fixtureData.Locations[i].CapacityYesterday,
		)
		assert.Equal(
			t,
			rsData.YesterdayFullAt,
			fixtureData.Locations[i].YesterdayFullAt,
		)
		assert.Equal(t, rsData.TimeZone, fixtureData.Locations[i].TimeZone)
		assert.Equal(t, rsData.UserID, fixtureData.Locations[i].UserID)
	}
}

func TestDeleteLocationNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	id, _ := uuid.NewUUID()
	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/locations/"+id.String(), nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["id"].(string),
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
	)
}

func TestDeleteLocationNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/locations/8000", nil)
	req.AddCookie(tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestDeleteLocationAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		nil,
	)

	req2, _ := http.NewRequest(
		http.MethodDelete,
		ts.URL+"/locations/"+fixtureData.Locations[0].ID,
		nil,
	)
	req2.AddCookie(tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}
