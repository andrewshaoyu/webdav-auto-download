FROM golang:latest as builder
COPY ./ ./
RUN go build


FROM docker.io/alpine
COPY  --from=builder /webdav-manager webdev-manager

ENTRYPOINT ["./webdev-manager"]