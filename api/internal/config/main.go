//nolint:gomnd //no magic number
package config

type Config struct {
	Env           string
	Port          int
	Throttle      bool
	WebURL        string
	SentryDsn     string
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
	var config Config

	config.Env = GetEnvStr("ENV", ProdEnv)
	config.Port = GetEnvInt("PORT", 8000)
	config.Throttle = GetEnvBool("THROTTLE", true)
	config.WebURL = GetEnvStr("WEB_URL", "http://localhost:3000")
	config.SentryDsn = GetEnvStr("SENTRY_DSN", "")

	config.AccessExpiry = GetEnvStr("ACCESS_EXPIRY", "1h")
	config.RefreshExpiry = GetEnvStr("REFRESH_EXPIRY", "7d")

	config.DB.Dsn = GetEnvStr("DB_DSN", "postgres://postgres@localhost/postgres")
	config.DB.MaxConns = GetEnvInt("DB_MAX_CONNS", 25)
	config.DB.MaxIdleTime = GetEnvStr("DB_MAX_IDLE_TIME", "15m")

	config.Release = GetEnvStr("RELEASE", DevEnv)

	return config
}
