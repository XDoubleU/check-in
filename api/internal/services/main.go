package services

import (
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"

	"nhooyr.io/websocket"
)

type Services struct {
	Auth       AuthService
	CheckIns   CheckInService
	Locations  LocationService
	Schools    SchoolService
	Users      UserService
	WebSockets WebSocketService
}

func New(repositories repositories.Repositories) Services {
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

	websockets := WebSocketService{
		subscribers: make(map[*websocket.Conn]models.Subscriber),
	}

	return Services{
		Auth:       auth,
		CheckIns:   checkIns,
		Locations:  locations,
		Schools:    schools,
		Users:      users,
		WebSockets: websockets,
	}
}
