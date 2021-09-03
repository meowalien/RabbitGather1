FROM alpine:3.14.2
FROM golang:1.16.5
WORKDIR "/app"
ADD     . /app
RUN    go mod download
RUN    go build -o ./bin/theapp ./src/main.go
RUN    chmod +x ./bin/theapp
RUN    setcap 'cap_net_bind_service=+ep' ./bin/theapp

EXPOSE 2001

ENTRYPOINT "./bin/theapp"
