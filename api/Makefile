db ?= postgres://postgres@localhost/postgres

tools: tools/db tools/lint tools/swagger
	
tools/swagger:
	go install github.com/swaggo/swag/cmd/swag@v1.16.3

tools/db:
	go install github.com/pressly/goose/v3/cmd/goose@v3.21.1

tools/lint: tools/lint/go tools/lint/sql

tools/lint/go:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	go install github.com/segmentio/golines@v0.12.2
	go install github.com/daixiang0/gci@v0.13.4
	go install github.com/securego/gosec/v2/cmd/gosec@v2.20.0

tools/lint/sql:
	pip install sqlfluff

lint/sql: tools/lint/sql
	sqlfluff lint --dialect postgres ./migrations

lint: tools/lint
	golangci-lint run
	make lint/sql

lint/fix: tools/swagger
	swag fmt
	golines . -m 88 -w
	golangci-lint run --fix
	gci write --skip-generated -s standard -s default -s "prefix(check-in/api)" .
	sqlfluff fix --dialect postgres ./migrations

build: 
	go build -o=./bin/api ./cmd/api
	make swag

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

db/migrations/new: tools/db
	goose -s -dir=./migrations create ${name} sql

db/migrations/up: tools/db
	goose -dir=./migrations postgres ${db} up

db/migrations/down: tools/db
	goose -dir=./migrations postgres ${db} down

swag: tools/swagger
	swag init --ot json --parseDependency -g cmd/api/main.go 
