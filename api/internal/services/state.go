package services

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/XDoubleU/essentia/pkg/logging"
	"github.com/XDoubleU/essentia/pkg/sentry"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type CurrentState struct {
	value models.State
	mu    *sync.RWMutex
}

type StateService struct {
	logger    *slog.Logger
	state     repositories.StateRepository
	websocket *WebSocketService
	Current   *CurrentState
}

func NewStateService(
	ctx context.Context,
	logger *slog.Logger,
	repo repositories.StateRepository,
	websocket *WebSocketService,
) StateService {
	//nolint:exhaustruct //Current is set later
	service := StateService{
		logger:    logger,
		state:     repo,
		websocket: websocket,
	}

	state, err := service.get(ctx, true)
	if err != nil {
		panic(err)
	}

	service.Current = &CurrentState{
		value: *state,
		mu:    &sync.RWMutex{},
	}

	return service
}

func (service *StateService) InitializeWS(ctx context.Context) error {
	err := service.websocket.SetStateTopic(
		func(ctx context.Context) (*models.State, error) { return service.get(ctx, false) },
	)
	if err != nil {
		return err
	}

	go service.startPolling(ctx, service.logger)
	return nil
}

func (service *StateService) get(
	ctx context.Context,
	fetchPersistentState bool,
) (*models.State, error) {
	var state models.State
	var err error

	if fetchPersistentState {
		var newState *models.State
		newState, err = service.state.Get(ctx)
		if err != nil {
			return nil, err
		}
		state = *newState
	} else {
		state = service.Current.Get()
	}

	state.IsDatabaseActive = service.state.IsDatabaseActive(ctx)

	return &state, nil
}

func (service *StateService) startPolling(ctx context.Context, logger *slog.Logger) {
	sentry.GoRoutineWrapper(
		ctx,
		logger,
		"State Polling",
		func(ctx context.Context, logger *slog.Logger) error {
			for ctx.Err() != context.Canceled {
				newState, err := service.get(ctx, false)
				if err != nil {
					logger.Error(
						"something went wrong while fetching current state",
						logging.ErrAttr(err),
					)
					continue
				}

				_, changed := service.Current.update(*newState)
				if changed {
					service.websocket.NewAppState(*newState)
				}

				time.Sleep(10 * time.Second) //nolint:mnd //no magic number
			}
			return nil
		},
	)
}

func (service *StateService) UpdateState(
	ctx context.Context,
	stateDto dtos.StateDto,
) (*models.State, error) {
	err := service.state.UpdateKey(
		ctx,
		models.IsMaintenanceKey,
		strconv.FormatBool(stateDto.IsMaintenance),
	)
	if err != nil {
		return nil, err
	}

	newState, changed := service.Current.update(models.State{
		IsMaintenance:    stateDto.IsMaintenance,
		IsDatabaseActive: service.state.IsDatabaseActive(ctx),
	})

	if changed {
		service.websocket.NewAppState(newState)
	}

	return &newState, nil
}

func (state *CurrentState) Get() models.State {
	state.mu.RLock()
	defer state.mu.RUnlock()

	return state.value
}

func (state *CurrentState) update(newState models.State) (models.State, bool) {
	state.mu.Lock()
	defer state.mu.Unlock()

	changed := false

	if state.value != newState {
		state.value = newState
		changed = true
	}

	return state.value, changed
}
