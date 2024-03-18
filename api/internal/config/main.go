//nolint:gomnd //no magic number
package config

import "fmt"

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
	var config Config

	config.Env = GetEnvStr("ENV", ProdEnv)
	fmt.Println("config.Env: ", config.Env)

	config.Port = GetEnvInt("PORT", 8000)
	fmt.Println("config.Port: ", config.Port)

	config.Throttle = GetEnvBool("THROTTLE", true)
	fmt.Println("config.Throttle: ", config.Throttle)

	config.WebURL = GetEnvStr("WEB_URL", "http://localhost:3000")
	fmt.Println("config.WebURL: ", config.WebURL)

	config.SentryDsn = GetEnvStr("SENTRY_DSN", "")
	fmt.Println("config.SentryDsn: ", config.SentryDsn)

	config.SampleRate = GetEnvFloat("SAMPLE_RATE", 1.0)
	fmt.Println("config.SampleRate: ", config.SampleRate)

	config.AccessExpiry = GetEnvStr("ACCESS_EXPIRY", "1h")
	fmt.Println("config.AccessExpiry: ", config.AccessExpiry)

	config.RefreshExpiry = GetEnvStr("REFRESH_EXPIRY", "7d")
	fmt.Println("config.RefreshExpiry: ", config.RefreshExpiry)

	config.DB.Dsn = GetEnvStr("DB_DSN", "postgres://postgres@localhost/postgres")
	fmt.Println("config.DB.Dsn: ", config.DB.Dsn)

	config.DB.MaxConns = GetEnvInt("DB_MAX_CONNS", 25)
	fmt.Println("config.DB.MaxConns: ", config.DB.MaxConns)

	config.DB.MaxIdleTime = GetEnvStr("DB_MAX_IDLE_TIME", "15m")
	fmt.Println("config.DB.MaxIdleTime: ", config.DB.MaxIdleTime)

	config.Release = GetEnvStr("RELEASE", DevEnv)
	fmt.Println("config.Release: ", config.Release)

	return config
}
