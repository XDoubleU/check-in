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
	return Services{
		Auth:      AuthService{db: db},
		CheckIns:  CheckInService{db: db},
		Locations: LocationService{db: db},
		Schools:   SchoolService{db: db},
		Users:     UserService{db: db},
		WebSockets: WebSocketService{
			subscribers: make(map[*websocket.Conn]models.Subscriber),
		},
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
