FROM golang:latest

WORKDIR /app

ENV TOKEN=2057091218:AAG4dtsDK-THlBjpD2CruAd-LhQIWbBnjx0
ENV HOST=localhost
ENV PORT=3306
ENV USER=root
ENV PASSWORD=14112002Rich
ENV DBNAME=queue_test_docker
ENV CREATE_TABLE=yes
ENV DB_SWITCH=on
ENV ADMINPASSWORD=qwerty

COPY go.* ./
RUN go mod download
COPY .. ./

RUN go build -o /docker-queue-bot

EXPOSE 3306/tcp

CMD [ "/docker-queue-bot" ]