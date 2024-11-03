FROM node:22-alpine AS deps
ENV YARN_VERSION=4.x

RUN apk add --no-cache libc6-compat
RUN corepack enable && corepack prepare yarn@${YARN_VERSION}

WORKDIR /app

COPY --link package.json yarn.lock .yarnrc.yml ./
COPY --link ./.yarn/releases ./.yarn/releases
RUN yarn install --immutable && yarn cache clean

FROM node:22-alpine AS builder
WORKDIR /app

ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_EDIT_TIME_ZONE
ARG NEXT_PUBLIC_SENTRY_DSN
ARG NEXT_PUBLIC_SENTRY_ENVIRONMENT
ARG SENTRY_RELEASE
ARG SENTRY_ORG
ARG SENTRY_PROJECT
ARG SENTRY_AUTH_TOKEN

COPY --from=deps --link ./app/node_modules ./node_modules
COPY --link  . .
RUN yarn build

FROM nginx:alpine AS server

COPY --link ./nginx-docker.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/html

RUN rm -rf ./*

COPY --from=builder --link /app/out /usr/share/nginx/html