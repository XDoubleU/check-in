package repositories

import (
	"github.com/XDoubleU/essentia/pkg/database/postgres"
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

func New(db postgres.DB) Repositories {
	checkInsWriter := CheckInWriteRepository{db: db}
	checkIns := CheckInRepository{db: db}
	schools := SchoolRepository{db: db}
	locations := LocationRepository{db: db}
	auth := AuthRepository{db: db}
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
