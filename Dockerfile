FROM golang:latest

WORKDIR app/

COPY go.mod ./
COPY go.sum ./

RUN go get download

COPY bot/ ./
COPY main.go ./

RUN go build -o "/botApp"

EXPOSE 8080

CMD ["/botApp"]

