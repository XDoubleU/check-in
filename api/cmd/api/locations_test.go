package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/test"
	timetools "github.com/xdoubleu/essentia/pkg/time"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func TestYesterdayFullAt(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	loc, _ := time.LoadLocation("Europe/Brussels")

	now := time.Now().In(loc).AddDate(0, 0, -1)

	for i := 0; i < int(testEnv.Fixtures.DefaultLocation.Capacity); i++ {
		query := `
			INSERT INTO check_ins 
			(location_id, school_id, capacity, created_at)
			VALUES ($1, $2, $3, $4)
		`

		_, err := testEnv.tx.Exec(
			context.Background(),
			query,
			testEnv.Fixtures.DefaultLocation.ID,
			1,
			testEnv.Fixtures.DefaultLocation.Capacity,
			now,
		)
		if err != nil {
			panic(err)
		}
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/locations/%s",
		testEnv.Fixtures.DefaultLocation.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	var rsData models.Location
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 0, rsData.AvailableYesterday)
	assert.Equal(t, testEnv.Fixtures.DefaultLocation.Capacity, rsData.CapacityYesterday)
	assert.Equal(t, true, rsData.YesterdayFullAt.Valid)
	assert.Equal(t, now.Day(), rsData.YesterdayFullAt.Time.Day())
	assert.Equal(t, now.Hour(), rsData.YesterdayFullAt.Time.Hour())
}

func TestGetCheckInsLocationRangeRawSingle(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	testEnv.createCheckIns(testEnv.Fixtures.DefaultLocation, int64(1), 10)

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	startDate := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, utc)
	startDate = timetools.StartOfDay(startDate)

	endDate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, utc)
	endDate = timetools.StartOfDay(endDate)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
		testEnv.Fixtures.Tokens.DefaultAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/all-locations/checkins/range",
		)
		tReq.AddCookie(user)

		tReq.SetQuery(map[string]string{
			"ids":        testEnv.Fixtures.DefaultLocation.ID,
			"startDate":  startDate.Format(constants.DateFormat),
			"endDate":    endDate.Format(constants.DateFormat),
			"returnType": "raw",
		})

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		assert.Equal(t, 0, rsData[startDate.Format(time.RFC3339)].Capacities.Len())

		value, present := rsData[startDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, 0, value)
		assert.Equal(t, true, present)

		capacity, _ := rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].Capacities.Get(
			testEnv.Fixtures.DefaultLocation.ID,
		)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.Capacity,
			capacity,
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

func TestGetCheckInsLocationRangeRawMultiple(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	testEnv.createCheckIns(testEnv.Fixtures.DefaultLocation, int64(1), 10)
	testEnv.createCheckIns(location, int64(1), 10)

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	startDate := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, utc)
	startDate = timetools.StartOfDay(startDate)

	endDate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, utc)
	endDate = timetools.StartOfDay(endDate)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/all-locations/checkins/range",
		)
		tReq.AddCookie(user)

		id := fmt.Sprintf(
			"%s,%s",
			testEnv.Fixtures.DefaultLocation.ID,
			location.ID,
		)

		tReq.SetQuery(map[string]string{
			"ids":        id,
			"startDate":  startDate.Format(constants.DateFormat),
			"endDate":    endDate.Format(constants.DateFormat),
			"returnType": "raw",
		})

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		assert.Equal(t, 0, rsData[startDate.Format(time.RFC3339)].Capacities.Len())

		value, present := rsData[startDate.Format(time.RFC3339)].Schools.Get(
			"Andere",
		)
		assert.Equal(t, 0, value)
		assert.Equal(t, true, present)

		capacity0, _ := rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].
			Capacities.Get(
			testEnv.Fixtures.DefaultLocation.ID,
		)
		capacity1, _ := rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].
			Capacities.Get(
			location.ID,
		)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.Capacity,
			capacity0,
		)

		assert.Equal(
			t,
			location.Capacity,
			capacity1,
		)

		value, present = rsData[startDate.AddDate(0, 0, 1).Format(time.RFC3339)].
			Schools.Get(
			"Andere",
		)
		assert.Equal(t, 20, value)
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
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	startDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 2).Format(constants.DateFormat)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
		testEnv.Fixtures.Tokens.DefaultAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/all-locations/checkins/range",
		)
		tReq.AddCookie(user)

		tReq.SetQuery(map[string]string{
			"ids":        testEnv.Fixtures.DefaultLocation.ID,
			"startDate":  startDate,
			"endDate":    endDate,
			"returnType": "csv",
		})

		var rsData [][]string
		rs := tReq.Do(t, &rsData, httptools.ReadCSV)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, "text/csv", rs.Header.Get("content-type"))
		//todo check rsData
	}
}

func TestGetCheckInsLocationRangeNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	startDate := time.Now().Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	id, _ := uuid.NewUUID()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/range",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        id.String(),
		"startDate":  startDate,
		"endDate":    endDate,
		"returnType": "raw",
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetCheckInsLocationRangeNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	startDate := time.Now().Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/range",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        location.ID,
		"startDate":  startDate,
		"endDate":    endDate,
		"returnType": "raw",
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", location.ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetCheckInsLocationRangeStartDateMissing(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/range",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        location.ID,
		"endDate":    endDate,
		"returnType": "raw",
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'startDate'", rsData.Message.(string))
}

func TestGetCheckInsLocationRangeEndDateMissing(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	startDate := time.Now().Format(constants.DateFormat)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/range",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        location.ID,
		"startDate":  startDate,
		"returnType": "raw",
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'endDate'", rsData.Message.(string))
}

func TestGetCheckInsLocationRangeReturnTypeMissing(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	startDate := time.Now().Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/range",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":       location.ID,
		"startDate": startDate,
		"endDate":   endDate,
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'returnType'", rsData.Message.(string))
}

func TestGetCheckInsLocationRangeNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	startDate := time.Now().Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/range",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        "8000",
		"startDate":  startDate,
		"endDate":    endDate,
		"returnType": "raw",
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestGetCheckInsLocationRangeAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/range",
	)

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReq, test.NewCaseResponse(http.StatusUnauthorized))

	mt.Do(t)
}

func TestGetCheckInsLocationDayRawSingle(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	amount := 10
	testEnv.createCheckIns(testEnv.Fixtures.DefaultLocation, int64(1), amount)

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utc)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
		testEnv.Fixtures.Tokens.DefaultAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/all-locations/checkins/day",
		)
		tReq.AddCookie(user)

		tReq.SetQuery(map[string]string{
			"ids":        testEnv.Fixtures.DefaultLocation.ID,
			"date":       date.Format(constants.DateFormat),
			"returnType": "raw",
		})

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		var checkInDate string
		for k := range rsData {
			checkInDate = k
			break
		}

		capacity, _ := rsData[checkInDate].Capacities.Get(
			testEnv.Fixtures.DefaultLocation.ID,
		)
		assert.Equal(t, capacity, int64(20))

		value, present := rsData[checkInDate].Schools.Get("Andere")
		assert.Equal(t, amount, value)
		assert.Equal(t, true, present)
	}
}

func TestGetCheckInsLocationDayRawMultiple(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	amount := 10
	testEnv.createCheckIns(testEnv.Fixtures.DefaultLocation, int64(1), amount)
	testEnv.createCheckIns(location, int64(1), amount)

	loc, _ := time.LoadLocation("Europe/Brussels")
	utc, _ := time.LoadLocation("UTC")

	now := time.Now().In(loc)

	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utc)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/all-locations/checkins/day",
		)
		tReq.AddCookie(user)

		id := fmt.Sprintf(
			"%s,%s",
			testEnv.Fixtures.DefaultLocation.ID,
			location.ID,
		)
		tReq.SetQuery(map[string]string{
			"ids":        id,
			"date":       date.Format(constants.DateFormat),
			"returnType": "raw",
		})

		var rsData map[string]dtos.CheckInsLocationEntryRaw
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		var checkInDate string
		for k := range rsData {
			checkInDate = k
			break
		}

		capacity0, _ := rsData[checkInDate].Capacities.Get(
			testEnv.Fixtures.DefaultLocation.ID,
		)
		capacity1, _ := rsData[checkInDate].Capacities.Get(location.ID)
		assert.Equal(t, int64(20), capacity0)
		assert.Equal(
			t,
			location.Capacity,
			capacity1,
		)

		value, present := rsData[checkInDate].Schools.Get("Andere")
		assert.Equal(t, 2*amount, value)
		assert.Equal(t, true, present)
	}
}

func TestGetCheckInsLocationDayCSV(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	date := time.Now().Format(constants.DateFormat)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
		testEnv.Fixtures.Tokens.DefaultAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/all-locations/checkins/day",
		)
		tReq.AddCookie(user)

		tReq.SetQuery(map[string]string{
			"ids":        testEnv.Fixtures.DefaultLocation.ID,
			"date":       date,
			"returnType": "csv",
		})

		var rsData [][]string
		rs := tReq.Do(t, &rsData, httptools.ReadCSV)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, "text/csv", rs.Header.Get("content-type"))
		//todo check rsData
	}
}

func TestGetCheckInsLocationDayNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	date := time.Now().Format(constants.DateFormat)

	id, _ := uuid.NewUUID()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/day",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        id.String(),
		"date":       date,
		"returnType": "raw",
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetCheckInsLocationDayNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	date := time.Now().Format(constants.DateFormat)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/day",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        location.ID,
		"date":       date,
		"returnType": "raw",
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", location.ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetCheckInsLocationDateMissing(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/day",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        location.ID,
		"returnType": "raw",
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'date'", rsData.Message.(string))
}

func TestGetCheckInsLocationReturnTypeMissing(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	date := time.Now().Format(constants.DateFormat)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/day",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":  location.ID,
		"date": date,
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(t, "missing query param 'returnType'", rsData.Message.(string))
}

func TestGetCheckInsLocationDayNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	date := time.Now().Format(constants.DateFormat)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/day",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	tReq.SetQuery(map[string]string{
		"ids":        "8000",
		"date":       date,
		"returnType": "raw",
	})

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestGetCheckInsLocationDayAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations/checkins/day",
	)

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReq, test.NewCaseResponse(http.StatusUnauthorized))

	mt.Do(t)
}

func TestGetAllCheckInsToday(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	amount := 5
	testEnv.createCheckIns(testEnv.Fixtures.DefaultLocation, int64(1), amount)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
		testEnv.Fixtures.Tokens.DefaultAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/locations/%s/checkins",
			testEnv.Fixtures.DefaultLocation.ID,
		)
		tReq.AddCookie(user)

		var rsData []dtos.CheckInDto
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		loc, _ := time.LoadLocation("Europe/Brussels")

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, amount, len(rsData))
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.ID, rsData[0].LocationID)
		assert.Equal(t, "Andere", rsData[0].SchoolName)
		assert.Equal(
			t,
			time.Now().In(loc).Format(constants.DateFormat),
			rsData[0].CreatedAt.Time.Format(constants.DateFormat),
		)
	}
}

func TestGetAllCheckInsTodayNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	id, _ := uuid.NewUUID()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/locations/%s/checkins",
		id.String(),
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetAllCheckInsTodayNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/locations/%s/checkins",
		location.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", location.ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetAllCheckInsTodayNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/locations/8000/checkins",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestGetCheckInsTodayAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/locations/%s/checkins",
		location.ID,
	)

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReq, test.NewCaseResponse(http.StatusUnauthorized))

	mt.Do(t)
}

func TestDeleteCheckIn(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	checkIns := testEnv.createCheckIns(testEnv.Fixtures.DefaultLocation, 1, 10)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
	}

	for i, user := range users {
		id := strconv.FormatInt(
			checkIns[i].ID,
			10,
		)
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodDelete,
			"/locations/%s/checkins/%s",
			testEnv.Fixtures.DefaultLocation.ID,
			id,
		)
		tReq.AddCookie(user)

		var rsData dtos.CheckInDto
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		loc, _ := time.LoadLocation("Europe/Brussels")

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, checkIns[i].ID, rsData.ID)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.ID, rsData.LocationID)
		assert.Equal(t, "Andere", rsData.SchoolName)
		assert.Equal(
			t,
			time.Now().In(loc).Format(constants.DateFormat),
			rsData.CreatedAt.Time.Format(constants.DateFormat),
		)
	}
}

func TestDeleteCheckInLocationNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	id, _ := uuid.NewUUID()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/locations/%s/checkins/1",
		id.String(),
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteCheckInCheckInNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/locations/%s/checkins/8000",
		testEnv.Fixtures.DefaultLocation.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		"checkIn with id '8000' doesn't exist",
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteCheckInNotToday(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	query := `
			INSERT INTO check_ins 
			(location_id, school_id, capacity, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`

	var checkIn models.CheckIn
	err := testEnv.tx.QueryRow(
		context.Background(),
		query,
		testEnv.Fixtures.DefaultLocation.ID,
		1,
		testEnv.Fixtures.DefaultLocation.Capacity,
		time.Now().AddDate(0, 0, -1),
	).Scan(&checkIn.ID)
	if err != nil {
		panic(err)
	}

	id := strconv.FormatInt(
		checkIn.ID,
		10,
	)
	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/locations/%s/checkins/%s",
		testEnv.Fixtures.DefaultLocation.ID,
		id,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Equal(
		t,
		"checkIn didn't occur today and thus can't be deleted",
		rsData.Message,
	)
}

func TestDeleteCheckInNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/locations/800O/checkins/1",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestDeleteCheckInAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/locations/%s/checkins/1",
		location.ID,
	)

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestGetPaginatedLocationsDefaultPage(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	amount := 20
	testEnv.createLocations(amount)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/locations")
		tReq.AddCookie(user)

		var rsData dtos.PaginatedLocationsDto
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		assert.EqualValues(t, 1, rsData.Pagination.Current)
		assert.EqualValues(
			t,
			math.Ceil(float64(amount)/3),
			rsData.Pagination.Total,
		)
		assert.Equal(t, 3, len(rsData.Data))

		assert.Equal(t, testEnv.Fixtures.DefaultLocation.ID, rsData.Data[0].ID)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.Name, rsData.Data[0].Name)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.NormalizedName,
			rsData.Data[0].NormalizedName,
		)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.Available,
			rsData.Data[0].Available,
		)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.Capacity,
			rsData.Data[0].Capacity,
		)
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
			testEnv.Fixtures.DefaultLocation.YesterdayFullAt,
			rsData.Data[0].YesterdayFullAt,
		)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.TimeZone,
			rsData.Data[0].TimeZone,
		)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.UserID, rsData.Data[0].UserID)
	}
}

func TestGetPaginatedLocationsSpecificPage(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	amount := 20
	locations := testEnv.createLocations(amount)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/locations")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetQuery(map[string]string{
		"page": "2",
	})

	var rsData dtos.PaginatedLocationsDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	assert.EqualValues(t, 2, rsData.Pagination.Current)
	assert.EqualValues(
		t,
		math.Ceil(float64(amount)/3),
		rsData.Pagination.Total,
	)
	assert.Equal(t, 3, len(rsData.Data))

	assert.Equal(t, locations[10].ID, rsData.Data[0].ID)
	assert.Equal(t, locations[10].Name, rsData.Data[0].Name)
	assert.Equal(
		t,
		locations[10].NormalizedName,
		rsData.Data[0].NormalizedName,
	)
	assert.Equal(t, locations[10].Available, rsData.Data[0].Available)
	assert.Equal(t, locations[10].Capacity, rsData.Data[0].Capacity)
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
		locations[10].YesterdayFullAt,
		rsData.Data[0].YesterdayFullAt,
	)
	assert.Equal(t, locations[10].TimeZone, rsData.Data[0].TimeZone)
	assert.Equal(t, locations[10].UserID, rsData.Data[0].UserID)
}

func TestGetPaginatedLocationsPageFull(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	amount := 20
	testEnv.createLocations(amount)

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/locations")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	test.PaginatedEndpointTester(
		t,
		tReq,
		"page",
		int(math.Ceil(float64(amount)/4)),
	)
}

func TestGetPaginatedLocationsAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(testApp.routes(), http.MethodGet, "/locations")

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestGetAllLocations(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	testEnv.createLocations(20)
	amount, err := testEnv.services.Locations.GetTotalCount(context.Background())
	require.Nil(t, err)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/all-locations",
		)
		tReq.AddCookie(user)

		var rsData []models.Location
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		assert.Equal(t, http.StatusOK, rs.StatusCode)

		assert.EqualValues(t, *amount, len(rsData))
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.ID, rsData[0].ID)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.Name, rsData[0].Name)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.NormalizedName,
			rsData[0].NormalizedName,
		)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.Available, rsData[0].Available)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.Capacity, rsData[0].Capacity)
		assert.NotEqual(t, 0, rsData[0].AvailableYesterday)
		assert.NotEqual(t, 0, rsData[0].CapacityYesterday)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.YesterdayFullAt,
			rsData[0].YesterdayFullAt,
		)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.TimeZone, rsData[0].TimeZone)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.UserID, rsData[0].UserID)
	}
}

func TestGetAllLocationsAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/all-locations",
	)

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestGetLocation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
		testEnv.Fixtures.Tokens.DefaultAccessToken,
	}

	for _, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodGet,
			"/locations/%s",
			testEnv.Fixtures.DefaultLocation.ID,
		)
		tReq.AddCookie(user)

		var rsData models.Location
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.ID, rsData.ID)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.Name, rsData.Name)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.NormalizedName,
			rsData.NormalizedName,
		)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.Available, rsData.Available)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.Capacity, rsData.Capacity)
		assert.NotEqual(t, 0, rsData.AvailableYesterday)
		assert.NotEqual(t, 0, rsData.CapacityYesterday)
		assert.Equal(
			t,
			testEnv.Fixtures.DefaultLocation.YesterdayFullAt,
			rsData.YesterdayFullAt,
		)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.TimeZone, rsData.TimeZone)
		assert.Equal(t, testEnv.Fixtures.DefaultLocation.UserID, rsData.UserID)
	}
}

func TestGetLocationNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	id, _ := uuid.NewUUID()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/locations/%s",
		id.String(),
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetLocationNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/locations/%s",
		location.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", location.ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestGetLocationNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/locations/8000",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestGetLocationAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/locations/%s",
		location.ID,
	)

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReq, test.NewCaseResponse(http.StatusUnauthorized))

	mt.Do(t)
}

func TestCreateLocation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
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

		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodPost,
			"/locations",
		)
		tReq.AddCookie(user)

		tReq.SetBody(data)

		var rsData models.Location
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

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
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	data := dtos.CreateLocationDto{
		Name:     testEnv.Fixtures.DefaultLocation.Name,
		Capacity: 10,
		Username: "test",
		Password: "testpassword",
		TimeZone: "Europe/Brussels",
	}

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/locations")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with name '%s' already exists", data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestCreateLocationNormalizedNameExists(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	data := dtos.CreateLocationDto{
		Name:     fmt.Sprintf("$%s$", testEnv.Fixtures.DefaultLocation.Name),
		Capacity: 10,
		Username: "test",
		Password: "testpassword",
		TimeZone: "Europe/Brussels",
	}

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/locations")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with name '%s' already exists", data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestCreateLocationUserNameExists(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	data := dtos.CreateLocationDto{
		Name:     "test",
		Capacity: 10,
		Username: testEnv.Fixtures.DefaultUser.Username,
		Password: "testpassword",
		TimeZone: "Europe/Brussels",
	}

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/locations")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with username '%s' already exists", data.Username),
		rsData.Message.(map[string]interface{})["username"].(string),
	)
}

func TestCreateLocationFailValidation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(testApp.routes(), http.MethodPost, "/locations")
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	mt := test.CreateMatrixTester()

	tReq1 := tReq.Copy()
	tReq1.SetBody(dtos.CreateLocationDto{
		Name:     "test",
		Capacity: -1,
		Username: "test",
		Password: "testpassword",
		TimeZone: "Europe/Brussels",
	})

	tRes1 := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes1.SetBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"capacity": "must be greater than 0",
		}),
	)

	mt.AddTestCase(tReq1, tRes1)

	tReq2 := tReq.Copy()
	tReq2.SetBody(dtos.CreateLocationDto{
		Name:     "test",
		Capacity: 10,
		Username: "test",
		Password: "testpassword",
		TimeZone: "wrong",
	})

	tRes2 := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes2.SetBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"timeZone": "must be a valid IANA value",
		}),
	)

	mt.AddTestCase(tReq2, tRes2)

	mt.Do(t)
}

func TestCreateLocationAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPost,
		"/locations",
	)

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}

func TestUpdateLocation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
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

		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodPatch,
			"/locations/%s",
			location.ID,
		)
		tReq.AddCookie(user)

		tReq.SetBody(data)

		var rsData models.Location
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, location.ID, rsData.ID)
		assert.Equal(t, *data.Name, rsData.Name)
		assert.Equal(t, *data.Name, rsData.NormalizedName)
		assert.EqualValues(t, *data.Capacity, rsData.Available)
		assert.Equal(t, *data.Capacity, rsData.Capacity)
		assert.Equal(
			t,
			*data.Capacity,
			rsData.AvailableYesterday,
		)
		assert.Equal(
			t,
			*data.Capacity,
			rsData.CapacityYesterday,
		)
		assert.Equal(t, false, rsData.YesterdayFullAt.Valid)
		assert.Equal(t, *data.TimeZone, rsData.TimeZone)
		assert.Equal(t, location.UserID, rsData.UserID)
	}
}

func TestUpdateLocationNameExists(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	name, username, password, timeZone := testEnv.Fixtures.DefaultLocation.Name, "test",
		"testpassword", "Europe/Brussels"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/locations/%s",
		location.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with name '%s' already exists", *data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestUpdateLocationNormalizedNameExists(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	name, username, password, timeZone := fmt.Sprintf(
		"$%s$",
		testEnv.Fixtures.DefaultLocation.Name,
	), "test",
		"testpassword", "Europe/Brussels"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/locations/%s",
		location.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with name '%s' already exists", *data.Name),
		rsData.Message.(map[string]interface{})["name"].(string),
	)
}

func TestUpdateLocationUserNameExists(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	name, username, password, timeZone := "test", testEnv.Fixtures.DefaultUser.Username,
		"testpassword", "Europe/Brussels"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/locations/%s",
		location.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusConflict, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("user with username '%s' already exists", *data.Username),
		rsData.Message.(map[string]interface{})["username"].(string),
	)
}

func TestUpdateLocationFailValidation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/locations/%s",
		location.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	mt := test.CreateMatrixTester()

	tReq1 := tReq.Copy()

	name, username, password, timeZone1 := "test", "test", "testpassword", "Europe/Brussels"
	var capacity1 int64 = -1
	tReq1.SetBody(dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity1,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone1,
	})

	tRes1 := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes1.SetBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"capacity": "must be greater than 0",
		}),
	)

	mt.AddTestCase(tReq1, tRes1)

	tReq2 := tReq.Copy()

	timeZone2 := "wrong"
	var capacity2 int64 = 10
	tReq2.SetBody(dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity2,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone2,
	})

	tRes2 := test.NewCaseResponse(http.StatusUnprocessableEntity)
	tRes2.SetBody(
		errortools.NewErrorDto(http.StatusUnprocessableEntity, map[string]interface{}{
			"timeZone": "must be a valid IANA value",
		}),
	)

	mt.AddTestCase(tReq2, tRes2)

	mt.Do(t)
}

func TestUpdateLocationNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	name, username, password, timeZone := "test", "test", "password", "Europe/Brussels"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone,
	}

	id, _ := uuid.NewUUID()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/locations/%s",
		id.String(),
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateLocationNotFoundNotOwner(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	name, username, password, timeZone := "test", "test", "testpassword", "Europe/Brussels"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/locations/%s",
		location.ID,
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", location.ID),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestUpdateLocationNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	name, username, password, timeZone := "test", "test", "password", "Europe/Brussels"
	var capacity int64 = 10
	data := dtos.UpdateLocationDto{
		Name:     &name,
		Capacity: &capacity,
		Username: &username,
		Password: &password,
		TimeZone: &timeZone,
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/locations/8000",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	tReq.SetBody(data)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestUpdateLocationAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPatch,
		"/locations/%s",
		location.ID,
	)

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReq, test.NewCaseResponse(http.StatusUnauthorized))

	mt.Do(t)
}

func TestDeleteLocation(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	locations := testEnv.createLocations(3)

	users := []*http.Cookie{
		testEnv.Fixtures.Tokens.AdminAccessToken,
		testEnv.Fixtures.Tokens.ManagerAccessToken,
	}

	for i, user := range users {
		tReq := test.CreateRequestTester(
			testApp.routes(),
			http.MethodDelete,
			"/locations/%s",
			locations[i].ID,
		)
		tReq.AddCookie(user)

		var rsData models.Location
		rs := tReq.Do(t, &rsData, httptools.ReadJSON)

		assert.Equal(t, http.StatusOK, rs.StatusCode)
		assert.Equal(t, locations[i].ID, rsData.ID)
		assert.Equal(t, locations[i].Name, rsData.Name)
		assert.Equal(t, locations[i].NormalizedName, rsData.NormalizedName)
		assert.Equal(t, locations[i].Available, rsData.Available)
		assert.Equal(t, locations[i].Capacity, rsData.Capacity)
		assert.Equal(
			t,
			locations[i].AvailableYesterday,
			rsData.AvailableYesterday,
		)
		assert.Equal(
			t,
			locations[i].CapacityYesterday,
			rsData.CapacityYesterday,
		)
		assert.Equal(
			t,
			locations[i].YesterdayFullAt,
			rsData.YesterdayFullAt,
		)
		assert.Equal(t, locations[i].TimeZone, rsData.TimeZone)
		assert.Equal(t, locations[i].UserID, rsData.UserID)
	}
}

func TestDeleteLocationNotFound(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	id, _ := uuid.NewUUID()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/locations/%s",
		id.String(),
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
	assert.Equal(
		t,
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
		rsData.Message.(map[string]interface{})["id"].(string),
	)
}

func TestDeleteLocationNotUUID(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/locations/8000",
	)
	tReq.AddCookie(testEnv.Fixtures.Tokens.ManagerAccessToken)

	var rsData errortools.ErrorDto
	rs := tReq.Do(t, &rsData, httptools.ReadJSON)

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
	assert.Contains(t, rsData.Message.(string), "should be a UUID")
}

func TestDeleteLocationAccess(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	location := testEnv.createLocations(1)[0]

	tReqBase := test.CreateRequestTester(
		testApp.routes(),
		http.MethodDelete,
		"/locations/%s",
		location.ID,
	)

	mt := test.CreateMatrixTester()
	mt.AddTestCase(tReqBase, test.NewCaseResponse(http.StatusUnauthorized))

	tReq2 := tReqBase.Copy()
	tReq2.AddCookie(testEnv.Fixtures.Tokens.DefaultAccessToken)

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden))

	mt.Do(t)
}
