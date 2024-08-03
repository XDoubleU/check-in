package services

import (
	"check-in/api/internal/config"
	"check-in/api/internal/repositories"
)

type Services struct {
	Auth           AuthService
	CheckInsWriter CheckInWriterService
	Locations      LocationService
	Schools        SchoolService
	Users          UserService
}

func New(config config.Config, repositories repositories.Repositories) Services {
	websocket := NewWebSocketService(config.WebURL)

	users := UserService{
		users: repositories.Users,
	}
	auth := AuthService{
		auth:  repositories.Auth,
		users: users,
	}
	schools := SchoolService{
		schools: repositories.Schools,
	}
	locations := LocationService{
		locations: repositories.Locations,
		checkins:  repositories.CheckIns,
		schools:   schools,
		users:     users,
		websocket: websocket,
	}
	checkInsWriter := CheckInWriterService{
		checkins:  repositories.CheckInsWriter,
		locations: locations,
		schools:   schools,
	}

	err := locations.InitializeWS()
	if err != nil {
		panic(err)
	}

	return Services{
		Auth:           auth,
		CheckInsWriter: checkInsWriter,
		Locations:      locations,
		Schools:        schools,
		Users:          users,
	}
}
