package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	configtools "github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
	"github.com/xdoubleu/essentia/pkg/logging"

	"check-in/api/internal/config"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
	"check-in/api/internal/services"
)

type TestEnv struct {
	tx       pgx.Tx
	cfg      config.Config
	Tokens   Tokens
	Fixtures FixtureData
}

type Tokens struct {
	AdminAccessToken    *http.Cookie
	ManagerAccessToken  *http.Cookie
	DefaultAccessToken  *http.Cookie
	DefaultRefreshToken *http.Cookie
}

type FixtureData struct {
	AdminUser            *models.User
	ManagerUser          *models.User
	DefaultUser          *models.User
	Schools              []*models.School
	ManagerUsers         []*models.User
	DefaultUsers         []*models.User
	Locations            []*models.Location
	CheckIns             []*models.CheckIn
	DefaultLocation      *models.Location
	AmountOfLocations    int
	AmountOfSchools      int
	AmountOfManagerUsers int
}

var db postgres.DB

func (env *TestEnv) userFixtures(services services.Services) {
	password := "testpassword"

	adminUser, err := services.Users.Create(context.Background(),
		"Admin",
		password,
		models.AdminRole,
	)
	if err != nil {
		panic(err)
	}

	managerUser, err := services.Users.Create(context.Background(),
		"Manager",
		password,
		models.ManagerRole,
	)
	if err != nil {
		panic(err)
	}
	env.Fixtures.AmountOfManagerUsers++

	env.Fixtures.AdminUser = adminUser
	env.Fixtures.ManagerUser = managerUser

	for i := 0; i < 10; i++ {
		var newUser *models.User
		newUser, err = services.Users.Create(context.Background(),
			fmt.Sprintf("TestManagerUser%d", i),
			password,
			models.ManagerRole,
		)
		if err != nil {
			panic(err)
		}

		env.Fixtures.AmountOfManagerUsers++

		env.Fixtures.ManagerUsers = append(env.Fixtures.ManagerUsers, newUser)
	}

	adminAccessToken, err := services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		adminUser.ID,
		env.cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	managerAccessToken, err := services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		managerUser.ID,
		env.cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	env.Tokens.AdminAccessToken = adminAccessToken
	env.Tokens.ManagerAccessToken = managerAccessToken
}

func (env *TestEnv) locationFixtures(services services.Services) {
	timezone, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		panic(err)
	}

	env.Fixtures.DefaultLocation, err = services.Locations.Create(
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
	env.Fixtures.AmountOfLocations++

	env.Fixtures.DefaultUser, err = services.Users.GetByID(
		context.Background(),
		env.Fixtures.DefaultLocation.UserID,
		models.DefaultRole,
	)
	if err != nil {
		panic(err)
	}

	env.Tokens.DefaultAccessToken, err = services.Auth.CreateCookie(
		context.Background(),
		models.AccessScope,
		env.Fixtures.DefaultUser.ID,
		env.cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	env.Tokens.DefaultRefreshToken, err = services.Auth.CreateCookie(
		context.Background(),
		models.RefreshScope,
		env.Fixtures.DefaultUser.ID,
		env.cfg.RefreshExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		newCap := env.Fixtures.DefaultLocation.Capacity + 1
		err = services.Locations.Update(
			context.Background(),
			env.Fixtures.DefaultLocation,
			env.Fixtures.AdminUser,
			//nolint:exhaustruct // other fields are optional
			dtos.UpdateLocationDto{
				Capacity: &newCap,
			},
		)
		if err != nil {
			panic(err)
		}

		var checkIn *models.CheckIn
		checkIn, err = services.CheckIns.Create(
			context.Background(),
			env.Fixtures.DefaultLocation,
			//nolint:exhaustruct // other fields are optional
			&models.School{ID: 1},
		)
		if err != nil {
			panic(err)
		}

		env.Fixtures.CheckIns = append(env.Fixtures.CheckIns, checkIn)
	}

	for i := 0; i < 20; i++ {
		var location *models.Location
		location, err = services.Locations.Create(
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
		env.Fixtures.AmountOfLocations++

		var user *models.User
		user, err = services.Users.GetByID(
			context.Background(),
			location.UserID,
			models.DefaultRole,
		)
		if err != nil {
			panic(err)
		}
		env.Fixtures.DefaultUsers = append(env.Fixtures.DefaultUsers, user)

		for j := 0; j < 5; j++ {
			_, err = services.CheckIns.Create(
				context.Background(),
				location,
				//nolint:exhaustruct // other fields are optional
				&models.School{ID: 1},
			)
			if err != nil {
				panic(err)
			}
		}

		env.Fixtures.Locations = append(env.Fixtures.Locations, location)
	}
}

func (env *TestEnv) schoolFixtures(services services.Services) {
	for i := 0; i < 20; i++ {
		school, err := services.Schools.Create(context.Background(),
			fmt.Sprintf("TestSchool%d", i))
		if err != nil {
			panic(err)
		}
		env.Fixtures.Schools = append(env.Fixtures.Schools, school)
		env.Fixtures.AmountOfSchools++
	}
}

func (env *TestEnv) fixtures() {
	services := services.New(env.cfg, repositories.New(env.tx))

	env.userFixtures(services)
	env.locationFixtures(services)
	env.schoolFixtures(services)
}

func TestMain(t *testing.T) {
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
}

func setup(t *testing.T) (TestEnv, *Application) {
	t.Parallel()

	cfg := config.New()
	cfg.Env = configtools.TestEnv
	cfg.Throttle = false

	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}

	testEnv := TestEnv{
		tx:  tx,
		cfg: cfg,
		Fixtures: FixtureData{
			Schools:      []*models.School{},
			ManagerUsers: []*models.User{},
			DefaultUsers: []*models.User{},
			Locations:    []*models.Location{},
			CheckIns:     []*models.CheckIn{},
		},
	}

	testEnv.fixtures()

	testApp := NewApp(logging.NewNopLogger(), cfg, testEnv.tx)

	return testEnv, testApp
}

func (env *TestEnv) teardown() {
	err := env.tx.Rollback(context.Background())
	if err != nil {
		panic(err)
	}
}
