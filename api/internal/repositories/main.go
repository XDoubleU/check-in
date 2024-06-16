package repositories

import (
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"nhooyr.io/websocket"

	"check-in/api/internal/models"
)

type Repositories struct {
	Auth       AuthRepository
	CheckIns   CheckInRepository
	Locations  LocationRepository
	Schools    SchoolRepository
	Users      UserRepository
	WebSockets WebSocketRepository
}

func New(db postgres.DB) Repositories {
	checkIns := CheckInRepository{db: db}
	schools := SchoolRepository{db: db}
	locations := LocationRepository{db: db, schools: schools, checkins: checkIns}
	auth := AuthRepository{db: db, locations: locations}
	users := UserRepository{db: db, locations: locations}
	websockets := WebSocketRepository{
		subscribers: make(map[*websocket.Conn]models.Subscriber),
	}

	return Repositories{
		Auth:       auth,
		CheckIns:   checkIns,
		Locations:  locations,
		Schools:    schools,
		Users:      users,
		WebSockets: websockets,
	}
}
