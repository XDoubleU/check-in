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

1. Create migrations: `pnpm db:migration-create`
2. Apply migration: `pnpm db:migration-up`
3. Undo migration: `pnpm db:migration-down`

## Other

1. `pnpm cli createadmin` create admin

## Run tests

1. Start database in docker: `docker-compose up -d`
2. Setup database: `pnpm db:migrate-test`
3. Run tests on: `pnpm test`
