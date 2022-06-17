FROM golang:1.18-buster as builder

WORKDIR /queueBot
COPY go.* ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bot cmd/queueBot/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /bot /bot

CMD ["/bot"]