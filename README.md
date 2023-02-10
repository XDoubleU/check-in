TODO

## How to run?

1. `docker-compose up -d --build` for running the db, api and web-client
2. Go to `http://localhost:3000` for the web-client and  `http://localhost:8000` for the api

Below commands work only on api

## Edit schema?
1. `npx prisma generate`
2. `npx prisma migrate dev` or `npx prisma db push` (ONLY FOR PROTOTYPING)

## Other
1. `npx prisma db seed` seed database
2. `npm run cli createadmin` create admin
