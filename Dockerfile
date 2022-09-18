FROM golang:1.18-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o queueBot ./cmd/queueBot

FROM alpine:3.15.4

MAINTAINER Richard Ragusski <me@richardhere.dev>

WORKDIR /app
COPY --from=builder /app/queueBot /app/queueBot

CMD ["sh", "-c", "/app/queueBot"]