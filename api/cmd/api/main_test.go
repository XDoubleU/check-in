package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"check-in/api/internal/config"
	"check-in/api/internal/database"
	"check-in/api/internal/models"
	"check-in/api/internal/services"
	"check-in/api/internal/tests"
)

type Tokens struct {
	AdminAccessToken    http.Cookie
	ManagerAccessToken  http.Cookie
	DefaultAccessToken  http.Cookie
	DefaultRefreshToken http.Cookie
}

type FixtureData struct {
	AdminUser         models.User
	ManagerUser       models.User
	DefaultUser       models.User
	Schools           []models.School
	ManagerUsers      []models.User
	DefaultUsers      []models.User
	Locations         []models.Location
	DefaultLocation   *models.Location
	AmountOfLocations int
	AmountOfSchools   int
	AmountOfUsers     int
}

var mainTestEnv *tests.MainTestEnv //nolint:gochecknoglobals //global var for tests
var tokens Tokens                  //nolint:gochecknoglobals //global var for tests
var cfg config.Config              //nolint:gochecknoglobals //global var for tests
var logger *log.Logger             //nolint:gochecknoglobals //global var for tests
var fixtureData FixtureData        //nolint:gochecknoglobals //global var for tests

func userFixtures(services services.Services) (*Tokens, error) {
	fixtureData.AmountOfUsers = 0

	password := "testpassword"

	adminUser, err := services.Users.Create(context.Background(),
		"Admin",
		password,
		models.AdminRole,
	)
	if err != nil {
		return nil, err
	}

	managerUser, err := services.Users.Create(context.Background(),
		"Manager",
		password,
		models.ManagerRole,
	)
	if err != nil {
		return nil, err
	}

	defaultUser, err := services.Users.Create(context.Background(),
		"Default",
		password,
		models.DefaultRole,
	)
	if err != nil {
		return nil, err
	}

	fixtureData.AmountOfUsers += 3
	fixtureData.AdminUser = *adminUser
	fixtureData.ManagerUser = *managerUser
	fixtureData.DefaultUser = *defaultUser

	for i := 0; i < 20; i++ {
		var newUser *models.User
		newUser, err = services.Users.Create(context.Background(),
			fmt.Sprintf("TestDefaultUser%d", i),
			password,
			models.DefaultRole,
		)
		if err != nil {
			return nil, err
		}

		fixtureData.AmountOfUsers++

		fixtureData.DefaultUsers = append(fixtureData.DefaultUsers, *newUser)
	}

	for i := 0; i < 10; i++ {
		var newUser *models.User
		newUser, err = services.Users.Create(context.Background(),
			fmt.Sprintf("TestManagerUser%d", i),
			password,
			models.ManagerRole,
		)
		if err != nil {
			return nil, err
		}

		fixtureData.AmountOfUsers++

		fixtureData.ManagerUsers = append(fixtureData.ManagerUsers, *newUser)
	}

	adminAccessToken, err := services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		adminUser.ID,
		cfg.AccessExpiry,
		false,
	)
	if err != nil {
		return nil, err
	}

	managerAccessToken, err := services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		managerUser.ID,
		cfg.AccessExpiry,
		false,
	)
	if err != nil {
		return nil, err
	}

	defaultAccessToken, err := services.Auth.CreateCookie(context.Background(),
		models.AccessScope,
		defaultUser.ID,
		cfg.AccessExpiry,
		false,
	)
	if err != nil {
		return nil, err
	}

	defaultRefreshToken, err := services.Auth.CreateCookie(context.Background(),
		models.RefreshScope,
		defaultUser.ID,
		cfg.RefreshExpiry,
		false,
	)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AdminAccessToken:    *adminAccessToken,
		ManagerAccessToken:  *managerAccessToken,
		DefaultAccessToken:  *defaultAccessToken,
		DefaultRefreshToken: *defaultRefreshToken,
	}, nil
}

func locationFixtures(services services.Services) error {
	locations, err := services.Locations.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, location := range locations {
		err = services.Locations.Delete(context.Background(), location.ID)
		if err != nil {
			return err
		}
	}

	fixtureData.AmountOfLocations = 0

	fixtureData.DefaultLocation, err = services.Locations.Create(context.Background(),
		"TestLocation", 20, fixtureData.DefaultUser.ID)
	if err != nil {
		return err
	}
	fixtureData.AmountOfLocations++

	for i := 0; i < 5; i++ {
		_, err = services.CheckIns.Create(
			context.Background(),
			fixtureData.DefaultLocation.ID,
			1,
			fixtureData.DefaultLocation.Capacity+int64(i),
		)
		if err != nil {
			return err
		}
	}

	fixtureData.DefaultLocation.Available -= 5

	for i := 0; i < 20; i++ {
		var location *models.Location
		location, err = services.Locations.Create(context.Background(),
			fmt.Sprintf("TestLocation%d", i), 20, fixtureData.DefaultUsers[i].ID)
		if err != nil {
			return err
		}
		fixtureData.AmountOfLocations++

		for j := 0; j < 5; j++ {
			_, err = services.CheckIns.Create(context.Background(),
				location.ID, 1, location.Capacity)
			if err != nil {
				return err
			}
		}

		location.Available -= 5

		fixtureData.Locations = append(fixtureData.Locations, *location)
	}

	return nil
}

func schoolFixtures(services services.Services) error {
	schools, err := services.Schools.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, school := range schools {
		if school.ID == 1 {
			continue
		}

		err = services.Schools.Delete(context.Background(), school.ID)
		if err != nil {
			return err
		}
	}

	fixtureData.AmountOfSchools = 0

	for i := 0; i < 20; i++ {
		var school *models.School
		school, err = services.Schools.Create(context.Background(),
			fmt.Sprintf("TestSchool%d", i))
		if err != nil {
			return err
		}
		fixtureData.Schools = append(fixtureData.Schools, *school)
		fixtureData.AmountOfSchools++
	}

	return nil
}

func fixtures(tx database.DB) Tokens {
	services := services.New(tx)

	tokens, err := userFixtures(services)
	if err != nil {
		panic(err)
	}

	err = locationFixtures(services)
	if err != nil {
		panic(err)
	}

	err = schoolFixtures(services)
	if err != nil {
		panic(err)
	}

	return *tokens
}

func TestMain(m *testing.M) {
	var err error

	cfg = config.New()
	cfg.Env = config.TestEnv

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	mainTestEnv, err = tests.SetupGlobal(
		cfg.DB.Dsn,
		cfg.DB.MaxConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		panic(err)
	}

	tokens = fixtures(mainTestEnv.TestTx)

	exitCode := m.Run()
	err = tests.TeardownGlobal(mainTestEnv)
	if err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

func setupTest(
	_ *testing.T,
	mainTestEnv *tests.MainTestEnv,
) (tests.TestEnv, *application) {
	// t.Parallel()
	testEnv := tests.SetupSingle(mainTestEnv)

	testApp := &application{
		config:   cfg,
		logger:   logger,
		services: services.New(testEnv.TestTx),
	}

	return testEnv, testApp
}
