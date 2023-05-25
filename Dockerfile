FROM golang/golang:lastest as builder
WORKDIR /app
COPY ./ ./
RUN GOOS=linux GOARCH=amd64 go build


FROM docker.io/alpine
LABEL authors="shaoyu"
WORKDIR /app
COPY webdav-manager webdev-manager

ENTRYPOINT ["./app/webdev-manager"]