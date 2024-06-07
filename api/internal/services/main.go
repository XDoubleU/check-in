package services

import (
	"errors"

	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"nhooyr.io/websocket"

	"check-in/api/internal/models"
)

type Services struct {
	Auth       AuthService
	CheckIns   CheckInService
	Locations  LocationService
	Schools    SchoolService
	Users      UserService
	WebSockets WebSocketService
}

func New(db postgres.DB) Services {
	checkIns := CheckInService{db: db}
	schools := SchoolService{db: db}
	locations := LocationService{db: db, schools: schools, checkins: checkIns}
	auth := AuthService{db: db, locations: locations}
	users := UserService{db: db, locations: locations}
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

func handleError(err error) error {
	var pgxError *pgconn.PgError
	errors.As(err, &pgxError)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return http_tools.ErrRecordNotFound
	case pgxError.Code == "23503":
		return http_tools.ErrRecordNotFound
	case pgxError.Code == "23505":
		return http_tools.ErrRecordUniqueValue
	default:
		return err
	}
}
