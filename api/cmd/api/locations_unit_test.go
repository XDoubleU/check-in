package main

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"check-in/api/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeName(t *testing.T) {
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

func TestSetCheckInRelatedFields(t *testing.T) {
	createdAt := time.Now()
	noCheckIns := []*models.CheckIn{}
	fiveCheckIns := generateCheckIns(5, 10, createdAt)
	tenCheckIns := generateCheckIns(10, 10, createdAt)

	// Case 1: check-ins yesterday
	location1 := models.Location{
		Capacity: 15,
	}
	location1.SetCheckInRelatedFields(noCheckIns, fiveCheckIns)
	assert.Equal(t, location1.Available, int64(15))
	assert.Equal(t, location1.Capacity, int64(15))
	assert.Equal(t, location1.AvailableYesterday, int64(5))
	assert.Equal(t, location1.CapacityYesterday, int64(10))
	assert.Equal(t, location1.YesterdayFullAt, pgtype.Timestamptz{})

	// Case 2: yesterday full
	location2 := models.Location{
		Capacity: 15,
	}
	location2.SetCheckInRelatedFields(noCheckIns, tenCheckIns)
	assert.Equal(t, location2.Available, int64(15))
	assert.Equal(t, location2.Capacity, int64(15))
	assert.Equal(t, location2.AvailableYesterday, int64(0))
	assert.Equal(t, location2.CapacityYesterday, int64(10))
	assert.Equal(t, location2.YesterdayFullAt.Time, createdAt)

	// Case 3: yesterday no check-ins, today check-ins
	location3 := models.Location{
		Capacity: 15,
	}
	location3.SetCheckInRelatedFields(fiveCheckIns, noCheckIns)
	assert.Equal(t, location3.Available, int64(10))
	assert.Equal(t, location3.Capacity, int64(15))
	assert.Equal(t, location3.AvailableYesterday, int64(10))
	assert.Equal(t, location3.CapacityYesterday, int64(10))
	assert.Equal(t, location3.YesterdayFullAt, pgtype.Timestamptz{})

	// Case 4: yesterday no check-ins, today no check-ins (yet?)
	location4 := models.Location{
		Capacity: 15,
	}
	location4.SetCheckInRelatedFields(noCheckIns, noCheckIns)
	assert.Equal(t, location4.Available, int64(15))
	assert.Equal(t, location4.Capacity, int64(15))
	assert.Equal(t, location4.AvailableYesterday, int64(15))
	assert.Equal(t, location4.CapacityYesterday, int64(15))
	assert.Equal(t, location4.YesterdayFullAt, pgtype.Timestamptz{})
}

func generateCheckIns(amount int, capacity int, createdAt time.Time) []*models.CheckIn {
	checkIns := []*models.CheckIn{}

	for i := 0; i < amount; i++ {
		checkIn := models.CheckIn{
			Capacity: int64(capacity),
			CreatedAt: pgtype.Timestamptz{
				Time: createdAt,
			},
		}
		checkIns = append(checkIns, &checkIn)
	}

	return checkIns
}
