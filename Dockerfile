FROM golang:1.23 AS builder


WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download


COPY ./backend ./backend


WORKDIR /app/backend/cmd


RUN CGO_ENABLED=0 go build -o /app/app .


FROM debian:bookworm-slim


WORKDIR /app


COPY --from=builder /app/app .


RUN chmod +x /app/app


EXPOSE 8080


CMD ["./app"]
