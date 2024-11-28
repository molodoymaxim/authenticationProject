FROM golang:1.23-alpine AS builder

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --dir ./cmd/ --output ./docs

RUN go build -o auth-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-service .
COPY ./migrations ./migrations
COPY .env .env

ENV APP_ENV=production

EXPOSE 8080

# Запуск приложения
CMD ["./auth-service"]
