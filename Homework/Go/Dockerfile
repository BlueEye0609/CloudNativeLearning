FROM golang:1.17-alpine

RUN mkdir /scripts
COPY httpServer/* /scripts/

expose 80

ENTRYPOINT [ "go run /scripts/main.go" ]
