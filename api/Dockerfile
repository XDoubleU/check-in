FROM golang:1.20 as base

WORKDIR /app

COPY Makefile go.mod go.sum ./
COPY ./internal/  ./internal/
COPY ./cmd/api/ ./cmd/api/

RUN go mod download
RUN make build

EXPOSE 8000

CMD ["/api"]