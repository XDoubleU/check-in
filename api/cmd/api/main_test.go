package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"check-in/api/internal/config"
	"check-in/api/internal/database"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/services"
	"check-in/api/internal/tests"

	"github.com/XDoubleU/essentia/pkg/http_tools"
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

var mainTestEnv *tests.MainTestEnv //nolint:gochecknoglobals //global var for tests
var tokens Tokens                  //nolint:gochecknoglobals //global var for tests
var cfg config.Config              //nolint:gochecknoglobals //global var for tests
var fixtureData FixtureData        //nolint:gochecknoglobals //global var for tests

func clearAll(services services.Services) {
	user, err := services.Users.GetByUsername(context.Background(), "Admin")
	if user != nil {
		err = services.Users.Delete(context.Background(), user.ID, user.Role)
	}

	if err != http_tools.ErrRecordNotFound && err != nil {
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

	fixtureData.AmountOfSchools = 0
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

	tokens.DefaultAccessToken, err = services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		fixtureData.DefaultUser.ID,
		cfg.AccessExpiry,
		false,
	)
	if err != nil {
		panic(err)
	}

	tokens.DefaultRefreshToken, err = services.Auth.CreateCookie(context.Background(),
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

func fixtures(db database.DB) {
	services := services.New(db)

	clearAll(services)
	userFixtures(services)
	locationFixtures(services)
	schoolFixtures(services)
}

func removeFixtures(db database.DB) {
	services := services.New(db)

	clearAll(services)
}

func TestMain(m *testing.M) {
	var err error

	cfg = config.New()
	cfg.Env = config.TestEnv
	cfg.Throttle = false

	mainTestEnv, err = tests.SetupGlobal(
		cfg.DB.Dsn,
		cfg.DB.MaxConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		panic(err)
	}

	fixtures(mainTestEnv.TestDB)
	exitCode := m.Run()
	removeFixtures(mainTestEnv.TestDB)

	err = tests.TeardownGlobal(mainTestEnv)
	if err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

func setupTest(
	t *testing.T,
	mainTestEnv *tests.MainTestEnv,
) (tests.TestEnv, *application) {
	t.Parallel()
	testEnv := tests.SetupSingle(mainTestEnv)

	testApp := &application{
		config:     cfg,
		services:   services.New(testEnv.TestTx),
		hideErrors: false,
	}

	return testEnv, testApp
}
