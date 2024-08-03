package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	wstools "github.com/xdoubleu/essentia/pkg/communication/ws"
	contexttools "github.com/xdoubleu/essentia/pkg/context"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	timetools "github.com/xdoubleu/essentia/pkg/time"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type LocationService struct {
	locations repositories.LocationRepository
	checkins  repositories.CheckInRepository
	schools   SchoolService
	users     UserService
	// WebSocketService is an internal service
	websocket *WebSocketService
}

// todo: refactor
func (service LocationService) GetCheckInsEntriesDay(ctx context.Context, locationIDs []string, date time.Time) (*orderedmap.OrderedMap[string, *dtos.CheckInsLocationEntryRaw], error) {
	checkIns, err := service.GetAllCheckInsOfDay(
		ctx,
		locationIDs,
		date,
	)
	if err != nil {
		return nil, err
	}

	schools, err := service.schools.GetAll(ctx)
	if err != nil {
		return nil, err
	}

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

	return checkInEntries, nil
}

// todo: refactor
func (service LocationService) GetCheckInsEntriesRange(
	ctx context.Context,
	locationIDs []string,
	startDate time.Time,
	endDate time.Time,
) (*orderedmap.OrderedMap[string, *dtos.CheckInsLocationEntryRaw], error) {
	startDate = timetools.StartOfDay(startDate)
	endDate = timetools.EndOfDay(endDate)

	checkIns, err := service.GetAllCheckInsInRange(
		ctx,
		locationIDs,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, err
	}

	schools, err := service.schools.GetAll(ctx)
	if err != nil {
		return nil, err
	}

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

	return checkInEntries, nil
}

func (service LocationService) GetAllCheckInsOfDay(
	ctx context.Context,
	locationIDs []string,
	date time.Time,
) ([]*models.CheckIn, error) {
	return service.GetAllCheckInsInRange(
		ctx,
		locationIDs,
		timetools.StartOfDay(date),
		timetools.EndOfDay(date),
	)
}

func (service LocationService) GetAllCheckInsInRange(
	ctx context.Context,
	locationIDs []string,
	startDate time.Time,
	endDate time.Time,
) ([]*models.CheckIn, error) {
	user := contexttools.GetValue[models.User](ctx, constants.UserContextKey)

	locations, err := service.GetByIDs(ctx, locationIDs)
	if err != nil {
		return nil, err
	}
	for _, location := range locations {
		if user.Role == models.DefaultRole && location.UserID != user.ID {
			return nil, errortools.ErrResourceNotFound
		}
	}

	return service.checkins.GetAllInRange(
		ctx,
		locationIDs,
		startDate,
		endDate,
	)
}

func (service LocationService) GetCheckInByID(
	ctx context.Context,
	location *models.Location,
	id int64,
) (*models.CheckIn, error) {
	return service.checkins.GetByID(ctx, location, id)
}

func (service LocationService) DeleteCheckIn(ctx context.Context, id int64) error {
	return service.checkins.Delete(ctx, id)
}

func (service *LocationService) InitializeWS() error {
	locations, err := service.GetAll(context.Background())
	if err != nil {
		return err
	}

	service.websocket.allTopic, err = service.websocket.handler.AddTopic(
		"*",
		func(_ *wstools.Topic) (any, error) { return service.GetAllStates(context.Background()) },
	)
	if err != nil {
		return err
	}

	for _, location := range locations {
		err = service.websocket.AddLocation(location)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service LocationService) GetWSHandler() http.HandlerFunc {
	return service.websocket.Handler()
}

func (service LocationService) NewCheckIn(location models.Location) {
	location.Available--
	service.websocket.NewLocationState(location)
}

func (service LocationService) GetTotalCount(ctx context.Context) (*int64, error) {
	return service.locations.GetTotalCount(ctx)
}

func (service LocationService) GetAll(ctx context.Context) ([]*models.Location, error) {
	locations, err := service.locations.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	locationIDs := []string{}
	for _, location := range locations {
		locationIDs = append(locationIDs, location.ID)
	}

	checkInsToday, err := service.GetAllCheckInsOfDay(ctx, locationIDs, time.Now())
	if err != nil {
		return nil, err
	}

	checkInsYesterday, err := service.GetAllCheckInsOfDay(ctx, locationIDs, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	for i := range locations {
		err = locations[i].SetFields(checkInsToday, checkInsYesterday)
		if err != nil {
			return nil, err
		}
	}

	return locations, nil
}

func (service LocationService) GetAllStates(ctx context.Context) ([]dtos.LocationStateDto, error) {
	locations, err := service.GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	var result []dtos.LocationStateDto
	for _, location := range locations {
		result = append(result, dtos.NewLocationStateDto(*location))
	}

	return result, nil
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

	locationIDs := []string{}
	for _, location := range locations {
		locationIDs = append(locationIDs, location.ID)
	}

	checkInsToday, err := service.GetAllCheckInsOfDay(ctx, locationIDs, time.Now())
	if err != nil {
		return nil, err
	}

	checkInsYesterday, err := service.GetAllCheckInsOfDay(ctx, locationIDs, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	for i := range locations {
		err = locations[i].SetFields(checkInsToday, checkInsYesterday)
		if err != nil {
			return nil, err
		}
	}

	return locations, nil
}

func (service LocationService) GetByIDs(ctx context.Context, ids []string) ([]*models.Location, error) {
	locations, err := service.locations.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	checkInsToday, err := service.GetAllCheckInsOfDay(ctx, ids, time.Now())
	if err != nil {
		return nil, err
	}

	checkInsYesterday, err := service.GetAllCheckInsOfDay(ctx, ids, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	for i := range locations {
		err = locations[i].SetFields(checkInsToday, checkInsYesterday)
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
	location, err := service.locations.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user := contexttools.GetValue[models.User](ctx, constants.UserContextKey)
	if user.Role == models.DefaultRole && location.UserID != user.ID {
		return nil, errortools.ErrResourceNotFound
	}

	checkInsToday, err := service.GetAllCheckInsOfDay(ctx, []string{location.ID}, time.Now())
	if err != nil {
		return nil, err
	}

	checkInsYesterday, err := service.GetAllCheckInsOfDay(ctx, []string{location.ID}, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	err = location.SetFields(checkInsToday, checkInsYesterday)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (service LocationService) GetByUserID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	location, err := service.locations.GetByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	checkInsToday, err := service.GetAllCheckInsOfDay(ctx, []string{location.ID}, time.Now())
	if err != nil {
		return nil, err
	}

	checkInsYesterday, err := service.GetAllCheckInsOfDay(ctx, []string{location.ID}, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	err = location.SetFields(checkInsToday, checkInsYesterday)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (service LocationService) GetDefaultUserByUserID(ctx context.Context, id string) (*models.User, error) {
	user, err := service.users.GetByID(ctx, id, models.DefaultRole)
	if err != nil {
		return nil, err
	}

	var location *models.Location
	location, err = service.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.Location = location

	return user, nil
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
	user, err := service.users.Create(ctx, &dtos.CreateUserDto{
		Username: username,
		Password: password,
	}, models.DefaultRole)
	if err != nil {
		return nil, err
	}

	location, err := service.locations.Create(
		ctx,
		name,
		capacity,
		timeZone,
		user.ID,
	)
	if err != nil {
		service.users.Delete(ctx, user.ID, models.DefaultRole)
		return nil, err
	}

	checkInsToday, err := service.GetAllCheckInsOfDay(ctx, []string{location.ID}, time.Now())
	if err != nil {
		return nil, err
	}

	checkInsYesterday, err := service.GetAllCheckInsOfDay(ctx, []string{location.ID}, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	err = location.SetFields(checkInsToday, checkInsYesterday)
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

	checkInsToday, err := service.GetAllCheckInsOfDay(ctx, []string{location.ID}, time.Now())
	if err != nil {
		return err
	}

	checkInsYesterday, err := service.GetAllCheckInsOfDay(ctx, []string{location.ID}, time.Now().Add(-24*time.Hour))
	if err != nil {
		return err
	}

	err = location.SetFields(checkInsToday, checkInsYesterday)
	if err != nil {
		return err
	}

	if updateLocationDto.Name != nil {
		err = service.websocket.UpdateLocation(location)
		if err != nil {
			return err
		}
	}

	service.websocket.NewLocationState(*location)

	return nil
}

func (service LocationService) Delete(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	location, err := service.GetByID(ctx, id)
	if err != nil {
		return nil, errortools.ErrResourceNotFound
	}

	user, err := service.users.Delete(ctx, location.UserID, models.DefaultRole)
	if err != nil {
		return nil, err
	}

	err = service.locations.Delete(ctx, location)
	if err != nil {
		service.users.Recreate(ctx, user)
		return nil, err
	}

	return location, service.websocket.DeleteLocation(location)
}
