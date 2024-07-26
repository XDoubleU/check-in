package repositories

import (
	"github.com/XDoubleU/essentia/pkg/database/postgres"
)

type Repositories struct {
	Auth      AuthRepository
	CheckIns  CheckInRepository
	Locations LocationRepository
	Schools   SchoolRepository
	Users     UserRepository
}

func New(db postgres.DB) Repositories {
	checkIns := CheckInRepository{db: db}
	schools := SchoolRepository{db: db}
	locations := LocationRepository{db: db}
	auth := AuthRepository{db: db}
	users := UserRepository{db: db}

	return Repositories{
		Auth:      auth,
		CheckIns:  checkIns,
		Locations: locations,
		Schools:   schools,
		Users:     users,
	}
}
