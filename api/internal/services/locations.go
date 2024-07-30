package services

import (
	"context"
	"encoding/json"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/tools"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type LocationService struct {
	locations repositories.LocationRepository
	schools   SchoolService
	checkins  CheckInService
	websocket *WebSocketService
}

// todo: refactor
func (service LocationService) GetCheckInsEntriesDay(
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
func (service LocationService) GetCheckInsEntriesRange(
	startDate time.Time,
	endDate time.Time,
	checkIns []*models.CheckIn,
	schools []*models.School,
) *orderedmap.OrderedMap[string, *dtos.CheckInsLocationEntryRaw] {
	schoolsIDNameMap, _ := service.schools.GetSchoolMaps(schools)

	checkInEntries := orderedmap.New[string, *dtos.CheckInsLocationEntryRaw]()
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dVal := tools.StartOfDay(d)

		_, schoolsMap := service.schools.GetSchoolMaps(schools)

		checkInEntry := &dtos.CheckInsLocationEntryRaw{
			Capacities: orderedmap.New[string, int64](),
			Schools:    schoolsMap,
		}

		checkInEntries.Set(dVal.Format(time.RFC3339), checkInEntry)
	}

	for i := range checkIns {
		datetime := tools.StartOfDay(checkIns[i].CreatedAt.Time)
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

func (service LocationService) GetTotalCount(ctx context.Context) (*int64, error) {
	return service.locations.GetTotalCount(ctx)
}

func (service LocationService) GetAll(ctx context.Context) ([]*models.Location, error) {
	locations, err := service.locations.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for i := range locations {
		err = service.prepareLocation(ctx, locations[i])
		if err != nil {
			return nil, err
		}
	}

	return locations, nil
}

func (service LocationService) GetAllPaginated(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]*models.Location, error) {
	locations, err := service.locations.GetAllPaginated(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range locations {
		err = service.prepareLocation(ctx, locations[i])
		if err != nil {
			return nil, err
		}
	}

	return locations, nil
}

func (service LocationService) GetByID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	location, err := service.locations.GetBy(ctx, "WHERE locations.id = $1", id)
	if err != nil {
		return nil, err
	}

	err = service.prepareLocation(ctx, location)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (service LocationService) GetByUserID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	location, err := service.locations.GetBy(ctx, "WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}

	err = service.prepareLocation(ctx, location)
	if err != nil {
		return nil, err
	}

	return location, nil
}

// todo: refactor
func (service LocationService) prepareLocation(
	ctx context.Context,
	location *models.Location,
) error {
	var checkInsToday []*models.CheckIn
	var checkInsYesterday []*models.CheckIn
	var err error

	loc, _ := time.LoadLocation(location.TimeZone)
	today := time.Now().In(loc)
	yesterday := today.AddDate(0, 0, -1)

	checkInsToday, err = service.checkins.GetAllOfDay(ctx, location.ID, today)
	if err != nil {
		return err
	}

	checkInsYesterday, err = service.checkins.GetAllOfDay(ctx, location.ID, yesterday)
	if err != nil {
		return err
	}

	location.SetCheckInRelatedFields(
		checkInsToday,
		checkInsYesterday,
	)

	err = location.NormalizeName()
	if err != nil {
		return err
	}

	return nil
}

func (service LocationService) GetByName(
	ctx context.Context,
	name string,
) (*models.Location, error) {
	locations, err := service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, location := range locations {
		var output bool

		output, err = location.CompareNormalizedName(name)
		if err != nil {
			return nil, err
		}

		if output {
			return location, nil
		}
	}

	return nil, errortools.ErrResourceNotFound
}

func (service LocationService) Create(
	ctx context.Context,
	name string,
	capacity int64,
	timeZone string,
	username string,
	password string,
) (*models.Location, error) {
	location, err := service.locations.Create(
		ctx,
		name,
		capacity,
		timeZone,
		username,
		password,
	)
	if err != nil {
		return nil, err
	}

	err = service.prepareLocation(ctx, location)
	if err != nil {
		return nil, err
	}

	err = service.websocket.AddLocation(location)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (service LocationService) Update(
	ctx context.Context,
	location *models.Location,
	user *models.User,
	updateLocationDto dtos.UpdateLocationDto,
) error {
	err := service.locations.Update(ctx, location, user, updateLocationDto)
	if err != nil {
		return err
	}

	err = service.prepareLocation(ctx, location)
	if err != nil {
		return err
	}

	// todo this check sucks, all subs lose their subscription
	if updateLocationDto.Name != nil {
		err = service.websocket.DeleteLocation(location)
		if err != nil {
			return err
		}
		err = service.websocket.AddLocation(location)
		if err != nil {
			return err
		}
	}

	service.websocket.NewLocationState(*location)

	return nil
}

func (service LocationService) Delete(
	ctx context.Context,
	location *models.Location,
) error {
	err := service.locations.Delete(ctx, location)
	if err != nil {
		return err
	}

	return service.websocket.DeleteLocation(location)
}
