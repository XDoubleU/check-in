package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	configtools "github.com/XDoubleU/essentia/pkg/config"
	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"

	"check-in/api/internal/config"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/shared"
)

type TestEnv struct {
	ctx      context.Context
	app      Application
	fixtures Fixtures
}

type Tokens struct {
	AdminAccessToken    *http.Cookie
	ManagerAccessToken  *http.Cookie
	ManagerRefreshToken *http.Cookie
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

var cfg config.Config        //nolint:gochecknoglobals //required
var postgresDB *pgxpool.Pool //nolint:gochecknoglobals //required

var timesToCheck = []shared.LocalNowTimeProvider{ //nolint:gochecknoglobals //required
	time.Now,
	func() time.Time { return getTimeNow(23, false, "Europe/Brussels") },
	func() time.Time { return getTimeNow(00, true, "Europe/Brussels") },
	func() time.Time { return getTimeNow(01, true, "Europe/Brussels") },
	func() time.Time { return getTimeNow(23, false, "UTC") },
	func() time.Time { return getTimeNow(00, false, "UTC") },
	func() time.Time { return getTimeNow(01, false, "UTC") },
}

func (env *TestEnv) defaultFixtures() {
	var err error

	_, err = env.app.services.State.UpdateState(context.Background(), dtos.StateDto{
		IsMaintenance: false,
	})
	if err != nil {
		panic(err)
	}

	password := "testpassword"
	env.fixtures.AdminUser, err = env.app.services.Users.Create(env.ctx,

		dtos.CreateUserDto{
			Username: "Admin",
			Password: password,
		},
		models.AdminRole,
	)
	if err != nil {
		panic(err)
	}

	env.ctx = env.app.contextSetUser(env.ctx, *env.fixtures.AdminUser)

	env.fixtures.ManagerUser, err = env.app.services.Users.Create(env.ctx,

		dtos.CreateUserDto{
			Username: "Manager",
			Password: password,
		},
		models.ManagerRole,
	)
	if err != nil {
		panic(err)
	}

	env.fixtures.Tokens.AdminAccessToken = env.createAccessToken(
		*env.fixtures.AdminUser,
	)

	env.fixtures.Tokens.ManagerAccessToken = env.createAccessToken(
		*env.fixtures.ManagerUser,
	)
	env.fixtures.Tokens.ManagerRefreshToken = env.createRefreshToken(
		*env.fixtures.ManagerUser,
	)

	env.fixtures.DefaultLocation, err = env.app.services.Locations.Create(
		env.ctx,
		env.fixtures.AdminUser,

		dtos.CreateLocationDto{
			Name:     "TestLocation",
			Capacity: 20,
			TimeZone: "Europe/Brussels",
			Username: "Default",
			Password: "testpassword",
		},
	)
	if err != nil {
		panic(err)
	}

	env.fixtures.DefaultUser, err = env.app.services.Locations.GetDefaultUserByUserID(
		env.ctx,
		env.fixtures.DefaultLocation.UserID,
	)
	if err != nil {
		panic(err)
	}

	env.fixtures.Tokens.DefaultAccessToken = env.createAccessToken(
		*env.fixtures.DefaultUser,
	)
	env.fixtures.Tokens.DefaultRefreshToken = env.createRefreshToken(
		*env.fixtures.DefaultUser,
	)
}

func (env *TestEnv) createAccessToken(user models.User) *http.Cookie {
	var err error

	accessToken, err := env.app.services.Auth.CreateCookie(
		env.ctx,
		models.AccessScope,
		user.ID,
		env.app.config.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	return accessToken
}

func (env *TestEnv) createRefreshToken(user models.User) *http.Cookie {
	var err error

	refreshToken, err := env.app.services.Auth.CreateCookie(
		env.ctx,
		models.RefreshScope,
		user.ID,
		env.app.config.RefreshExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	return refreshToken
}

func (env *TestEnv) clearAllData() {
	var err error

	//nolint:exhaustruct //other fields are optional
	fakeAdminUser := &models.User{
		Role: models.AdminRole,
	}

	locations, _ := env.app.services.Locations.GetAll(env.ctx, nil, true)
	for _, location := range locations {
		_, err = env.app.services.Locations.Delete(
			env.ctx,
			fakeAdminUser,
			location.ID,
		)
		if err != nil {
			panic(err)
		}
	}

	users, _ := env.app.services.Users.GetAll(env.ctx)
	for _, user := range users {
		_, err = env.app.services.Users.Delete(env.ctx, user.ID, user.Role)
		if err != nil {
			panic(err)
		}
	}

	adminUser, _ := env.app.services.Users.GetByUsername(env.ctx, "Admin")
	if adminUser != nil {
		_, err = env.app.services.Users.Delete(
			env.ctx,
			adminUser.ID,
			adminUser.Role,
		)
		if err != nil {
			panic(err)
		}
	}

	schools, _ := env.app.services.Schools.GetAll(env.ctx)
	for _, school := range schools {
		if school.ID == 1 {
			continue
		}

		_, err = env.app.services.Schools.Delete(env.ctx, school.ID)
		if err != nil {
			panic(err)
		}
	}
}

func (env *TestEnv) createManagerUsers(amount int) []*models.User {
	var err error
	password := "testpassword"

	users := []*models.User{}
	for i := 0; i < amount; i++ {
		var newUser *models.User
		newUser, err = env.app.services.Users.Create(env.ctx,

			dtos.CreateUserDto{
				Username: fmt.Sprintf("TestManagerUser%d", i),
				Password: password,
			},
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

	locations := []*models.Location{}
	for i := 0; i < amount; i++ {
		var location *models.Location
		location, err = env.app.services.Locations.Create(
			env.ctx,
			env.fixtures.AdminUser,

			dtos.CreateLocationDto{
				Name:     fmt.Sprintf("TestLocation%d", i),
				Capacity: 20,
				TimeZone: "Europe/Brussels",
				Username: fmt.Sprintf("TestDefaultUser%d", i),
				Password: "testpassword",
			},
		)
		if err != nil {
			panic(err)
		}

		locations = append(locations, location)
	}

	return locations
}

func (env *TestEnv) createCheckIns(
	location *models.Location,
	schoolID int64,
	amount int,
) []*dtos.CheckInDto {
	var err error

	defaultUser, err := env.app.services.Locations.GetDefaultUserByUserID(
		env.ctx,
		location.UserID,
	)
	if err != nil {
		panic(err)
	}

	checkIns := []*dtos.CheckInDto{}
	for i := 0; i < amount; i++ {
		var checkIn *dtos.CheckInDto
		checkIn, err = env.app.services.CheckInsWriter.Create(
			env.ctx,

			dtos.CreateCheckInDto{
				SchoolID: schoolID,
			},
			defaultUser,
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
		name := fmt.Sprintf("TestSchool%d", i)

		school, err := env.app.services.Schools.GetByName(env.ctx, name)
		if err != nil && !errors.Is(err, database.ErrResourceNotFound) {
			panic(err)
		}

		if school == nil {
			school, err = env.app.services.Schools.Create(env.ctx,
				dtos.SchoolDto{
					Name: name,
				})
			if err != nil {
				panic(err)
			}
		}

		schools = append(schools, school)
	}

	return schools
}

func TestMain(m *testing.M) {
	var err error

	cfg = config.New(logging.NewNopLogger())
	cfg.Env = configtools.TestEnv
	cfg.Throttle = false

	postgresDB, err = postgres.Connect(
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

	ApplyMigrations(logging.NewNopLogger(), postgresDB)

	os.Exit(m.Run())
}

func runForAllTimes(
	t *testing.T,
	testFunc func(t *testing.T, testEnv TestEnv, testApp Application),
) {
	for _, timeNow := range timesToCheck {
		testEnv, testApp := setupSpecificTimeProvider(timeNow)
		testFunc(t, testEnv, testApp)
		testEnv.teardown()
	}
}

func getTimeNow(hour int, nextDay bool, tz string) time.Time {
	now := time.Now()

	if nextDay {
		now = now.Add(24 * time.Hour)
	}

	location, _ := time.LoadLocation(tz)

	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		hour,
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		location,
	)
}

func setupSpecificTimeProvider(
	timeProvider shared.LocalNowTimeProvider,
) (TestEnv, Application) {
	testApp := NewApp(slog.New(slog.NewTextHandler(os.Stdout, nil)), cfg, postgresDB, timeProvider)
	testEnv := TestEnv{
		ctx: context.Background(),
		app: *testApp,
		//nolint:exhaustruct //other fields are optional
		fixtures: Fixtures{},
	}

	testEnv.clearAllData()
	testEnv.defaultFixtures()

	return testEnv, *testApp
}

func setup(_ *testing.T) (TestEnv, Application) {
	return setupSpecificTimeProvider(time.Now)
}

func (env *TestEnv) teardown() {
	env.clearAllData()
}
