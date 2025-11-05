# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notification-service ./cmd/server

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build /app/notification-service /app/notification-service

ENV HTTP_PORT=8089

EXPOSE 8089

ENTRYPOINT ["/app/notification-service"]
