package services

import (
	"check-in/api/internal/config"
	"check-in/api/internal/repositories"
)

type Services struct {
	Auth           AuthService
	CheckInsWriter CheckInWriterService
	CheckIns       CheckInService
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
	checkIns := CheckInService{
		checkins: repositories.CheckIns,
		schools:  schools,
	}
	locations := LocationService{
		locations: repositories.Locations,
		checkins:  checkIns,
		users:     users,
		websocket: websocket,
	}
	checkInsWriter := CheckInWriterService{
		checkins:  repositories.CheckIns,
		schools:   schools,
		locations: locations,
	}

	err := locations.InitializeWS()
	if err != nil {
		panic(err)
	}

	return Services{
		Auth:           auth,
		CheckInsWriter: checkInsWriter,
		CheckIns:       checkIns,
		Locations:      locations,
		Schools:        schools,
		Users:          users,
	}
}
