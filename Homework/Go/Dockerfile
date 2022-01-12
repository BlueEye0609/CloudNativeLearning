FROM golang:latest

RUN mkdir /scripts
COPY httpServer/* /scripts/

expose 80

ENTRYPOINT [ "go", "run", "/scripts/main.go" ]