package services

import (
	"context"
	"log/slog"

	"check-in/api/internal/config"
	"check-in/api/internal/repositories"
	"check-in/api/internal/shared"
)

type Services struct {
	Auth           AuthService
	CheckInsWriter CheckInWriterService
	Locations      LocationService
	Schools        SchoolService
	Users          UserService
	State          StateService
	WebSocket      *WebSocketService
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	config config.Config,
	repositories repositories.Repositories,
	utcNowTimeProvider shared.UTCNowTimeProvider,
) Services {
	websocket := NewWebSocketService(logger, []string{config.WebURL})
	state := NewStateService(ctx, logger, repositories.State, websocket)

	users := UserService{
		users: repositories.Users,
	}
	schools := SchoolService{
		schools:         repositories.Schools,
		schoolIDNameMap: make(map[int64]string),
	}
	locations := LocationService{
		locations:     repositories.Locations,
		checkins:      repositories.CheckIns,
		schools:       schools,
		users:         users,
		websocket:     websocket,
		getTimeNowUTC: utcNowTimeProvider,
	}
	auth := AuthService{
		auth:          repositories.Auth,
		users:         users,
		locations:     locations,
		getTimeNowUTC: utcNowTimeProvider,
	}
	checkInsWriter := CheckInWriterService{
		checkins:  repositories.CheckInsWriter,
		locations: locations,
		schools:   schools,
	}

	err := locations.InitializeWS(ctx)
	if err != nil {
		panic(err)
	}

	err = state.InitializeWS(ctx)
	if err != nil {
		panic(err)
	}

	return Services{
		Auth:           auth,
		CheckInsWriter: checkInsWriter,
		Locations:      locations,
		Schools:        schools,
		Users:          users,
		State:          state,
		WebSocket:      websocket,
	}
}
