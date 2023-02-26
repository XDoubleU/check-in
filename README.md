TODO

## How to run?

1. `docker-compose up -d --build` for running the db, api and web-client
2. Go to `http://localhost:3000` for the web-client and  `http://localhost:8000` for the api

Below commands work only on api

## Building
3. Build: `turbo build`

## Linting
1. Linting: `turbo lint` and `turbo lint:fix`

## Edit schema?
1. `turbo db:generate`
2. Create migrations: `turbo db:migrate-dev -- --name [NAME]` (run with docker-compose)

## Seeding database
1. `turbo db:seed`

## Other
1. `npm run cli createadmin` create admin


## Commands in Docker
1. `docker-compose exec api npx turbo [cmd]`
2. Provide arguments to pass after '--': `docker-compose exec api npx turbo [cmd] -- [args]`
3. Create admin: `docker-compose exec api npx turbo cli -- createadmin -u username -p password`

## Run tests
1. Start docker: `docker-compose up -d --build`
2. Migrate database schema: `docker-compose exec api npx turbo db:test`
3. Run tests on api: `docker-compose exec api npx turbo test`
4. Run tests on web: `docker-compose exec web npx turbo test`