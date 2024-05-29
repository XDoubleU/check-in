//nolint:gomnd //no magic number
package config

import (
	"fmt"

	"github.com/XDoubleU/essentia/pkg/config"
)

type Config struct {
	Env           string
	Port          int
	Throttle      bool
	WebURL        string
	SentryDsn     string
	SampleRate    float64
	AccessExpiry  string
	RefreshExpiry string
	DB            struct {
		Dsn         string
		MaxConns    int
		MaxIdleTime string
	}
	Release string
}

const (
	ProdEnv string = "production"
	TestEnv string = "test"
	DevEnv  string = "development"
)

//nolint:forbidigo //returns value of config
func New() Config {
	var cfg Config

	cfg.Env = config.GetEnvStr("ENV", ProdEnv)
	fmt.Println("cfg.Env: ", cfg.Env)

	cfg.Port = config.GetEnvInt("PORT", 8000)
	fmt.Println("cfg.Port: ", cfg.Port)

	cfg.Throttle = config.GetEnvBool("THROTTLE", true)
	fmt.Println("cfg.Throttle: ", cfg.Throttle)

	cfg.WebURL = config.GetEnvStr("WEB_URL", "http://localhost:3000")
	fmt.Println("cfg.WebURL: ", cfg.WebURL)

	cfg.SentryDsn = config.GetEnvStr("SENTRY_DSN", "")
	fmt.Println("cfg.SentryDsn: ", cfg.SentryDsn)

	cfg.SampleRate = config.GetEnvFloat("SAMPLE_RATE", 1.0)
	fmt.Println("cfg.SampleRate: ", cfg.SampleRate)

	cfg.AccessExpiry = config.GetEnvStr("ACCESS_EXPIRY", "1h")
	fmt.Println("cfg.AccessExpiry: ", cfg.AccessExpiry)

	cfg.RefreshExpiry = config.GetEnvStr("REFRESH_EXPIRY", "7d")
	fmt.Println("cfg.RefreshExpiry: ", cfg.RefreshExpiry)

	cfg.DB.Dsn = config.GetEnvStr("DB_DSN", "postgres://postgres@localhost/postgres")
	fmt.Println("cfg.DB.Dsn: ", cfg.DB.Dsn)

	cfg.DB.MaxConns = config.GetEnvInt("DB_MAX_CONNS", 25)
	fmt.Println("cfg.DB.MaxConns: ", cfg.DB.MaxConns)

	cfg.DB.MaxIdleTime = config.GetEnvStr("DB_MAX_IDLE_TIME", "15m")
	fmt.Println("cfg.DB.MaxIdleTime: ", cfg.DB.MaxIdleTime)

	cfg.Release = config.GetEnvStr("RELEASE", DevEnv)
	fmt.Println("cfg.Release: ", cfg.Release)

	return cfg
}
