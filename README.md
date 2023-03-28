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

## Static testing

1. Linting: `pnpm lint` and `pnpm lint:fix`
2. Knip: `pnpm knip`
3. Duplicates: `pnpm jscpd`
4. Formatting: `pnpm format`

## Edit schema?

1. Create migrations: `pnpm db:migration-create`
2. Apply migration: `pnpm db:migration-up`
3. Undo migration: `pnpm db:migration-down`

## Other

1. `pnpm cli createadmin` create admin

## Run tests

1. Start database in docker: `docker-compose up -d`
2. Setup database: `pnpm db:test`
3. Run tests on: `pnpm test`

## Deploy (on DigitalOcean)

Mostly as a reference to myself but might be useful for others too.
CI needs step where it creates a package-lock.json using `npm i --package-lock-only`

0. Database with pooling
1. Don't forget env vars (see .env)
2. Web (static site): `npm run export --filter=web...`
3. API (server):
   1. Build: `npm run build --filter=api...`
   2. Run: `npm run prod --filter=api...`
4. Manually on API: `npm db:migration-up`
