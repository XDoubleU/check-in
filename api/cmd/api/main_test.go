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
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

type TestEnv struct {
	ctx context.Context
	tx  *postgres.PgxSyncTx
	app *Application
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

var mainTx *postgres.PgxSyncTx //nolint:gochecknoglobals //needed for tests
var cfg config.Config          //nolint:gochecknoglobals //needed for tests
var fixtures Fixtures          //nolint:gochecknoglobals //needed for tests
var mainTestApp *Application   //nolint:gochecknoglobals //needed for tests
var testCtx context.Context    //nolint:gochecknoglobals //needed for tests

func defaultFixtures(ctx context.Context, app *Application) {
	var err error

	password := "testpassword"
	fixtures.AdminUser, err = app.services.Users.Create(ctx,
		&dtos.CreateUserDto{
			Username: "Admin",
			Password: password,
		},
		models.AdminRole,
	)
	if err != nil {
		panic(err)
	}

	ctx = app.contextSetUser(ctx, *fixtures.AdminUser)

	fixtures.ManagerUser, err = app.services.Users.Create(ctx,
		&dtos.CreateUserDto{
			Username: "Manager",
			Password: password,
		},
		models.ManagerRole,
	)
	if err != nil {
		panic(err)
	}

	fixtures.Tokens.AdminAccessToken, err = app.services.Auth.CreateCookie(
		ctx,
		models.AccessScope,
		fixtures.AdminUser.ID,
		app.config.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	fixtures.Tokens.ManagerAccessToken, err = app.services.Auth.CreateCookie(
		ctx,
		models.AccessScope,
		fixtures.ManagerUser.ID,
		app.config.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	timezone, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		panic(err)
	}

	fixtures.DefaultLocation, err = app.services.Locations.Create(
		ctx,
		fixtures.AdminUser,
		&dtos.CreateLocationDto{
			Name:     "TestLocation",
			Capacity: 20,
			TimeZone: timezone.String(),
			Username: "Default",
			Password: "testpassword",
		},
	)
	if err != nil {
		panic(err)
	}

	fixtures.DefaultUser, err = app.services.Locations.GetDefaultUserByUserID(
		ctx,
		fixtures.DefaultLocation.UserID,
	)
	if err != nil {
		panic(err)
	}

	fixtures.Tokens.DefaultAccessToken, err = app.services.Auth.CreateCookie(
		ctx,
		models.AccessScope,
		fixtures.DefaultUser.ID,
		app.config.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	fixtures.Tokens.DefaultRefreshToken, err = app.services.Auth.CreateCookie(
		ctx,
		models.RefreshScope,
		fixtures.DefaultUser.ID,
		app.config.RefreshExpiry,
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
		newUser, err = env.app.services.Users.Create(env.ctx,
			&dtos.CreateUserDto{
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

	timezone, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		panic(err)
	}

	locations := []*models.Location{}
	for i := 0; i < amount; i++ {
		var location *models.Location
		location, err = env.app.services.Locations.Create(
			env.ctx,
			fixtures.AdminUser,
			&dtos.CreateLocationDto{
				Name:     fmt.Sprintf("TestLocation%d", i),
				Capacity: 20,
				TimeZone: timezone.String(),
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

	defaultUser, err := env.app.services.Locations.GetDefaultUserByUserID(env.ctx, location.UserID)
	if err != nil {
		panic(err)
	}

	checkIns := []*dtos.CheckInDto{}
	for i := 0; i < amount; i++ {
		var checkIn *dtos.CheckInDto
		checkIn, err = env.app.services.CheckInsWriter.Create(
			env.ctx,
			&dtos.CreateCheckInDto{
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
		school, err := env.app.services.Schools.Create(env.ctx,
			&dtos.SchoolDto{
				Name: fmt.Sprintf("TestSchool%d", i),
			})
		if err != nil {
			panic(err)
		}
		schools = append(schools, school)
	}

	return schools
}

func TestMain(m *testing.M) {
	var err error

	cfg = config.New()
	cfg.Env = configtools.TestEnv
	cfg.Throttle = false

	postgresDB, err := postgres.Connect(
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

	mainTx = postgres.CreatePgxSyncTx(context.Background(), postgresDB)

	testCtx = context.Background()
	mainTestApp = NewApp(logging.NewNopLogger(), cfg, mainTx)

	defaultFixtures(testCtx, mainTestApp)

	code := m.Run()
	err = mainTx.Rollback(context.Background())
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func setup(_ *testing.T) (*TestEnv, *Application) {
	tx := postgres.CreatePgxSyncTx(context.Background(), mainTx)

	testApp := *mainTestApp
	testApp.SetDB(tx)

	testEnv := &TestEnv{
		ctx: testCtx,
		tx:  tx,
		app: &testApp,
	}

	return testEnv, &testApp
}

func (env *TestEnv) teardown() {
	err := env.tx.Rollback(context.Background())
	if err != nil {
		panic(err)
	}

	env.app.ctxCancel()
}
