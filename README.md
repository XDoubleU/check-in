# Check-In

## How to run ?

1. `docker-compose up -d` for running the database
2. `yarn start` or `yarn dev` for running the api and web apps
3. Go to `http://localhost:3000` for the web-client and `http://localhost:8000` for the api

## How to deploy ?

### Docker

There are Docker-files present in both apps.

### Commands

1. Building web: `yarn build --filter=web...`
2. Building api: `yarn build --filter=api...`
3. For running both apps: `yarn prod`

## Static testing

1. Linting: `yarn lint` and `yarn lint:fix`
2. Knip: `yarn knip`
3. Duplicates: `yarn jscpd`
4. Formatting: `yarn format`

## Edit schema?

1. Create migrations: `yarn db:migration-create`
2. Apply migration: `yarn db:migration-up`
3. Undo migration: `yarn db:migration-down`

## Other

1. `yarn cli createadmin` create admin

## Run tests

1. Start database in docker: `docker-compose up -d`
2. Setup database: `yarn db:test`
3. Run tests on: `yarn test`

## Deploy (on DigitalOcean)

Mostly as a reference to myself but might be useful for others too.

1. Database with pooling
2. Don't forget env vars (see .env)
3. Web (static site): `yarn export --filter=web...`
4. API (server):
   1. Build: `yarn build --filter=api...`
   2. Run: `yarn prod --filter=api...`
5. 
6. Manually on API: `yarn db:migration-up`
