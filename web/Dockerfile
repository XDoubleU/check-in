# stage-0: generate static HTML files
FROM node:18.16-alpine

WORKDIR /app

COPY --link package.json yarn.lock
RUN yarn --frozen-lockfile

COPY --link  . .
RUN yarn export

# stage-1: serve static HTML files
FROM nginx:alpine

COPY --link ./nginx-docker.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/html

RUN rm -rf ./*

COPY --from=0 /app/out /usr/share/nginx/html