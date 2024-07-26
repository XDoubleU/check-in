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

func New() Config {
	var cfg Config

	cfg.Env = config.GetEnvStr("ENV", ProdEnv)
	cfg.Port = config.GetEnvInt("PORT", 8000)
	cfg.Throttle = config.GetEnvBool("THROTTLE", true)
	cfg.WebURL = config.GetEnvStr("WEB_URL", "http://localhost:3000")
	cfg.SentryDsn = config.GetEnvStr("SENTRY_DSN", "")
	cfg.SampleRate = config.GetEnvFloat("SAMPLE_RATE", 1.0)
	cfg.AccessExpiry = config.GetEnvStr("ACCESS_EXPIRY", "1h")
	cfg.RefreshExpiry = config.GetEnvStr("REFRESH_EXPIRY", "7d")
	cfg.DB.Dsn = config.GetEnvStr("DB_DSN", "postgres://postgres@localhost/postgres")
	cfg.DB.MaxConns = config.GetEnvInt("DB_MAX_CONNS", 25)
	cfg.DB.MaxIdleTime = config.GetEnvStr("DB_MAX_IDLE_TIME", "15m")
	cfg.Release = config.GetEnvStr("RELEASE", DevEnv)

	return cfg
}

func (cfg Config) String() string {
	return fmt.Sprintf(`config:
	cfg.Env: %s
	cfg.Port: %d
	cfg.Throttle: %t
	cfg.WebURL: %s
	cfg.SentryDsn: %s
	cfg.SampleRate: %f
	cfg.AccessExpiry: %s
	cfg.RefreshExpiry: %s
	cfg.DB.Dsn: %s
	cfg.DB.MaxConns: %d
	cfg.DB.MaxIdleTime: %s
	cfg.Release: %s`,
		cfg.Env,
		cfg.Port,
		cfg.Throttle,
		cfg.WebURL,
		cfg.SentryDsn,
		cfg.SampleRate,
		cfg.AccessExpiry,
		cfg.RefreshExpiry,
		cfg.DB.Dsn,
		cfg.DB.MaxConns,
		cfg.DB.MaxIdleTime,
		cfg.Release,
	)
}
