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

type StateService struct {
	logger    *slog.Logger
	state     repositories.StateRepository
	websocket *WebSocketService
	Current   models.State
}

func NewStateService(logger *slog.Logger, ctx context.Context, repo repositories.StateRepository, websocket *WebSocketService) StateService {
	service := StateService{
		logger:    logger,
		state:     repo,
		websocket: websocket,
	}

	state, err := service.get(ctx, true)
	if err != nil {
		panic(err)
	}

	service.Current = *state

	return service
}

func (service *StateService) InitializeWS(ctx context.Context) error {
	err := service.websocket.SetStateTopic(func(ctx context.Context) (*models.State, error) { return service.get(ctx, false) })
	if err != nil {
		return err
	}

	go service.startPolling(service.logger, ctx)
	return nil
}

func (service *StateService) get(ctx context.Context, fetchPersistentState bool) (*models.State, error) {
	state := &service.Current
	var err error

	if fetchPersistentState {
		state, err = service.state.Get(ctx)
		if err != nil {
			return nil, err
		}
	}

	state.IsDatabaseActive = service.state.IsDatabaseActive(ctx)

	return state, nil
}

func (service StateService) startPolling(logger *slog.Logger, ctx context.Context) {
	sentry.GoRoutineErrorHandler("State Polling", ctx, func(ctx context.Context) error {
		for ctx.Err() != context.Canceled {
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
		return nil
	})
}

func (service *StateService) UpdateState(ctx context.Context, stateDto *dtos.StateDto) (*models.State, error) {
	err := service.state.UpdateKey(ctx, models.IsMaintenanceKey, strconv.FormatBool(stateDto.IsMaintenance))
	if err != nil {
		return nil, err
	}

	service.Current = models.State{
		IsMaintenance:    stateDto.IsMaintenance,
		IsDatabaseActive: service.state.IsDatabaseActive(ctx),
	}

	service.websocket.NewAppState(service.Current)

	return &service.Current, nil
}
