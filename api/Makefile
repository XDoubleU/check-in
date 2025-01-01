db ?= postgres://postgres@localhost/postgres
UNAME_S := $(shell uname -s)

tools: tools/db tools/lint tools/swagger
	
tools/swagger:
	go install github.com/swaggo/swag/cmd/swag@v1.16.4

tools/db:
	go install github.com/pressly/goose/v3/cmd/goose@v3.23.0

tools/lint: tools/lint/go tools/lint/sql

tools/lint/go:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.1
	go install github.com/segmentio/golines@v0.12.2
	go install github.com/daixiang0/gci@v0.13.5
	go install github.com/securego/gosec/v2/cmd/gosec@v2.21.4

tools/lint/sql:
ifeq ($(UNAME_S),Darwin)
	brew install sqlfluff
endif
ifeq ($(UNAME_S),Linux)
	pip install sqlfluff
endif

lint/sql: tools/lint/sql
	sqlfluff lint --dialect postgres ./cmd/api/migrations

lint: tools/lint
	golangci-lint run
	make lint/sql

lint/fix: tools/swagger
	swag fmt
	golines . -m 88 -w
	golangci-lint run --fix
	gci write --skip-generated -s standard -s default -s "prefix(check-in/api)" .
	sqlfluff fix --dialect postgres ./cmd/api/migrations

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
	go test ./cmd/api -v
	
test/race:
	go test ./cmd/api -race -v

test/pprof:
	go test ./cmd/api -cpuprofile cpu.prof -memprofile mem.prof -bench ./cmd/api

test/cov/report:
	go test ./cmd/api -coverpkg=./cmd/api,./internal/... -covermode=set -coverprofile=coverage.out

test/cov: test/cov/report
	go tool cover -html=coverage.out -o=coverage.html
	make test/cov/open

test/cov/open:
	open ./coverage.html

db/migrations/new: tools/db
	goose -s -dir=./migrations create ${name} sql

swag: tools/swagger
	swag init --ot json --parseDependency -g cmd/api/main.go 
	cp -R ./docs ./../web/docs 
