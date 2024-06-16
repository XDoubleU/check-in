package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/XDoubleU/essentia/pkg/tools"
	"github.com/google/uuid"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestYesterdayFullAt(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/%s", fixtureData.DefaultLocation.ID)
	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData models.Location
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 0, rsData.AvailableYesterday)
	assert.Equal(t, fixtureData.DefaultLocation.Capacity, rsData.CapacityYesterday)
	assert.Equal(t, true, rsData.YesterdayFullAt.Valid)
	assert.Equal(t, now.Day(), rsData.YesterdayFullAt.Time.Day())
	assert.Equal(t, now.Hour(), rsData.YesterdayFullAt.Time.Hour())
}

func TestGetCheckInsLocationRangeRawSingle(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	startDate := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, utc)
	startDate = tools.StartOfDay(startDate)

	endDate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, utc)
	endDate = tools.StartOfDay(endDate)

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

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		assert.Equal(t, 0, rsData[startDate.Format(time.RFC3339)].Capacities.Len())

		value, present := rsData[startDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, 0, value)
		assert.Equal(t, true, present)

		capacity, _ := rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].Capacities.Get(
			fixtureData.DefaultLocation.ID,
		)
		assert.Equal(
			t,
			fixtureData.DefaultLocation.Capacity,
			capacity,
		)

		value, present = rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].
			Schools.Get(
			"Andere",
		)
		assert.Equal(t, 5, value)
		assert.Equal(t, true, present)

		assert.Equal(t, 0, rsData[endDate.Format(time.RFC3339)].Capacities.Len())

		value, present = rsData[endDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, 0, value)
		assert.Equal(t, true, present)
	}
}

func TestGetCheckInsLocationRangeRawMultiple(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	startDate := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, utc)
	startDate = tools.StartOfDay(startDate)

	endDate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, utc)
	endDate = tools.StartOfDay(endDate)

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

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		assert.Equal(t, 0, rsData[startDate.Format(time.RFC3339)].Capacities.Len())

		value, present := rsData[startDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, 0, value)
		assert.Equal(t, true, present)

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
			fixtureData.DefaultLocation.Capacity,
			capacity0,
		)

		assert.Equal(
			t,
			fixtureData.Locations[0].Capacity,
			capacity1,
		)

		value, present = rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].
			Schools.Get(
			"Andere",
		)
		assert.Equal(t, 10, value)
		assert.Equal(t, true, present)

		assert.Equal(t, 0, rsData[endDate.Format(time.RFC3339)].Capacities.Len())

		value, present = rsData[endDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, 0, value)
		assert.Equal(t, true, present)
	}
}

func TestGetCheckInsLocationRangeCSV(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, "text/csv", rs.Header.Get("content-type"))
	}
}

func TestGetCheckInsLocationRangeNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetCheckInsLocationRangeNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetCheckInsLocationRangeStartDateMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'startDate'", rsData.Message.(string))
}

func TestGetCheckInsLocationRangeEndDateMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'endDate'", rsData.Message.(string))
}

func TestGetCheckInsLocationRangeReturnTypeMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'returnType'", rsData.Message.(string))
}

func TestGetCheckInsLocationRangeNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetCheckInsLocationRangeAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/range")

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)

	mt.Do(t)
}

func TestGetCheckInsLocationDayRawSingle(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

		assert.Equal(t, http.StatusOK, rs.StatusCode)

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
		//todo assert.Equal(t, 5, value)
		//assert.Equal(t, true, present)
	}
}

func TestGetCheckInsLocationDayRawMultiple(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

		assert.Equal(t, http.StatusOK, rs.StatusCode)

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
		//todoassert.Equal(t, 10, value)
		//assert.Equal(t, true, present)
	}
}

func TestGetCheckInsLocationDayCSV(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, "text/csv", rs.Header.Get("content-type"))
	}
}

func TestGetCheckInsLocationDayNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetCheckInsLocationDayNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetCheckInsLocationDateMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        fixtureData.Locations[0].ID,
		"returnType": "raw",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'date'", rsData.Message.(string))
}

func TestGetCheckInsLocationReturnTypeMissing(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	date := time.Now().Format(constants.DateFormat)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":  fixtureData.Locations[0].ID,
		"date": date,
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'returnType'", rsData.Message.(string))
}

func TestGetCheckInsLocationDayNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetCheckInsLocationDayAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations/checkins/day")

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)

	mt.Do(t)
}

func TestGetAllCheckInsToday(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/%s/checkins", fixtureData.DefaultLocation.ID)
		tReq.AddCookie(user)

		var rsData []dtos.CheckInDto
		rs := tReq.Do(t, &rsData)

		loc, _ := time.LoadLocation("Europe/Brussels")

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, 5, len(rsData))
		assert.Equal(t, fixtureData.DefaultLocation.ID, rsData[0].LocationID)
		assert.Equal(t, "Andere", rsData[0].SchoolName)
		assert.Equal(
			t,
			time.Now().In(loc).Format(constants.DateFormat),
			rsData[0].CreatedAt.Time.Format(constants.DateFormat),
		)
	}
}

func TestGetAllCheckInsTodayNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/%s/checkins", id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetAllCheckInsTodayNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/%s/checkins", fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetAllCheckInsTodayNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/8000/checkins")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetCheckInsTodayAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/%s/checkins", fixtureData.Locations[0].ID)

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)

	mt.Do(t)
}

func TestDeleteCheckIn(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		id := strconv.FormatInt(
			fixtureData.CheckIns[i].ID,
			10,
		)
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/%s/checkins/%s", fixtureData.DefaultLocation.ID, id)
		tReq.AddCookie(user)

		var rsData dtos.CheckInDto
		rs := tReq.Do(t, &rsData)

		loc, _ := time.LoadLocation("Europe/Brussels")

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, fixtureData.CheckIns[i].ID, rsData.ID)
		assert.Equal(t, fixtureData.DefaultLocation.ID, rsData.LocationID)
		assert.Equal(t, "Andere", rsData.SchoolName)
		assert.Equal(
			t,
			time.Now().In(loc).Format(constants.DateFormat),
			rsData.CreatedAt.Time.Format(constants.DateFormat),
		)
	}
}

func TestDeleteCheckInLocationNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/%s/checkins/1", id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteCheckInCheckInNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/%s/checkins/8000", fixtureData.DefaultLocation.ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"checkIn with id '8000' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteCheckInNotToday(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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
	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/%s/checkins/%s", fixtureData.DefaultLocation.ID, id)
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(
		t,
		"checkIn didn't occur today and thus can't be deleted",
		rsData.Message,
	)
}

func TestDeleteCheckInNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/800O/checkins/1")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestDeleteCheckInAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/%s/checkins/1", fixtureData.Locations[0].ID)

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestGetPaginatedLocationsDefaultPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations")
		tReq.AddCookie(user)

		var rsData dtos.PaginatedLocationsDto
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		assert.EqualValues(t, 1, rsData.Pagination.Current)
		assert.EqualValues(
			t,
			math.Ceil(float64(fixtureData.AmountOfLocations)/3),
			rsData.Pagination.Total,
		)
		assert.Equal(t, 3, len(rsData.Data))

		assert.Equal(t, fixtureData.DefaultLocation.ID, rsData.Data[0].ID)
		assert.Equal(t, fixtureData.DefaultLocation.Name, rsData.Data[0].Name)
		assert.Equal(
			t,
			fixtureData.DefaultLocation.NormalizedName,
			rsData.Data[0].NormalizedName,
		)
		assert.Equal(t, fixtureData.DefaultLocation.Available, rsData.Data[0].Available)
		assert.Equal(t, fixtureData.DefaultLocation.Capacity, rsData.Data[0].Capacity)
		assert.NotEqual(
			t,
			0,
			rsData.Data[0].AvailableYesterday,
		)
		assert.NotEqual(
			t,
			0,
			rsData.Data[0].CapacityYesterday,
		)
		assert.Equal(
			t,
			fixtureData.DefaultLocation.YesterdayFullAt,
			rsData.Data[0].YesterdayFullAt,
		)
		assert.Equal(t, fixtureData.DefaultLocation.TimeZone, rsData.Data[0].TimeZone)
		assert.Equal(t, fixtureData.DefaultLocation.UserID, rsData.Data[0].UserID)
	}
}

func TestGetPaginatedLocationsSpecificPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	var rsData dtos.PaginatedLocationsDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 2, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(fixtureData.AmountOfLocations)/3),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 3, len(rsData.Data))

	assert.Equal(t, fixtureData.Locations[10].ID, rsData.Data[0].ID)
	assert.Equal(t, fixtureData.Locations[10].Name, rsData.Data[0].Name)
	assert.Equal(
		t,
		fixtureData.Locations[10].NormalizedName,
		rsData.Data[0].NormalizedName,
	)
	assert.Equal(t, fixtureData.Locations[10].Available, rsData.Data[0].Available)
	assert.Equal(t, fixtureData.Locations[10].Capacity, rsData.Data[0].Capacity)
	assert.NotEqual(
		t,
		0,
		rsData.Data[0].AvailableYesterday,
	)
	assert.NotEqual(
		t,
		0,
		rsData.Data[0].CapacityYesterday,
	)
	assert.Equal(
		t,
		fixtureData.Locations[10].YesterdayFullAt,
		rsData.Data[0].YesterdayFullAt,
	)
	assert.Equal(t, fixtureData.Locations[10].TimeZone, rsData.Data[0].TimeZone)
	assert.Equal(t, fixtureData.Locations[10].UserID, rsData.Data[0].UserID)
}

func TestGetPaginatedLocationsPageZero(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations")
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "0",
	})

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "invalid query param 'page' with value '0', can't be '0'", rsData.Message)
}

func TestGetPaginatedLocationsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations")

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestGetAllLocations(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations")
		tReq.AddCookie(user)

		var rsData []models.Location
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		assert.Equal(t, 21, len(rsData))
		assert.Equal(t, fixtureData.DefaultLocation.ID, rsData[0].ID)
		assert.Equal(t, fixtureData.DefaultLocation.Name, rsData[0].Name)
		assert.Equal(
			t,
			fixtureData.DefaultLocation.NormalizedName,
			rsData[0].NormalizedName,
		)
		assert.Equal(t, fixtureData.DefaultLocation.Available, rsData[0].Available)
		assert.Equal(t, fixtureData.DefaultLocation.Capacity, rsData[0].Capacity)
		assert.NotEqual(t, 0, rsData[0].AvailableYesterday)
		assert.NotEqual(t, 0, rsData[0].CapacityYesterday)
		assert.Equal(
			t,
			fixtureData.DefaultLocation.YesterdayFullAt,
			rsData[0].YesterdayFullAt,
		)
		assert.Equal(t, fixtureData.DefaultLocation.TimeZone, rsData[0].TimeZone)
		assert.Equal(t, fixtureData.DefaultLocation.UserID, rsData[0].UserID)
	}
}

func TestGetAllLocationsAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/all-locations")

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestGetLocation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/%s", fixtureData.DefaultLocation.ID)
		tReq.AddCookie(user)

		var rsData models.Location
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, fixtureData.DefaultLocation.ID, rsData.ID)
		assert.Equal(t, fixtureData.DefaultLocation.Name, rsData.Name)
		assert.Equal(
			t,
			fixtureData.DefaultLocation.NormalizedName,
			rsData.NormalizedName,
		)
		assert.Equal(t, fixtureData.DefaultLocation.Available, rsData.Available)
		assert.Equal(t, fixtureData.DefaultLocation.Capacity, rsData.Capacity)
		assert.NotEqual(t, 0, rsData.AvailableYesterday)
		assert.NotEqual(t, 0, rsData.CapacityYesterday)
		assert.Equal(
			t,
			fixtureData.DefaultLocation.YesterdayFullAt,
			rsData.YesterdayFullAt,
		)
		assert.Equal(t, fixtureData.DefaultLocation.TimeZone, rsData.TimeZone)
		assert.Equal(t, fixtureData.DefaultLocation.UserID, rsData.UserID)
	}
}

func TestGetLocationNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/%s", id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetLocationNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/%s", fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.DefaultAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetLocationNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestGetLocationAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodGet, "/locations/%s", fixtureData.Locations[0].ID)

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)

	mt.Do(t)
}

func TestCreateLocation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

		assert.Equal(t, http.StatusCreated, rs.StatusCode)
		assert.Nil(t, uuid.Validate(rsData.ID))
		assert.Equal(t, data.Name, rsData.Name)
		assert.Equal(t, data.Name, rsData.NormalizedName)
		assert.Equal(t, data.Capacity, rsData.Available)
		assert.Equal(t, data.Capacity, rsData.Capacity)
		assert.Equal(t, data.Capacity, rsData.AvailableYesterday)
		assert.Equal(t, data.Capacity, rsData.CapacityYesterday)
		assert.Equal(t, false, rsData.YesterdayFullAt.Valid)
		assert.Equal(t, data.TimeZone, rsData.TimeZone)
		assert.Nil(t, uuid.Validate(rsData.ID))
	}
}

func TestCreateLocationNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with name '%s' already exists", data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestCreateLocationNormalizedNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with name '%s' already exists", data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestCreateLocationUserNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with username '%s' already exists", data.Username),
		rsData.Message.(map[string]interface{})["username"].(string),
	)
}

func TestCreateLocationFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/locations")
	tReq.AddCookie(tokens.ManagerAccessToken)

	mt := test.CreateMatrixTester(t, tReq)

	data1 := dtos.CreateLocationDto{
		Name:     "test",
		Capacity: -1,
		Username: "test",
		Password: "testpassword",
		TimeZone: "Europe/Brussels",
	}

	mt.AddTestCaseErrorMessage(data1, map[string]interface{}{
		"capacity": "must be greater than 0",
	})

	data2 := dtos.CreateLocationDto{
		Name:     "test",
		Capacity: 10,
		Username: "test",
		Password: "testpassword",
		TimeZone: "wrong",
	}

	mt.AddTestCaseErrorMessage(data2, map[string]interface{}{
		"timeZone": "must be a valid IANA value",
	})

	mt.Do(t)
}

func TestCreateLocationAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPost, "/locations")

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}

func TestUpdateLocation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

		tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/%s", fixtureData.Locations[0].ID)
		tReq.AddCookie(user)

		tReq.SetReqData(data)

		var rsData models.Location
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, fixtureData.Locations[0].ID, rsData.ID)
		assert.Equal(t, *data.Name, rsData.Name)
		assert.Equal(t, *data.Name, rsData.NormalizedName)
		assert.EqualValues(t, 0, rsData.Available)
		assert.Equal(t, *data.Capacity, rsData.Capacity)
		assert.Equal(
			t,
			fixtureData.Locations[0].AvailableYesterday,
			rsData.AvailableYesterday,
		)
		assert.Equal(
			t,
			fixtureData.Locations[0].CapacityYesterday,
			rsData.CapacityYesterday,
		)
		assert.Equal(t, false, rsData.YesterdayFullAt.Valid)
		assert.Equal(t, *data.TimeZone, rsData.TimeZone)
		assert.Equal(t, fixtureData.Locations[0].UserID, rsData.UserID)
	}
}

func TestUpdateLocationNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	name, username, password := "TestLocation1", "test", "testpassword"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/%s", fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with name '%s' already exists", *data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestUpdateLocationNormalizedNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	name, username, password := "$TestLocation1$", "test", "testpassword"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/%s", fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with name '%s' already exists", *data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestUpdateLocationUserNameExists(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	name, username, password := "test", "TestDefaultUser1", "testpassword"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/%s", fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with username '%s' already exists", *data.Username),
		rsData.Message.(map[string]interface{})["username"].(string),
	)
}

func TestUpdateLocationFailValidation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/%s", fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.ManagerAccessToken)

	mt := test.CreateMatrixTester(t, tReq)

	name, username, password := "test", "test", "testpassword"
	var capacity1 int64 = -1
	data1 := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity1,
		Username: &username,
		Password: &password,
	}

	mt.AddTestCaseErrorMessage(data1, map[string]interface{}{
		"capacity": "must be greater than 0",
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

	mt.AddTestCaseErrorMessage(data2, map[string]interface{}{
		"timeZone": "must be a valid IANA value",
	})

	mt.Do(t)
}

func TestUpdateLocationNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	name, username, password := "test", "test", "password"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/%s", id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateLocationNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	name, username, password := "test", "test", "testpassword"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
	}

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/%s", fixtureData.Locations[0].ID)
	tReq.AddCookie(tokens.DefaultAccessToken)

	tReq.SetReqData(data)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateLocationNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

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

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestUpdateLocationAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodPatch, "/locations/%s", fixtureData.Locations[0].ID)

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)

	mt.Do(t)
}

func TestDeleteLocation(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	users := []*http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
	}

	for i, user := range users {
		tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/%s", fixtureData.Locations[i].ID)
		tReq.AddCookie(user)

		var rsData models.Location
		rs := tReq.Do(t, &rsData)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, fixtureData.Locations[i].ID, rsData.ID)
		assert.Equal(t, fixtureData.Locations[i].Name, rsData.Name)
		assert.Equal(t, fixtureData.Locations[i].NormalizedName, rsData.NormalizedName)
		assert.Equal(t, fixtureData.Locations[i].Available, rsData.Available)
		assert.Equal(t, fixtureData.Locations[i].Capacity, rsData.Capacity)
		assert.Equal(
			t,
			fixtureData.Locations[i].AvailableYesterday,
			rsData.AvailableYesterday,
		)
		assert.Equal(
			t,
			fixtureData.Locations[i].CapacityYesterday,
			rsData.CapacityYesterday,
		)
		assert.Equal(
			t,
			fixtureData.Locations[i].YesterdayFullAt,
			rsData.YesterdayFullAt,
		)
		assert.Equal(t, fixtureData.Locations[i].TimeZone, rsData.TimeZone)
		assert.Equal(t, fixtureData.Locations[i].UserID, rsData.UserID)
	}
}

func TestDeleteLocationNotFound(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	id, _ := uuid.NewUUID()

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/%s", id.String())
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteLocationNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/8000")
	tReq.AddCookie(tokens.ManagerAccessToken)

	var rsData http_tools.ErrorDto
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "invalid UUID")
}

func TestDeleteLocationAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer test.TeardownSingle(testEnv)

	tReq := test.CreateTestRequest(testApp.routes(), http.MethodDelete, "/locations/%s", fixtureData.Locations[0].ID)

	mt := test.CreateMatrixTester(t, tReq)
	mt.AddTestCaseCookieStatusCode(nil, http.StatusUnauthorized)
	mt.AddTestCaseCookieStatusCode(tokens.DefaultAccessToken, http.StatusForbidden)

	mt.Do(t)
}
