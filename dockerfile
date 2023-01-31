# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /home/upay/GolandProjects/awesomeProject

COPY go.mod ./
COPY go.sum ./
COPY ./files ./files
COPY sample.yaml ./
COPY manifest.yaml ./


RUN go mod download

COPY *.go ./


RUN go build -o /docker-gs-ping9

EXPOSE 8080

CMD [ "/docker-gs-ping9" ]