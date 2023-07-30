FROM node:18.16-alpine AS deps
RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --link package.json yarn.lock ./
RUN yarn install --frozen-lockfile --production && yarn cache clean

FROM node:18.16-alpine AS builder
WORKDIR /app

ARG NODE_ENV
ARG PORT
ARG NEXT_PUBLIC_API_URL
ARG SENTRY_RELEASE
ARG NEXT_PUBLIC_SENTRY_DSN
ARG SENTRY_ORG
ARG SENTRY_PROJECT
ARG SENTRY_AUTH_TOKEN

COPY --from=deps --link /app/node_modules ./node_modules
COPY --link  . .
RUN yarn build

FROM nginx:alpine AS server

COPY --link ./nginx-docker.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/html

RUN rm -rf ./*

COPY --from=builder --link /app/out /usr/share/nginx/html