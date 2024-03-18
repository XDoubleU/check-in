FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY --link go.mod go.sum ./
RUN go mod download

COPY --link ./internal/  ./internal/
COPY --link ./cmd/api/ ./cmd/api/
RUN CGO_ENABLED=0 go build -o=./bin/api ./cmd/api

FROM scratch AS server
WORKDIR /app

COPY --from=builder --link /app/bin/api ./bin/api

EXPOSE 8000
CMD ["./bin/api"]