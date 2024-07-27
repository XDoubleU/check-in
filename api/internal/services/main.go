package services

import (
	"check-in/api/internal/config"
	"check-in/api/internal/repositories"
)

type Services struct {
	Auth      AuthService
	CheckIns  CheckInService
	Locations LocationService
	Schools   SchoolService
	Users     UserService
	WebSocket WebSocketService
}

func New(config config.Config, repositories repositories.Repositories) Services {
	checkIns := CheckInService{
		checkins: repositories.CheckIns,
	}
	schools := SchoolService{
		schools: repositories.Schools,
	}
	locations := LocationService{
		locations: repositories.Locations,
		schools:   schools,
		checkins:  checkIns,
	}
	users := UserService{
		users:     repositories.Users,
		locations: locations,
	}
	auth := AuthService{
		auth:  repositories.Auth,
		users: users,
	}

	websocket, err := NewWebSocketService(config.WebURL, locations)
	if err != nil {
		panic(err)
	}

	return Services{
		Auth:      auth,
		CheckIns:  checkIns,
		Locations: locations,
		Schools:   schools,
		Users:     users,
		WebSocket: *websocket,
	}
}
