FROM golang:1.18-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bot ./cmd

FROM alpine:3.15.4
WORKDIR /app
COPY --from=builder /app/bot /app/bot
COPY --from=builder /app/.env /app/.env

CMD ["sh", "-c", "/app/bot"]