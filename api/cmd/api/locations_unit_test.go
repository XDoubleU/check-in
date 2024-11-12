package main

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	"check-in/api/internal/models"
)

func TestNormalizeName(t *testing.T) {
	//nolint:exhaustruct // other fields are optional
	location1 := models.Location{
		Name: "Test name $14",
	}

	//nolint:exhaustruct // other fields are optional
	location2 := models.Location{
		Name: " Test name $14",
	}

	//nolint:exhaustruct // other fields are optional
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
	createdAt := time.Now().UTC()
	noCheckIns := []*models.CheckIn{}
	fiveCheckIns := generateCheckIns(5, 10, createdAt)
	tenCheckIns := generateCheckIns(10, 10, createdAt)

	// Case 1: check-ins yesterday
	//nolint:exhaustruct // other fields are optional
	location1 := models.Location{
		Capacity: 15,
	}
	location1.SetCheckInRelatedFields(noCheckIns, fiveCheckIns)
	assert.EqualValues(t, location1.Available, 15)
	assert.EqualValues(t, location1.Capacity, 15)
	assert.EqualValues(t, location1.AvailableYesterday, 5)
	assert.EqualValues(t, location1.CapacityYesterday, 10)
	//nolint:exhaustruct // other fields are optional
	assert.Equal(t, location1.YesterdayFullAt, pgtype.Timestamptz{})

	// Case 2: yesterday full
	//nolint:exhaustruct // other fields are optional
	location2 := models.Location{
		Capacity: 15,
	}
	location2.SetCheckInRelatedFields(noCheckIns, tenCheckIns)
	assert.EqualValues(t, location2.Available, 15)
	assert.EqualValues(t, location2.Capacity, 15)
	assert.EqualValues(t, location2.AvailableYesterday, 0)
	assert.EqualValues(t, location2.CapacityYesterday, 10)
	assert.Equal(t, location2.YesterdayFullAt.Time, createdAt)

	// Case 3: yesterday no check-ins, today check-ins
	//nolint:exhaustruct // other fields are optional
	location3 := models.Location{
		Capacity: 15,
	}
	location3.SetCheckInRelatedFields(fiveCheckIns, noCheckIns)
	assert.EqualValues(t, location3.Available, 10)
	assert.EqualValues(t, location3.Capacity, 15)
	assert.EqualValues(t, location3.AvailableYesterday, 10)
	assert.EqualValues(t, location3.CapacityYesterday, 10)
	//nolint:exhaustruct // other fields are optional
	assert.Equal(t, location3.YesterdayFullAt, pgtype.Timestamptz{})

	// Case 4: yesterday no check-ins, today no check-ins (yet?)
	//nolint:exhaustruct // other fields are optional
	location4 := models.Location{
		Capacity: 15,
	}
	location4.SetCheckInRelatedFields(noCheckIns, noCheckIns)
	assert.EqualValues(t, location4.Available, 15)
	assert.EqualValues(t, location4.Capacity, 15)
	assert.EqualValues(t, location4.AvailableYesterday, 15)
	assert.EqualValues(t, location4.CapacityYesterday, 15)
	//nolint:exhaustruct // other fields are optional
	assert.Equal(t, location4.YesterdayFullAt, pgtype.Timestamptz{})
}

func generateCheckIns(amount int, capacity int, createdAt time.Time) []*models.CheckIn {
	checkIns := []*models.CheckIn{}

	for i := 0; i < amount; i++ {
		//nolint:exhaustruct // other fields are optional
		checkIn := models.CheckIn{
			Capacity: int64(capacity),
			//nolint:exhaustruct // other fields are optional
			CreatedAt: pgtype.Timestamptz{
				Time: createdAt,
			},
		}
		checkIns = append(checkIns, &checkIn)
	}

	return checkIns
}
