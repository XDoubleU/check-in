FROM golang:1.20

WORKDIR /app

COPY --link Makefile go.mod go.sum ./
COPY --link ./internal/  ./internal/
COPY --link ./cmd/api/ ./cmd/api/

RUN go mod download
RUN make build

EXPOSE 8000
CMD ["./bin/api"]