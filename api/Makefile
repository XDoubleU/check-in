db ?= postgres://postgres@localhost/postgres

init:
	go install github.com/segmentio/golines@latest
	go install github.com/daixiang0/gci@latest
	go install github.com/swaggo/swag/cmd/swag@master

lint:
	golangci-lint run
	sqlfluff lint --dialect postgres ./migrations

lint/fix:
	swag fmt
	golines . -m 88 -w
	golangci-lint run --fix
	gci write --skip-generated -s standard -s default -s "prefix(check-in/api)" .
	sqlfluff fix --dialect postgres ./migrations

build:
	go build -o=./bin/api ./cmd/api

run/api:
	go run ./cmd/api

run/cli/createadmin:
	go run ./cmd/cli -db=${db} -u=${u} -p=${p} createadmin

test:
	go test ./cmd/api

test/v:
	go test -v ./cmd/api

test/cov/report:
	go test -covermode=set -coverprofile=coverage.out ./cmd/api -coverpkg ./cmd/api,./internal/...

test/cov:
	go test -covermode=set -coverprofile=coverage.out ./cmd/api -coverpkg ./cmd/api,./internal/...
	go tool cover -html=coverage.out -o=coverage.html
	make test/cov/open

test/cov/open:
	CMD /C start chrome /new-tab %CD%/coverage.html

db/migrations/new:
	goose -s -dir=./migrations create ${name} sql

db/migrations/up:
	goose -dir=./migrations postgres ${db} up

db/migrations/down:
	goose -dir=./migrations postgres ${db} down

swag:
	swag init --ot json --parseDependency -g cmd/api/main.go 