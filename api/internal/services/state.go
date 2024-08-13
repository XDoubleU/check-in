package services

import (
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/sentry"
)

type IsDatabaseActiveFunc = func(ctx context.Context) bool

type StateService struct {
	logger           *slog.Logger
	state            repositories.StateRepository
	websocket        *WebSocketService
	isDatabaseActive IsDatabaseActiveFunc
	Current          models.State
}

func NewStateService(logger *slog.Logger, repo repositories.StateRepository, websocket *WebSocketService, isDatabaseActive IsDatabaseActiveFunc) StateService {
	service := StateService{
		logger:           logger,
		state:            repo,
		websocket:        websocket,
		isDatabaseActive: isDatabaseActive,
	}

	state, err := service.get(context.Background(), true)
	if err != nil {
		panic(err)
	}

	service.Current = *state

	return service
}

func (service *StateService) InitializeWS() error {
	err := service.websocket.SetStateTopic()
	if err != nil {
		return err
	}

	go service.startPolling(service.logger)
	return nil
}

func (service StateService) get(ctx context.Context, fetchPersistentState bool) (*models.State, error) {
	state := &service.Current
	var err error

	if fetchPersistentState {
		state, err = service.state.Get(ctx)
		if err != nil {
			return nil, err
		}
	}

	state.IsDatabaseActive = service.isDatabaseActive(ctx)

	return state, nil
}

func (service StateService) startPolling(logger *slog.Logger) {
	sentry.GoRoutineErrorHandler("State Polling", context.Background(), func(ctx context.Context) error {
		for {
			newState, err := service.get(ctx, false)
			if err != nil {
				logger.Error("something went wrong while fetching current state", logging.ErrAttr(err))
				continue
			}

			if service.Current != *newState {
				service.Current = *newState
				service.websocket.NewAppState(*newState)
			}

			time.Sleep(10 * time.Second)
		}
	})
}

func (service StateService) UpdateState(ctx context.Context, stateDto *dtos.StateDto) (*models.State, error) {
	err := service.state.UpdateKey(ctx, models.IsMaintenanceKey, strconv.FormatBool(stateDto.IsMaintenance))
	if err != nil {
		return nil, err
	}

	service.Current = models.State{
		IsMaintenance:    stateDto.IsMaintenance,
		IsDatabaseActive: service.isDatabaseActive(ctx),
	}

	service.websocket.NewAppState(service.Current)

	return &service.Current, nil
}
