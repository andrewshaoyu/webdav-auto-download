FROM golang:latest as builder
WORKDIR /app
COPY ./ ./
RUN GOOS=linux GOARCH=amd64 go build


FROM docker.io/alpine
COPY  --from=builder /app/webdav-manager webdev-manager

ENTRYPOINT ["./webdev-manager"]