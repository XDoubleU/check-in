# Check-In API

## Lint Commands

(Prereq.) Install linter packages:  `make lint/init`
Run linter:                         `make lint`
Run autofix linter:                 `make lint/fix`
Clean mod file:                     `go mod tidy`

## Dev Commands

(Prereq.) Run database only:  `docker-compose up -d`
Run API:                      `make run/api`
Run CLI (for creating admin): `make run/cli/createadmin u=[USERNAME] p=[PASSWORD]`

## DB Migration Commands

New migration:  `make db/migrations/new name=[NAME]`
Apply:          `make db/migrations/up`
Undo:           `make db/migrations/down`

## Build Commands

Build: `make build`
Run from build:

- Windows: `bin\api`

## Test Commands

Running tests:                      `make test`
Running tests verbose:              `make test/v`
Generate coverage report of tests:  `make test/cov`
Reopen coverage report:             `make test/cov/open`
Run Artillery stress test:          `npx artillery run artillery/checkin.yml`

## Other Commands

Generate OpenAPI Spec: `make swag`

## CI Commands

Generate test report: `make test/cov/report`
