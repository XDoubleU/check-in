//nolint:gomnd //no magic number
package config

import (
	"github.com/XDoubleU/essentia/pkg/config"
	"github.com/XDoubleU/essentia/pkg/logger"
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

func (cfg Config) Print() {
	logger.GetLogger().Println("cfg.Env: ", cfg.Env)
	logger.GetLogger().Println("cfg.Port: ", cfg.Port)
	logger.GetLogger().Println("cfg.Throttle: ", cfg.Throttle)
	logger.GetLogger().Println("cfg.WebURL: ", cfg.WebURL)
	logger.GetLogger().Println("cfg.SentryDsn: ", cfg.SentryDsn)
	logger.GetLogger().Println("cfg.SampleRate: ", cfg.SampleRate)
	logger.GetLogger().Println("cfg.AccessExpiry: ", cfg.AccessExpiry)
	logger.GetLogger().Println("cfg.RefreshExpiry: ", cfg.RefreshExpiry)
	logger.GetLogger().Println("cfg.DB.Dsn: ", cfg.DB.Dsn)
	logger.GetLogger().Println("cfg.DB.MaxConns: ", cfg.DB.MaxConns)
	logger.GetLogger().Println("cfg.DB.MaxIdleTime: ", cfg.DB.MaxIdleTime)
	logger.GetLogger().Println("cfg.Release: ", cfg.Release)
}
