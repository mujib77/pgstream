FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .
RUN go build -mod=vendor -o pgstream .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pgstream .

CMD ["./pgstream"]