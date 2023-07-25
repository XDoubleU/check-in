# stage-0: deps
FROM node:18.16-alpine AS deps
WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --link package.json yarn.lock ./
RUN --mount=type=cache,target=/root/.yarn YARN_CACHE_FOLDER=/root/.yarn yarn install --frozen-lockfile

# stage-1: generate static HTML files
FROM node:18.16-alpine AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY --link  . .
RUN yarn export

# stage-2: serve static HTML files
FROM nginx:alpine AS server

COPY --link ./nginx-docker.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/html

RUN rm -rf ./*

COPY --from=builder /app/out /usr/share/nginx/html