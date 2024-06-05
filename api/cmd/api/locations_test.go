package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/google/uuid"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
	"check-in/api/internal/tests"

	"github.com/stretchr/testify/assert"
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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/"+fixtureData.DefaultLocation.ID)
	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData models.Location
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.EqualValues(t, rsData.AvailableYesterday, 0)
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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
		tReq.AddCookie(user)

		tReq.SetQuery(map[string]string{
			"ids":        fixtureData.DefaultLocation.ID,
			"startDate":  startDate.Format(constants.DateFormat),
			"endDate":    endDate.Format(constants.DateFormat),
			"returnType": "raw",
		})

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		rs := tReq.Do(t, &rsData)

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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
		tReq.AddCookie(user)

		id := fmt.Sprintf(
			"%s,%s",
			fixtureData.DefaultLocation.ID,
			fixtureData.Locations[0].ID,
		)

		tReq.SetQuery(map[string]string{
			"ids":        id,
			"startDate":  startDate.Format(constants.DateFormat),
			"endDate":    endDate.Format(constants.DateFormat),
			"returnType": "raw",
		})

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		rs := tReq.Do(t, &rsData)

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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
		tReq.AddCookie(user)

		tReq.SetQuery(map[string]string{
			"ids":        fixtureData.DefaultLocation.ID,
			"startDate":  startDate,
			"endDate":    endDate,
			"returnType": "csv",
		})

		rs := tReq.Do(t, nil)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        id.String(),
		"startDate":  startDate,
		"endDate":    endDate,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        fixtureData.Locations[0].ID,
		"startDate":  startDate,
		"endDate":    endDate,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        fixtureData.Locations[0].ID,
		"endDate":    endDate,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "missing startDate param in query")
}

func TestGetCheckInsLocationRangeEndDateMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	startDate := time.Now().Format(constants.DateFormat)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        fixtureData.Locations[0].ID,
		"startDate":  startDate,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":       fixtureData.Locations[0].ID,
		"startDate": startDate,
		"endDate":   endDate,
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        "8000",
		"startDate":  startDate,
		"endDate":    endDate,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetCheckInsLocationRangeAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq1 := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")
	rs1 := tReq1.Do(t, nil)

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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
		tReq.AddCookie(user)

		tReq.SetQuery(map[string]string{
			"ids":        fixtureData.DefaultLocation.ID,
			"date":       date.Format(constants.DateFormat),
			"returnType": "raw",
		})

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		var checkInDate string
		for k := range rsData {
			checkInDate = k
			break
		}

		capacity, _ := rsData[checkInDate].Capacities.Get(
			fixtureData.DefaultLocation.ID,
		)
		assert.GreaterOrEqual(t, capacity, int64(21))
		assert.LessOrEqual(t, capacity, int64(25))

		//value, present := rsData[checkInDate].Schools.Get("Andere")
		//todo assert.Equal(t, value, 5)
		//assert.Equal(t, present, true)
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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
		tReq.AddCookie(user)

		id := fmt.Sprintf(
			"%s,%s",
			fixtureData.DefaultLocation.ID,
			fixtureData.Locations[0].ID,
		)
		tReq.SetQuery(map[string]string{
			"ids":        id,
			"date":       date.Format(constants.DateFormat),
			"returnType": "raw",
		})

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		var checkInDate string
		for k := range rsData {
			checkInDate = k
			break
		}

		capacity0, _ := rsData[checkInDate].Capacities.Get(
			fixtureData.DefaultLocation.ID,
		)
		//capacity1, _ := rsData[checkInDate].Capacities.Get(fixtureData.Locations[0].ID)
		assert.GreaterOrEqual(t, capacity0, int64(21))
		assert.LessOrEqual(t, capacity0, int64(25))

		//todo assert.Equal(
		//	t,
		//	capacity1,
		//	fixtureData.Locations[0].Capacity,
		//)

		//value, present := rsData[checkInDate].Schools.Get("Andere")
		//todoassert.Equal(t, value, 10)
		//assert.Equal(t, present, true)
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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
		tReq.AddCookie(user)

		tReq.SetQuery(map[string]string{
			"ids":        fixtureData.DefaultLocation.ID,
			"date":       date,
			"returnType": "csv",
		})

		rs := tReq.Do(t, nil)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        id.String(),
		"date":       date,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        fixtureData.Locations[0].ID,
		"date":       date,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        fixtureData.Locations[0].ID,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "missing date param in query")
}

func TestGetCheckInsLocationReturnTypeMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	date := time.Now().Format(constants.DateFormat)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":  fixtureData.Locations[0].ID,
		"date": date,
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Equal(t, rsData.Message.(string), "missing returnType param in query")
}

func TestGetCheckInsLocationDayNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	date := time.Now().Format(constants.DateFormat)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        "8000",
		"date":       date,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusBadRequest)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetCheckInsLocationDayAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	tReq1 := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
	rs1 := tReq1.Do(t, nil)

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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/"+fixtureData.DefaultLocation.ID+"/checkins")
		tReq.AddCookie(user)

		var rsData []dtos.CheckInDto
		rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/"+id.String()+"/checkins")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/"+fixtureData.Locations[0].ID+"/checkins")
	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/8000/checkins")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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
		id := strconv.FormatInt(
			fixtureData.CheckIns[i].ID,
			10,
		)
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/"+fixtureData.DefaultLocation.ID+"/checkins/"+id)
		tReq.AddCookie(user)

		var rsData dtos.CheckInDto
		rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/"+id.String()+"/checkins/1")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/"+fixtureData.DefaultLocation.ID+"/checkins/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	id := strconv.FormatInt(
		checkIn.ID,
		10,
	)
	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/"+fixtureData.DefaultLocation.ID+"/checkins/"+id)
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/800O/checkins/1")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations")
		tReq.AddCookie(user)

		var rsData dtos.PaginatedLocationsDto
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		assert.EqualValues(t, rsData.Pagination.Current, 1)
		assert.EqualValues(
			t,
			rsData.Pagination.Total,
			math.Ceil(float64(fixtureData.AmountOfLocations)/3),
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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	var rsData dtos.PaginatedLocationsDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	assert.EqualValues(t, rsData.Pagination.Current, 2)
	assert.EqualValues(
		t,
		rsData.Pagination.Total,
		math.Ceil(float64(fixtureData.AmountOfLocations)/3),
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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "0",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations")
		tReq.AddCookie(user)

		var rsData []models.Location
		rs := tReq.Do(t, &rsData)

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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/"+fixtureData.DefaultLocation.ID)
		tReq.AddCookie(user)

		var rsData models.Location
		rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/"+id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/"+fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

		tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/locations")
		tReq.AddCookie(user)

		tReq.SetReqData(data)

		var rsData models.Location
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusCreated)
		assert.Nil(t, uuid.Validate(rsData.ID))
		assert.Equal(t, rsData.Name, data.Name)
		assert.Equal(t, rsData.NormalizedName, data.Name)
		assert.Equal(t, rsData.Available, data.Capacity)
		assert.Equal(t, rsData.Capacity, data.Capacity)
		assert.Equal(t, rsData.AvailableYesterday, data.Capacity)
		assert.Equal(t, rsData.CapacityYesterday, data.Capacity)
		assert.Equal(t, rsData.YesterdayFullAt.Valid, false)
		assert.Equal(t, rsData.TimeZone, data.TimeZone)
		assert.Nil(t, uuid.Validate(rsData.ID))
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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/locations")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/locations")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/locations")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["username"].(string),
		fmt.Sprintf("user with username '%s' already exists", data.Username),
	)
}

func TestCreateLocationFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	vt := test.CreateValidatorTester(t)

	tReq1 := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/locations")
	tReq1.AddCookie(tokens.ManagerAccessToken)

	tReq1.SetReqData(dtos.CreateLocationDto{
		Name:     "test",
		Capacity: -1,
		Username: "test",
		Password: "testpassword",
		TimeZone: "Europe/Brussels",
	})

	vt.AddTestCase(tReq1, map[string]interface{}{
		"capacity": "must be greater than zero",
	})

	tReq2 := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/locations")
	tReq2.AddCookie(tokens.ManagerAccessToken)

	tReq2.SetReqData(dtos.CreateLocationDto{
		Name:     "test",
		Capacity: 10,
		Username: "test",
		Password: "testpassword",
		TimeZone: "wrong",
	})

	vt.AddTestCase(tReq2, map[string]interface{}{
		"timeZone": "must be a valid IANA value",
	})

	vt.Do(t)
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

		tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/"+fixtureData.Locations[0].ID)
		tReq.AddCookie(user)

		tReq.SetReqData(data)

		var rsData models.Location
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.Locations[0].ID)
		assert.Equal(t, rsData.Name, *data.Name)
		assert.Equal(t, rsData.NormalizedName, *data.Name)
		assert.EqualValues(t, rsData.Available, 0)
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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/"+fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/"+fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/"+fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(map[string]interface{})["username"].(string),
		fmt.Sprintf("user with username '%s' already exists", *data.Username),
	)
}

func TestUpdateLocationFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	vt := test.CreateValidatorTester(t)

	name, username, password := "test", "test", "testpassword"
	var capacity1 int64 = -1
	data1 := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity1,
		Username: &username,
		Password: &password,
	}

	tReq1 := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/"+fixtureData.Locations[0].ID)
	tReq1.AddCookie(tokens.ManagerAccessToken)

	tReq1.SetReqData(data1)

	vt.AddTestCase(tReq1, map[string]interface{}{
		"capacity": "must be greater than zero",
	})

	name, username, password, timeZone := "test", "test", "testpassword", "wrong"
	var capacity2 int64 = 10
	data2 := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity2,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone,
	}

	tReq2 := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/"+fixtureData.Locations[0].ID)
	tReq2.AddCookie(tokens.ManagerAccessToken)

	tReq2.SetReqData(data2)

	vt.AddTestCase(tReq2, map[string]interface{}{
		"timeZone": "must be a valid IANA value",
	})

	vt.Do(t)
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

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/"+id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/"+fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/"+fixtureData.Locations[i].ID)
		tReq.AddCookie(user)

		var rsData models.Location
		rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/"+id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

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
