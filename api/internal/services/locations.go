package services

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/XDoubleU/essentia/pkg/database"
	errortools "github.com/XDoubleU/essentia/pkg/errors"
	"github.com/XDoubleU/essentia/pkg/grapher"
	timetools "github.com/XDoubleU/essentia/pkg/time"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
	"check-in/api/internal/shared"
)

type LocationService struct {
	locations     repositories.LocationRepository
	checkins      repositories.CheckInRepository
	schools       SchoolService
	users         UserService
	websocket     *WebSocketService
	getTimeNowUTC shared.UTCNowTimeProvider
}

func (service *LocationService) InitializeWS(ctx context.Context) error {
	locations, err := service.GetAll(ctx, nil, true)
	if err != nil {
		return err
	}

	err = service.websocket.SetAllLocationsTopic(service.GetAllStates)
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

func (service LocationService) GetCheckInsEntriesDay(
	ctx context.Context,
	user *models.User,
	locationIDs []string,
	date time.Time,
) ([]string, map[string][]int, map[string][]int, error) {
	_, checkIns, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		locationIDs,
		date,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	g := grapher.New[int](grapher.Cumulative, grapher.None, time.RFC3339, time.Second)
	capacitiesGrapher := grapher.New[int](grapher.Normal, grapher.None, time.RFC3339, time.Second)

	for _, checkIn := range checkIns {
		datetime := timetools.LocationIndependentTime(checkIn.CreatedAt.Time, "UTC")
		g.AddPoint(datetime, 1, checkIn.SchoolName)
		capacitiesGrapher.AddPoint(datetime, int(checkIn.Capacity), checkIn.LocationID)
	}

	dateStrings, valueMap := g.ToSlices()
	_, capacitiesMap := capacitiesGrapher.ToSlices()

	return dateStrings, capacitiesMap, valueMap, nil
}

func (service LocationService) GetCheckInsEntriesRange(
	ctx context.Context,
	user *models.User,
	locationIDs []string,
	startDate time.Time,
	endDate time.Time,
) ([]string, map[string][]int, map[string][]int, error) {
	startDate = timetools.StartOfDay(startDate)
	endDate = timetools.EndOfDay(endDate)

	_, checkIns, err := service.GetAllCheckInsInRange(
		ctx,
		user,
		false,
		locationIDs,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	schoolIDNameMap, err := service.schools.SchoolIDNameMap(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	g := grapher.New[int](grapher.CumulativeSameDate, grapher.None, time.RFC3339, time.Second)
	capacitiesGrapher := grapher.New[int](grapher.Normal, grapher.None, time.RFC3339, time.Second)

	for i := startDate; i.Before(endDate); i = i.AddDate(0, 0, 1) {
		for _, schoolName := range schoolIDNameMap {
			g.AddPoint(i, 0, schoolName)
		}

		for _, locationID := range locationIDs {
			capacitiesGrapher.AddPoint(i, 0, locationID)
		}
	}

	for i := range checkIns {
		datetime := timetools.StartOfDay(checkIns[i].CreatedAt.Time)

		g.AddPoint(datetime, 1, checkIns[i].SchoolName)
		capacitiesGrapher.AddPoint(
			datetime,
			int(checkIns[i].Capacity),
			checkIns[i].LocationID,
		)
	}

	dateStrings, valueMap := g.ToSlices()
	_, capacitiesMap := capacitiesGrapher.ToSlices()

	return dateStrings, capacitiesMap, valueMap, nil
}

func (service LocationService) GetAllCheckInsOfDay(
	ctx context.Context,
	user *models.User,
	allowAnonymous bool,
	locationIDs []string,
	date time.Time,
) ([]*models.CheckIn, []*dtos.CheckInDto, error) {
	return service.GetAllCheckInsInRange(
		ctx,
		user,
		allowAnonymous,
		locationIDs,
		timetools.StartOfDay(date),
		timetools.EndOfDay(date),
	)
}

func (service LocationService) GetAllCheckInsInRange(
	ctx context.Context,
	user *models.User,
	allowAnonymous bool,
	locationIDs []string,
	startDate time.Time,
	endDate time.Time,
) ([]*models.CheckIn, []*dtos.CheckInDto, error) {
	if len(locationIDs) == 0 {
		return make([]*models.CheckIn, 0), make([]*dtos.CheckInDto, 0), nil
	}

	locations, err := service.getByIDs(ctx, locationIDs)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, nil, errortools.NewNotFoundError(
				"locations",
				locationIDs,
				"ids",
			)
		}
		return nil, nil, err
	}

	if !allowAnonymous && user.Role == models.DefaultRole {
		for _, location := range locations {
			if location.UserID != user.ID {
				return nil, nil, errortools.NewNotFoundError(
					"location",
					location.ID,
					"id",
				)
			}
		}
	}

	checkIns := []*models.CheckIn{}

	for _, location := range locations {
		var locationCheckIns []*models.CheckIn
		locationCheckIns, err = service.checkins.GetAllInRange(
			ctx,
			location.ID,
			startDate,
			endDate,
		)
		if err != nil {
			return nil, nil, err
		}

		for _, checkIn := range locationCheckIns {
			checkIn.CreatedAt.Time = timetools.LocationIndependentTime(
				checkIn.CreatedAt.Time,
				"UTC",
			)
			checkIns = append(checkIns, checkIn)
		}
	}

	slices.SortFunc(checkIns, func(i, j *models.CheckIn) int {
		return i.CreatedAt.Time.Compare(j.CreatedAt.Time)
	})

	schoolIDNameMap, err := service.schools.SchoolIDNameMap(ctx)
	if err != nil {
		return nil, nil, err
	}

	checkInDtos := make([]*dtos.CheckInDto, 0)
	for _, checkIn := range checkIns {
		checkInDto := &dtos.CheckInDto{
			ID:         checkIn.ID,
			LocationID: checkIn.LocationID,
			SchoolName: schoolIDNameMap[checkIn.SchoolID],
			Capacity:   checkIn.Capacity,
			CreatedAt:  checkIn.CreatedAt,
		}
		checkInDtos = append(checkInDtos, checkInDto)
	}

	return checkIns, checkInDtos, nil
}

func (service LocationService) GetCheckInByID(
	ctx context.Context,
	location *models.Location,
	id int64,
) (*models.CheckIn, error) {
	return service.checkins.GetByID(ctx, location, id)
}

func (service LocationService) DeleteCheckIn(
	ctx context.Context,
	user *models.User,
	locationID string,
	checkInID int64,
) (*dtos.CheckInDto, error) {
	location, err := service.GetByID(ctx, user, locationID)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, errortools.NewNotFoundError("location", locationID, "id")
		}
		return nil, err
	}

	checkIn, err := service.GetCheckInByID(ctx, location, checkInID)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, errortools.NewNotFoundError("checkIn", checkInID, "id")
		}
		return nil, err
	}

	today := service.getTimeNowUTC()
	startOfToday := timetools.StartOfDay(today)
	endOfToday := timetools.EndOfDay(today)

	if !(checkIn.CreatedAt.Time.After(startOfToday) &&
		checkIn.CreatedAt.Time.Before(endOfToday)) {
		return nil, errortools.NewBadRequestError(
			errors.New("checkIn didn't occur today and thus can't be deleted"),
		)
	}

	err = service.checkins.Delete(ctx, checkInID)
	if err != nil {
		return nil, err
	}

	schoolIDNameMap, err := service.schools.SchoolIDNameMap(ctx)
	if err != nil {
		return nil, err
	}

	checkInDto := &dtos.CheckInDto{
		ID:         checkIn.ID,
		LocationID: checkIn.LocationID,
		SchoolName: schoolIDNameMap[checkIn.SchoolID],
		Capacity:   checkIn.Capacity,
		CreatedAt:  checkIn.CreatedAt,
	}

	return checkInDto, nil
}

func (service LocationService) NewCheckIn(location models.Location) {
	location.Available--
	service.websocket.NewLocationState(location)
}

func (service LocationService) GetTotalCount(ctx context.Context) (*int64, error) {
	return service.locations.GetTotalCount(ctx)
}

func (service LocationService) GetAll(
	ctx context.Context,
	user *models.User,
	allowAnonymous bool,
) ([]*models.Location, error) {
	locations, err := service.locations.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	locationIDs := []string{}
	for _, location := range locations {
		locationIDs = append(locationIDs, location.ID)
	}

	checkInsToday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		allowAnonymous,
		locationIDs,
		service.getTimeNowUTC(),
	)
	if err != nil {
		return nil, err
	}

	checkInsYesterday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		allowAnonymous,
		locationIDs,
		service.getTimeNowUTC().Add(-24*time.Hour),
	)
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

func (service LocationService) GetAllStates(
	ctx context.Context,
) ([]dtos.LocationStateDto, error) {
	locations, err := service.GetAll(ctx, nil, true)
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
	user *models.User,
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

	checkInsToday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		locationIDs,
		service.getTimeNowUTC(),
	)
	if err != nil {
		return nil, err
	}

	checkInsYesterday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		locationIDs,
		service.getTimeNowUTC().Add(-24*time.Hour),
	)
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

func (service LocationService) getByIDs(
	ctx context.Context,
	ids []string,
) ([]*models.Location, error) {
	return service.locations.GetByIDs(ctx, ids)
}

func (service LocationService) GetByID(
	ctx context.Context,
	user *models.User,
	id string,
) (*models.Location, error) {
	location, err := service.locations.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, errortools.NewNotFoundError("location", id, "id")
		}
		return nil, err
	}

	if user.Role == models.DefaultRole && location.UserID != user.ID {
		return nil, errortools.NewNotFoundError("location", location.ID, "id")
	}

	checkInsToday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		[]string{location.ID},
		service.getTimeNowUTC(),
	)
	if err != nil {
		return nil, err
	}

	checkInsYesterday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		[]string{location.ID},
		service.getTimeNowUTC().Add(-24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	err = location.SetFields(checkInsToday, checkInsYesterday)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (service LocationService) GetByUser(
	ctx context.Context,
	user *models.User,
) (*models.Location, error) {
	location, err := service.locations.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	checkInsToday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		[]string{location.ID},
		service.getTimeNowUTC(),
	)
	if err != nil {
		return nil, err
	}

	checkInsYesterday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		[]string{location.ID},
		service.getTimeNowUTC().Add(-24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	err = location.SetFields(checkInsToday, checkInsYesterday)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (service LocationService) GetDefaultUserByUserID(
	ctx context.Context,
	id string,
) (*models.User, error) {
	user, err := service.users.GetByID(ctx, id, models.DefaultRole)
	if err != nil {
		return nil, err
	}

	var location *models.Location
	location, err = service.GetByUser(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Location = location

	return user, nil
}

func (service LocationService) GetByName(
	ctx context.Context,
	user *models.User,
	name string,
) (*models.Location, error) {
	locations, err := service.GetAll(ctx, user, false)
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

	return nil, errortools.NewNotFoundError("location", name, "name")
}

func (service LocationService) Create(
	ctx context.Context,
	user *models.User,
	createLocationDto dtos.CreateLocationDto,
) (*models.Location, error) {
	err := service.checkForConflictsOnCreate(ctx, user, createLocationDto)
	if err != nil {
		return nil, err
	}

	defaultUser, err := service.users.Create(ctx, dtos.CreateUserDto{
		Username: createLocationDto.Username,
		Password: createLocationDto.Password,
	}, models.DefaultRole)
	if err != nil {
		return nil, err
	}

	location, err := service.locations.Create(
		ctx,
		createLocationDto.Name,
		createLocationDto.Capacity,
		createLocationDto.TimeZone,
		defaultUser.ID,
	)
	if err != nil {
		_, err2 := service.users.Delete(ctx, defaultUser.ID, models.DefaultRole)
		if err2 != nil {
			return nil, err2
		}

		return nil, err
	}

	checkInsToday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		[]string{location.ID},
		service.getTimeNowUTC(),
	)
	if err != nil {
		return nil, err
	}

	checkInsYesterday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		[]string{location.ID},
		service.getTimeNowUTC().Add(-24*time.Hour),
	)
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

func (service LocationService) Recreate(
	ctx context.Context,
	location *models.Location,
) (*models.Location, error) {
	return service.locations.Create(
		ctx,
		location.Name,
		location.Capacity,
		location.TimeZone,
		location.UserID,
	)
}

func (service LocationService) checkForConflictsOnCreate(
	ctx context.Context,
	user *models.User,
	createLocationDto dtos.CreateLocationDto,
) error {
	existingLocation, _ := service.GetByName(
		ctx,
		user,
		createLocationDto.Name,
	)
	if existingLocation != nil {
		return errortools.NewConflictError("location", createLocationDto.Name, "name")
	}

	existingUser, _ := service.users.GetByUsername(
		ctx,
		createLocationDto.Username,
	)

	if existingUser != nil {
		return errortools.NewConflictError(
			"user",
			createLocationDto.Username,
			"username",
		)
	}

	return nil
}

func (service LocationService) Update(
	ctx context.Context,
	user *models.User,
	id string,
	updateLocationDto dtos.UpdateLocationDto,
) (*models.Location, error) {
	err := service.checkForConflictsOnUpdate(ctx, user, updateLocationDto)
	if err != nil {
		return nil, err
	}

	oldLocation, err := service.GetByID(ctx, user, id)
	if err != nil {
		return nil, err
	}

	location, err := service.locations.Update(ctx, *oldLocation, updateLocationDto)
	if err != nil {
		return nil, err
	}

	_, err = service.users.Update(ctx, location.UserID, dtos.UpdateUserDto{
		Username: updateLocationDto.Username,
		Password: updateLocationDto.Password,
	}, models.DefaultRole)
	if err != nil {
		//nolint:exhaustruct //other fields are optional
		_, err2 := service.locations.Update(ctx, *location, dtos.UpdateLocationDto{
			Name:     &oldLocation.Name,
			Capacity: &oldLocation.Capacity,
			TimeZone: &oldLocation.TimeZone,
		})
		if err2 != nil {
			return nil, err2
		}

		return nil, err
	}

	checkInsToday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		[]string{location.ID},
		service.getTimeNowUTC(),
	)
	if err != nil {
		return nil, err
	}

	checkInsYesterday, _, err := service.GetAllCheckInsOfDay(
		ctx,
		user,
		false,
		[]string{location.ID},
		service.getTimeNowUTC().Add(-24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	err = location.SetFields(checkInsToday, checkInsYesterday)
	if err != nil {
		return nil, err
	}

	if updateLocationDto.Name != nil {
		err = service.websocket.UpdateLocation(location)
		if err != nil {
			return nil, err
		}
	}

	service.websocket.NewLocationState(*location)

	return location, nil
}

func (service LocationService) checkForConflictsOnUpdate(
	ctx context.Context,
	user *models.User,
	updateLocationDto dtos.UpdateLocationDto,
) error {
	if updateLocationDto.Name != nil {
		existingLocation, _ := service.GetByName(
			ctx,
			user,
			*updateLocationDto.Name,
		)

		if existingLocation != nil {
			return errortools.NewConflictError(
				"location",
				*updateLocationDto.Name,
				"name",
			)
		}
	}

	if updateLocationDto.Username != nil {
		existingUser, _ := service.users.GetByUsername(
			ctx,
			*updateLocationDto.Username,
		)

		if existingUser != nil {
			return errortools.NewConflictError(
				"user",
				*updateLocationDto.Username,
				"username",
			)
		}
	}

	return nil
}

func (service LocationService) Delete(
	ctx context.Context,
	user *models.User,
	id string,
) (*models.Location, error) {
	location, err := service.GetByID(ctx, user, id)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, errortools.NewNotFoundError("location", id, "id")
		}
		return nil, err
	}

	err = service.locations.Delete(ctx, location)
	if err != nil {
		return nil, err
	}

	_, err = service.users.Delete(ctx, location.UserID, models.DefaultRole)
	if err != nil {
		_, err2 := service.Recreate(ctx, location)
		if err2 != nil {
			return nil, err2
		}

		return nil, err
	}

	return location, service.websocket.DeleteLocation(location)
}
