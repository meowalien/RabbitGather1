FROM centos
FROM golang:1.16.5
WORKDIR "/app"
ADD     . /app
RUN    go mod download
RUN    go build -o ./bin/theapp ./src/main.go
RUN    chmod +x ./bin/theapp
RUN    setcap 'cap_net_bind_service=+ep' ./bin/theapp

EXPOSE 443

ENTRYPOINT "./bin/theapp"
