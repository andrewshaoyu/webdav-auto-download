FROM golang:latest as builder
WORKDIR /app
COPY ./ /app
RUN go mod tidy
RUN go build -o webdav-manager


FROM golang:latest
COPY  --from=builder /app/webdav-manager /webdav-manager
ENTRYPOINT ["./webdav-manager"]