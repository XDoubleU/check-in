package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	configtools "github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/database"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/logging"

	"check-in/api/internal/config"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
	"check-in/api/internal/services"
)

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

// todo get rid of these
var mainTestEnv *database.MainTestEnv[*pgxpool.Pool, postgres.PgxSyncTx] //nolint:gochecknoglobals //global var for tests
var tokens Tokens                                                        //nolint:gochecknoglobals //global var for tests
var cfg config.Config                                                    //nolint:gochecknoglobals //global var for tests
var fixtureData FixtureData                                              //nolint:gochecknoglobals //global var for tests

func clearAll(services services.Services) {
	user, err := services.Users.GetByUsername(context.Background(), "Admin")
	if user != nil {
		err = services.Users.Delete(context.Background(), user.ID, user.Role)
	}

	if err != nil && !errors.Is(err, errortools.ErrResourceNotFound) {
		panic(err)
	}

	users, err := services.Users.GetAll(context.Background())
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		err = services.Users.Delete(context.Background(), user.ID, user.Role)
		if err != nil {
			panic(err)
		}
	}

	fixtureData.AmountOfManagerUsers = 0

	locations, err := services.Locations.GetAll(context.Background())
	if err != nil {
		panic(err)
	}

	for _, location := range locations {
		err = services.Locations.Delete(context.Background(), location)
		if err != nil {
			panic(err)
		}
	}

	fixtureData.AmountOfLocations = 0

	schools, err := services.Schools.GetAll(context.Background())
	if err != nil {
		panic(err)
	}

	for _, school := range schools {
		if school.ID == 1 {
			continue
		}

		err = services.Schools.Delete(context.Background(), school.ID)
		if err != nil {
			panic(err)
		}
	}

	fixtureData.AmountOfSchools = 1
}

func userFixtures(services services.Services) {
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
	fixtureData.AmountOfManagerUsers++

	fixtureData.AdminUser = adminUser
	fixtureData.ManagerUser = managerUser

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

		fixtureData.AmountOfManagerUsers++

		fixtureData.ManagerUsers = append(fixtureData.ManagerUsers, newUser)
	}

	adminAccessToken, err := services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		adminUser.ID,
		cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	managerAccessToken, err := services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		managerUser.ID,
		cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	tokens.AdminAccessToken = adminAccessToken
	tokens.ManagerAccessToken = managerAccessToken
}

func locationFixtures(services services.Services) {
	timezone, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		panic(err)
	}

	fixtureData.DefaultLocation, err = services.Locations.Create(
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
	fixtureData.AmountOfLocations++

	fixtureData.DefaultUser, err = services.Users.GetByID(
		context.Background(),
		fixtureData.DefaultLocation.UserID,
		models.DefaultRole,
	)
	if err != nil {
		panic(err)
	}

	tokens.DefaultAccessToken, err = services.Auth.CreateCookie(
		context.Background(),
		models.AccessScope,
		fixtureData.DefaultUser.ID,
		cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	tokens.DefaultRefreshToken, err = services.Auth.CreateCookie(
		context.Background(),
		models.RefreshScope,
		fixtureData.DefaultUser.ID,
		cfg.RefreshExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		newCap := fixtureData.DefaultLocation.Capacity + 1
		err = services.Locations.Update(
			context.Background(),
			fixtureData.DefaultLocation,
			fixtureData.AdminUser,
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
			fixtureData.DefaultLocation,
			//nolint:exhaustruct // other fields are optional
			&models.School{ID: 1},
		)
		if err != nil {
			panic(err)
		}

		fixtureData.CheckIns = append(fixtureData.CheckIns, checkIn)
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
		fixtureData.AmountOfLocations++

		var user *models.User
		user, err = services.Users.GetByID(
			context.Background(),
			location.UserID,
			models.DefaultRole,
		)
		if err != nil {
			panic(err)
		}
		fixtureData.DefaultUsers = append(fixtureData.DefaultUsers, user)

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

		fixtureData.Locations = append(fixtureData.Locations, location)
	}
}

func schoolFixtures(services services.Services) {
	for i := 0; i < 20; i++ {
		school, err := services.Schools.Create(context.Background(),
			fmt.Sprintf("TestSchool%d", i))
		if err != nil {
			panic(err)
		}
		fixtureData.Schools = append(fixtureData.Schools, school)
		fixtureData.AmountOfSchools++
	}
}

func fixtures(db postgres.DB) {
	services := services.New(cfg, repositories.New(db))

	clearAll(services)
	userFixtures(services)
	locationFixtures(services)
	schoolFixtures(services)
}

func removeFixtures(db postgres.DB) {
	services := services.New(cfg, repositories.New(db))

	clearAll(services)
}

func TestMain(m *testing.M) {
	var err error

	cfg = config.New()
	cfg.Env = configtools.TestEnv
	cfg.Throttle = false

	db, err := postgres.Connect(
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

	mainTestEnv = database.CreateMainTestEnv(db, postgres.CreatePgxSyncTx)

	fixtures(mainTestEnv.TestDB)
	exitCode := m.Run()
	removeFixtures(mainTestEnv.TestDB)

	os.Exit(exitCode)
}

func setupTest(
	t *testing.T,
	mainTestEnv *database.MainTestEnv[*pgxpool.Pool, postgres.PgxSyncTx],
) (database.TestEnv[postgres.PgxSyncTx], *Application) {
	t.Parallel()
	testEnv := mainTestEnv.SetupSingle()

	testApp := NewApp(logging.NewNopLogger(), cfg, testEnv.Tx)

	return testEnv, testApp
}
