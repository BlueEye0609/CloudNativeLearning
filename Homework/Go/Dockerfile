FROM golang:1.16-alpine AS build

WORKDIR /src/

COPY httpServer/main.go httpServer/go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/httpserver

FROM scratch
COPY --from=build /bin/httpserver /bin/httpserver
ENTRYPOINT [ "/bin/httpserver" ]