module check-in/api

go 1.22

toolchain go1.22.4

require github.com/justinas/alice v1.2.0

require (
	github.com/dlclark/regexp2 v1.11.0
	github.com/getsentry/sentry-go v0.28.1
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.9.0
	github.com/wk8/go-ordered-map/v2 v2.1.8
	github.com/xdoubleu/essentia v0.0.2
	github.com/xhit/go-str2duration/v2 v2.1.0
)

replace github.com/xdoubleu/essentia v0.0.2 => ../../essentia

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/goddtriffin/helmet v1.0.2 // indirect
	github.com/jackc/pgerrcode v0.0.0-20240316143900-6e2875d9b438 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	nhooyr.io/websocket v1.8.11 // indirect
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.6.0
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/rs/cors v1.11.0 // indirect
	golang.org/x/crypto v0.24.0
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)
