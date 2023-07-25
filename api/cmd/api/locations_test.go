package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
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

func TestNormalizeName(t *testing.T) {
	// t.Parallel()

	location1 := models.Location{
		Name: "Test name $14",
	}

	location2 := models.Location{
		Name: " Test name $14",
	}

	location3 := models.Location{
		Name: "Test name $14 ",
	}

	err1 := location1.NormalizeName()
	if err1 != nil {
		panic(err1)
	}

	err2 := location2.NormalizeName()
	if err2 != nil {
		panic(err2)
	}

	err3 := location3.NormalizeName()
	if err3 != nil {
		panic(err3)
	}

	assert.Equal(t, location1.NormalizedName, "test-name-14")
	assert.Equal(t, location2.NormalizedName, "test-name-14")
	assert.Equal(t, location3.NormalizedName, "test-name-14")
}

func TestYesterdayFullAt(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	fullTime := time.Now().AddDate(0, 0, -1)
	for i := 0; i < int(fixtureData.DefaultLocation.Capacity); i++ {
		query := `
			INSERT INTO check_ins (location_id, school_id, capacity, created_at)
			VALUES ($1, $2, $3, $4)
		`

		_, _ = testEnv.TestTx.Exec(
			context.Background(),
			query,
			fixtureData.DefaultLocation.ID,
			1,
			fixtureData.DefaultLocation.Capacity,
			fullTime,
		)
	}

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/current-location", nil)
	req.AddCookie(&tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData models.Location
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rsData.YesterdayFullAt.Valid, true)
	assert.Equal(t, rsData.YesterdayFullAt.Time.Day(), fullTime.Day())
	assert.Equal(t, rsData.YesterdayFullAt.Time.Hour(), fullTime.Hour())
	assert.Equal(t, rsData.YesterdayFullAt.Time.Minute(), fullTime.Minute())
}

func TestGetCheckInsLocationRangeRaw(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	startDate := time.Now().AddDate(0, 0, -1)
	startDate = *helpers.StartOfDay(&startDate)

	endDate := time.Now().AddDate(0, 0, 1)
	endDate = *helpers.StartOfDay(&endDate)

	users := []http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/locations/"+fixtureData.DefaultLocation.ID+"/checkins/range",
			nil,
		)
		req.AddCookie(&user)

		query := req.URL.Query()
		query.Add("startDate", startDate.Format(constants.DateFormat))
		query.Add("endDate", endDate.Format(constants.DateFormat))
		query.Add("returnType", "raw")

		req.URL.RawQuery = query.Encode()

		rs, _ := ts.Client().Do(req)

		var rsData map[int64]dtos.CheckInsLocationEntryRaw
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		assert.Equal(t, rsData[startDate.Unix()].Capacity, 0)
		assert.Equal(t, rsData[startDate.Unix()].Schools["Andere"], 0)

		assert.Equal(
			t,
			rsData[startDate.AddDate(0, 0, 1).Unix()].Capacity,
			fixtureData.DefaultLocation.Capacity,
		)
		assert.Equal(t, rsData[startDate.AddDate(0, 0, 1).Unix()].Schools["Andere"], 5)

		assert.Equal(t, rsData[endDate.Unix()].Capacity, 0)
		assert.Equal(t, rsData[endDate.Unix()].Schools["Andere"], 0)
	}
}

func TestGetCheckInsLocationRangeCsv(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	startDate := time.Now().Format(constants.DateFormat)
	endDate := time.Now().AddDate(0, 0, 1).Format(constants.DateFormat)

	users := []http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/locations/"+fixtureData.DefaultLocation.ID+"/checkins/range",
			nil,
		)
		req.AddCookie(&user)

		query := req.URL.Query()
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
		ts.URL+"/locations/"+id.String()+"/checkins/range",
		nil,
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	query := req.URL.Query()
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
		rsData.Message.(string),
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
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/range",
		nil,
	)
	req.AddCookie(&tokens.DefaultAccessToken)

	query := req.URL.Query()
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
		rsData.Message.(string),
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
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/range",
		nil,
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	query := req.URL.Query()
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
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/range",
		nil,
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	query := req.URL.Query()
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
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/range",
		nil,
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	query := req.URL.Query()
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
		ts.URL+"/locations/8000/checkins/range",
		nil,
	)
	req.AddCookie(&tokens.DefaultAccessToken)

	query := req.URL.Query()
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
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/range",
		nil,
	)

	rs1, _ := ts.Client().Do(req1)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
}

func TestGetCheckInsLocationDayRaw(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	date := time.Now()

	users := []http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/locations/"+fixtureData.DefaultLocation.ID+"/checkins/day",
			nil,
		)
		req.AddCookie(&user)

		query := req.URL.Query()
		query.Add("date", date.Format(constants.DateFormat))
		query.Add("returnType", "raw")

		req.URL.RawQuery = query.Encode()

		rs, _ := ts.Client().Do(req)

		var rsData map[int64]dtos.CheckInsLocationEntryRaw
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)

		var checkInDate int64
		for k := range rsData {
			checkInDate = k
			break
		}

		assert.Equal(
			t,
			rsData[checkInDate].Capacity,
			fixtureData.DefaultLocation.Capacity,
		)
		assert.Equal(t, rsData[checkInDate].Schools["Andere"], 5)
	}
}

func TestGetCheckInsLocationDayCsv(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	date := time.Now().Format(constants.DateFormat)

	users := []http.Cookie{
		tokens.AdminAccessToken,
		tokens.ManagerAccessToken,
		tokens.DefaultAccessToken,
	}

	for _, user := range users {
		req, _ := http.NewRequest(
			http.MethodGet,
			ts.URL+"/locations/"+fixtureData.DefaultLocation.ID+"/checkins/day",
			nil,
		)
		req.AddCookie(&user)

		query := req.URL.Query()
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
		ts.URL+"/locations/"+id.String()+"/checkins/day",
		nil,
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	query := req.URL.Query()
	query.Add("date", date)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
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
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/day",
		nil,
	)
	req.AddCookie(&tokens.DefaultAccessToken)

	query := req.URL.Query()
	query.Add("date", date)
	query.Add("returnType", "raw")

	req.URL.RawQuery = query.Encode()

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
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
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/day",
		nil,
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	query := req.URL.Query()
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
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/day",
		nil,
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	query := req.URL.Query()
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
		ts.URL+"/locations/8000/checkins/day",
		nil,
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	query := req.URL.Query()
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
		ts.URL+"/locations/"+fixtureData.Locations[0].ID+"/checkins/day",
		nil,
	)

	rs1, _ := ts.Client().Do(req1)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
}

func TestGetLocationLoggedInUser(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/current-location", nil)
	req.AddCookie(&tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData models.Location
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rsData.ID, fixtureData.DefaultLocation.ID)
	assert.Equal(t, rsData.Name, fixtureData.DefaultLocation.Name)
	assert.Equal(t, rsData.NormalizedName, fixtureData.DefaultLocation.NormalizedName)
	assert.Equal(t, rsData.Available, fixtureData.DefaultLocation.Available)
	assert.Equal(t, rsData.Capacity, fixtureData.DefaultLocation.Capacity)
	assert.Equal(t, rsData.YesterdayFullAt, fixtureData.DefaultLocation.YesterdayFullAt)
	assert.Equal(t, rsData.UserID, fixtureData.DefaultLocation.UserID)
}

func TestGetLocationLoggedInUserAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodGet, ts.URL+"/current-location", nil)

	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/current-location", nil)
	req2.AddCookie(&tokens.ManagerAccessToken)

	req3, _ := http.NewRequest(http.MethodGet, ts.URL+"/current-location", nil)
	req3.AddCookie(&tokens.AdminAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)
	rs3, _ := ts.Client().Do(req3)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
	assert.Equal(t, rs3.StatusCode, http.StatusForbidden)
}

func TestGetPaginatedLocationsDefaultPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	users := []http.Cookie{tokens.AdminAccessToken, tokens.ManagerAccessToken}

	for _, user := range users {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations", nil)
		req.AddCookie(&user)

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
	}
}

func TestGetPaginatedLocationsSpecificPage(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations?page=2", nil)
	req.AddCookie(&tokens.ManagerAccessToken)

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
}

func TestGetPaginatedLocationsPageZero(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations?page=0", nil)
	req.AddCookie(&tokens.ManagerAccessToken)

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
	req2.AddCookie(&tokens.DefaultAccessToken)

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

	users := []http.Cookie{
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
		req.AddCookie(&user)

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
		assert.Equal(
			t,
			rsData.YesterdayFullAt,
			fixtureData.DefaultLocation.YesterdayFullAt,
		)
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
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
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
	req.AddCookie(&tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("location with id '%s' doesn't exist", fixtureData.Locations[0].ID),
	)
}

func TestGetLocationNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/locations/8000", nil)
	req.AddCookie(&tokens.ManagerAccessToken)

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

	users := []http.Cookie{tokens.AdminAccessToken, tokens.ManagerAccessToken}

	for i, user := range users {
		unique := fmt.Sprintf("test%d", i)

		data := dtos.CreateLocationDto{
			Name:     unique,
			Capacity: 10,
			Username: unique,
			Password: "testpassword",
		}

		body, _ := json.Marshal(data)

		req, _ := http.NewRequest(
			http.MethodPost,
			ts.URL+"/locations",
			bytes.NewReader(body),
		)
		req.AddCookie(&user)

		rs, _ := ts.Client().Do(req)

		var rsData models.Location
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusCreated)
		assert.IsUUID(t, rsData.ID)
		assert.Equal(t, rsData.Name, data.Name)
		assert.Equal(t, rsData.NormalizedName, data.Name)
		assert.Equal(t, rsData.Available, 10)
		assert.Equal(t, rsData.Capacity, 10)
		assert.Equal(t, rsData.YesterdayFullAt.Valid, false)
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
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/locations",
		bytes.NewReader(body),
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
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
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/locations",
		bytes.NewReader(body),
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
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
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		http.MethodPost,
		ts.URL+"/locations",
		bytes.NewReader(body),
	)
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
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
	req.AddCookie(&tokens.ManagerAccessToken)

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

func TestCreateLocationAccess(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req1, _ := http.NewRequest(http.MethodPost, ts.URL+"/locations", nil)

	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/locations", nil)
	req2.AddCookie(&tokens.DefaultAccessToken)

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

	users := []http.Cookie{tokens.AdminAccessToken, tokens.ManagerAccessToken}

	for i, user := range users {
		unique := fmt.Sprintf("test%d", i)
		name, username, password := unique, unique, "testpassword"
		var capacity int64 = 3
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
		req.AddCookie(&user)

		rs, _ := ts.Client().Do(req)

		var rsData models.Location
		_ = helpers.ReadJSON(rs.Body, &rsData)

		assert.Equal(t, rs.StatusCode, http.StatusOK)
		assert.Equal(t, rsData.ID, fixtureData.Locations[0].ID)
		assert.Equal(t, rsData.Name, *data.Name)
		assert.Equal(t, rsData.NormalizedName, *data.Name)
		assert.Equal(t, rsData.Available, 0)
		assert.Equal(t, rsData.Capacity, *data.Capacity)
		assert.Equal(t, rsData.YesterdayFullAt.Valid, false)
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
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
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
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
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
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusConflict)
	assert.Equal(
		t,
		rsData.Message.(string),
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
	req.AddCookie(&tokens.ManagerAccessToken)

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
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
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
	req.AddCookie(&tokens.DefaultAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
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
	req.AddCookie(&tokens.ManagerAccessToken)

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

	users := []http.Cookie{tokens.AdminAccessToken, tokens.ManagerAccessToken}

	for i, user := range users {
		req, _ := http.NewRequest(
			http.MethodDelete,
			ts.URL+"/locations/"+fixtureData.Locations[i].ID,
			nil,
		)
		req.AddCookie(&user)

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
			rsData.YesterdayFullAt,
			fixtureData.Locations[i].YesterdayFullAt,
		)
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
	req.AddCookie(&tokens.ManagerAccessToken)

	rs, _ := ts.Client().Do(req)

	var rsData dtos.ErrorDto
	_ = helpers.ReadJSON(rs.Body, &rsData)

	assert.Equal(t, rs.StatusCode, http.StatusNotFound)
	assert.Equal(
		t,
		rsData.Message.(string),
		fmt.Sprintf("location with id '%s' doesn't exist", id.String()),
	)
}

func TestDeleteLocationNotUUID(t *testing.T) {
	testEnv, testApp := setupTest(t, mainTestEnv)
	defer tests.TeardownSingle(testEnv)

	ts := httptest.NewTLSServer(testApp.routes())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/locations/8000", nil)
	req.AddCookie(&tokens.ManagerAccessToken)

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
	req2.AddCookie(&tokens.DefaultAccessToken)

	rs1, _ := ts.Client().Do(req1)
	rs2, _ := ts.Client().Do(req2)

	assert.Equal(t, rs1.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, rs2.StatusCode, http.StatusForbidden)
}
