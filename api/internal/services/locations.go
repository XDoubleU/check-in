package services

import (
	"context"
	"net/http"
	"time"

	wstools "github.com/xdoubleu/essentia/pkg/communication/ws"
	errortools "github.com/xdoubleu/essentia/pkg/errors"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type LocationService struct {
	locations repositories.LocationRepository
	checkins  CheckInService
	users     UserService
	// WebSocketService is an internal service
	websocket *WebSocketService
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

	for i := range locations {
		err = service.prepareLocation(ctx, locations[i])
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
	location *models.Location,
) error {
	err := service.locations.Delete(ctx, location)
	if err != nil {
		return err
	}

	return service.websocket.DeleteLocation(location)
}
