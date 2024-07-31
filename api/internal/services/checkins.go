package services

import (
	"context"
	"encoding/json"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	timetools "github.com/xdoubleu/essentia/pkg/time"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type CheckInService struct {
	checkins repositories.CheckInRepository
	schools  SchoolService
}

// todo: refactor
func (service CheckInService) GetCheckInsEntriesDay(
	checkIns []*models.CheckIn,
	schools []*models.School,
) *orderedmap.OrderedMap[string, *dtos.CheckInsLocationEntryRaw] {
	schoolsIDNameMap, _ := service.schools.GetSchoolMaps(schools)

	checkInEntries := orderedmap.New[string, *dtos.CheckInsLocationEntryRaw]()

	_, lastEntrySchoolsMap := service.schools.GetSchoolMaps(schools)
	capacities := orderedmap.New[string, int64]()
	for _, checkIn := range checkIns {
		schoolName := schoolsIDNameMap[checkIn.SchoolID]

		// Used to deep copy schoolsMap
		var schoolsMap dtos.SchoolsMap
		data, _ := json.Marshal(lastEntrySchoolsMap)
		_ = json.Unmarshal(data, &schoolsMap)

		capacities.Set(checkIn.LocationID, checkIn.Capacity)

		// Used to deep copy capacities
		var capacitiesCopy *orderedmap.OrderedMap[string, int64]
		data, _ = json.Marshal(capacities)
		_ = json.Unmarshal(data, &capacitiesCopy)

		checkInEntry := &dtos.CheckInsLocationEntryRaw{
			Capacities: capacitiesCopy,
			Schools:    schoolsMap,
		}

		schoolValue, _ := checkInEntry.Schools.Get(schoolName)
		schoolValue++
		checkInEntry.Schools.Set(schoolName, schoolValue)

		checkInEntries.Set(
			checkIn.CreatedAt.Time.Format(time.RFC3339),
			checkInEntry,
		)
		lastEntrySchoolsMap = checkInEntry.Schools
	}

	return checkInEntries
}

// todo: refactor
func (service CheckInService) GetCheckInsEntriesRange(
	startDate time.Time,
	endDate time.Time,
	checkIns []*models.CheckIn,
	schools []*models.School,
) *orderedmap.OrderedMap[string, *dtos.CheckInsLocationEntryRaw] {
	schoolsIDNameMap, _ := service.schools.GetSchoolMaps(schools)

	checkInEntries := orderedmap.New[string, *dtos.CheckInsLocationEntryRaw]()
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dVal := timetools.StartOfDay(d)

		_, schoolsMap := service.schools.GetSchoolMaps(schools)

		checkInEntry := &dtos.CheckInsLocationEntryRaw{
			Capacities: orderedmap.New[string, int64](),
			Schools:    schoolsMap,
		}

		checkInEntries.Set(dVal.Format(time.RFC3339), checkInEntry)
	}

	for i := range checkIns {
		datetime := timetools.StartOfDay(checkIns[i].CreatedAt.Time)
		schoolName := schoolsIDNameMap[checkIns[i].SchoolID]

		checkInEntry, _ := checkInEntries.Get(datetime.Format(time.RFC3339))

		schoolValue, _ := checkInEntry.Schools.Get(schoolName)
		schoolValue++
		checkInEntry.Schools.Set(schoolName, schoolValue)

		capacity, present := checkInEntry.Capacities.Get(checkIns[i].LocationID)
		if !present {
			capacity = 0
		}

		if checkIns[i].Capacity > capacity {
			capacity = checkIns[i].Capacity
		}

		checkInEntry.Capacities.Set(checkIns[i].LocationID, capacity)
	}

	return checkInEntries
}

func (service CheckInService) GetAllOfDay(
	ctx context.Context,
	locationID string,
	date time.Time,
) ([]*models.CheckIn, error) {
	return service.checkins.GetAllInRange(
		ctx,
		[]string{locationID},
		timetools.StartOfDay(date),
		timetools.EndOfDay(date),
	)
}

func (service CheckInService) GetAllInRange(
	ctx context.Context,
	locationIDs []string,
	startDate time.Time,
	endDate time.Time,
) ([]*models.CheckIn, error) {
	return service.checkins.GetAllInRange(
		ctx,
		locationIDs,
		startDate,
		endDate,
	)
}

func (service CheckInService) GetByID(
	ctx context.Context,
	location *models.Location,
	id int64,
) (*models.CheckIn, error) {
	return service.checkins.GetByID(ctx, location, id)
}

func (service CheckInWriterService) Delete(ctx context.Context, id int64) error {
	return service.checkins.Delete(ctx, id)
}
