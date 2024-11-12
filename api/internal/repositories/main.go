package repositories

import (
	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"check-in/api/internal/shared"
)

type Repositories struct {
	Auth           AuthRepository
	CheckIns       CheckInRepository
	CheckInsWriter CheckInWriteRepository
	Locations      LocationRepository
	Schools        SchoolRepository
	Users          UserRepository
	State          StateRepository
}

func New(db postgres.DB, utcNowTimeProvider shared.UTCNowTimeProvider) Repositories {
	checkInsWriter := CheckInWriteRepository{db: db, getTimeNowUTC: utcNowTimeProvider}
	checkIns := CheckInRepository{db: db}
	schools := SchoolRepository{db: db}
	locations := LocationRepository{db: db}
	auth := AuthRepository{db: db, getTimeNowUTC: utcNowTimeProvider}
	users := UserRepository{db: db}
	state := StateRepository{db: db}

	return Repositories{
		Auth:           auth,
		CheckIns:       checkIns,
		CheckInsWriter: checkInsWriter,
		Locations:      locations,
		Schools:        schools,
		Users:          users,
		State:          state,
	}
}
