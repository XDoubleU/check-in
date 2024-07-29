package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	configtools "github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
	"github.com/xdoubleu/essentia/pkg/logging"

	"check-in/api/internal/config"
	"check-in/api/internal/models"
	"check-in/api/internal/services"
)

type TestEnv struct {
	tx       postgres.PgxSyncTx
	cfg      config.Config
	services services.Services
	Fixtures Fixtures
}

type Tokens struct {
	AdminAccessToken    *http.Cookie
	ManagerAccessToken  *http.Cookie
	DefaultAccessToken  *http.Cookie
	DefaultRefreshToken *http.Cookie
}

type Fixtures struct {
	Tokens          Tokens
	AdminUser       *models.User
	ManagerUser     *models.User
	DefaultUser     *models.User
	DefaultLocation *models.Location
}

var db postgres.DB

func (env *TestEnv) defaultFixtures() {
	var err error

	password := "testpassword"
	env.Fixtures.AdminUser, err = env.services.Users.Create(context.Background(),
		"Admin",
		password,
		models.AdminRole,
	)
	if err != nil {
		panic(err)
	}

	env.Fixtures.ManagerUser, err = env.services.Users.Create(context.Background(),
		"Manager",
		password,
		models.ManagerRole,
	)
	if err != nil {
		panic(err)
	}

	env.Fixtures.Tokens.AdminAccessToken, err = env.services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		env.Fixtures.AdminUser.ID,
		env.cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	env.Fixtures.Tokens.ManagerAccessToken, err = env.services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		env.Fixtures.ManagerUser.ID,
		env.cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	timezone, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		panic(err)
	}

	env.Fixtures.DefaultLocation, err = env.services.Locations.Create(
		context.Background(),
		"TestLocation",
		20,
		timezone.String(),
		"Default",
		"testpassword",
	)
	if err != nil {
		panic(err)
	}

	env.Fixtures.DefaultUser, err = env.services.Users.GetByID(
		context.Background(),
		env.Fixtures.DefaultLocation.UserID,
		models.DefaultRole,
	)
	if err != nil {
		panic(err)
	}

	env.Fixtures.Tokens.DefaultAccessToken, err = env.services.Auth.CreateCookie(
		context.Background(),
		models.AccessScope,
		env.Fixtures.DefaultUser.ID,
		env.cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	env.Fixtures.Tokens.DefaultRefreshToken, err = env.services.Auth.CreateCookie(
		context.Background(),
		models.RefreshScope,
		env.Fixtures.DefaultUser.ID,
		env.cfg.RefreshExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}
}

func (env *TestEnv) createManagerUsers(amount int) []*models.User {
	var err error
	password := "testpassword"

	users := []*models.User{}
	for i := 0; i < amount; i++ {
		var newUser *models.User
		newUser, err = env.services.Users.Create(context.Background(),
			fmt.Sprintf("TestManagerUser%d", i),
			password,
			models.ManagerRole,
		)
		if err != nil {
			panic(err)
		}

		users = append(users, newUser)
	}

	return users
}

func (env *TestEnv) createLocations(amount int) []*models.Location {
	var err error

	timezone, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		panic(err)
	}

	locations := []*models.Location{}
	for i := 0; i < amount; i++ {
		var location *models.Location
		location, err = env.services.Locations.Create(
			context.Background(),
			fmt.Sprintf("TestLocation%d", i),
			20,
			timezone.String(),
			fmt.Sprintf("TestDefaultUser%d", i),
			"testpassword",
		)
		if err != nil {
			panic(err)
		}

		locations = append(locations, location)
	}

	return locations
}

func (env *TestEnv) createCheckIns(location *models.Location, schoolID int64, amount int) []*models.CheckIn {
	var err error

	checkIns := []*models.CheckIn{}
	for i := 0; i < amount; i++ {
		var checkIn *models.CheckIn
		checkIn, err = env.services.CheckIns.Create(
			context.Background(),
			location,
			//nolint:exhaustruct // other fields are optional
			&models.School{ID: schoolID},
		)
		if err != nil {
			panic(err)
		}

		checkIns = append(checkIns, checkIn)
	}

	return checkIns
}

func (env *TestEnv) createSchools(amount int) []*models.School {
	schools := []*models.School{}
	for i := 0; i < amount; i++ {
		school, err := env.services.Schools.Create(context.Background(),
			fmt.Sprintf("TestSchool%d", i))
		if err != nil {
			panic(err)
		}
		schools = append(schools, school)
	}

	return schools
}

func TestMain(m *testing.M) {
	var err error

	// only to acquire db dsn
	cfg := config.New()
	db, err = postgres.Connect(
		logging.NewNopLogger(),
		cfg.DBDsn,
		25,
		"15m",
		5,
		15*time.Second,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func setup(t *testing.T) (TestEnv, *Application) {
	t.Parallel()

	cfg := config.New()
	cfg.Env = configtools.TestEnv
	cfg.Throttle = false

	tx := postgres.CreatePgxSyncTx(context.Background(), db)

	testApp := NewApp(logging.NewNopLogger(), cfg, tx)

	testEnv := TestEnv{
		tx:       tx,
		cfg:      cfg,
		services: testApp.services,
		Fixtures: Fixtures{},
	}

	testEnv.defaultFixtures()

	return testEnv, testApp
}

func (env *TestEnv) teardown() {
	err := env.tx.Rollback(context.Background())
	if err != nil {
		panic(err)
	}
}
