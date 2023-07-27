package services

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"nhooyr.io/websocket"

	"check-in/api/internal/database"
	"check-in/api/internal/models"
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrRecordUniqueValue = errors.New("record unique value already used")
)

type Services struct {
	Auth       AuthService
	CheckIns   CheckInService
	Locations  LocationService
	Schools    SchoolService
	Users      UserService
	WebSockets WebSocketService
}

func New(db database.DB) Services {
	locations := LocationService{db: db}
	auth := AuthService{db: db, locations: locations}
	checkIns := CheckInService{db: db}
	schools := SchoolService{db: db}
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
		return ErrRecordNotFound
	case pgxError.Code == "23503":
		return ErrRecordNotFound
	case pgxError.Code == "23505":
		return ErrRecordUniqueValue
	default:
		return err
	}
}
