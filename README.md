# Check-In

## How to run ?

1. `docker-compose up -d` for running the database
2. `pnpm start` or `pnpm dev` for running the api and web apps
3. Go to `http://localhost:3000` for the web-client and `http://localhost:8000` for the api

## How to deploy ?

### Docker

There are Docker-files present in both apps.

### Commands

1. Building web: `pnpm build --filter=web...`
2. Building api: `pnpm build --filter=api...`
3. For running both apps: `pnpm prod`

## Linting

1. Linting: `pnpm lint` and `pnpm lint:fix`

## Edit schema?

1. `pnpm db:generate`
2. Create migrations: `pnpm db:migrate-dev -- --name [NAME]` (run with docker-compose)

## Other

1. `pnpm cli createadmin` create admin

## Run tests

1. Start docker: `docker-compose up -d --build`
2. Migrate database schema: `pnpm db:test`
3. Run tests on: `pnpm test`
