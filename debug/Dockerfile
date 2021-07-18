# Compile stage
FROM golang:1.16.6 AS build-env

# Build Delve
RUN go get github.com/go-delve/delve/cmd/dlv

COPY . /usr/local/go/src/debug
WORKDIR /usr/local/go/src/debug

RUN go build -gcflags="all=-N -l" -o server

# Final stage
FROM debian:buster

EXPOSE 8000 40000

WORKDIR /
COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /usr/local/go/src/debug/server /


ENV HTTP_PORT 8080
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]