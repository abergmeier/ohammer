FROM golang:1.12-alpine as build-env

RUN apk add --no-cache build-base git

COPY /usr/local/share/ca-certificates/* /usr/local/share/ca-certificates/

RUN update-ca-certificates

ENV GO111MODULE=on

# Initialize cache for remote sources
COPY go.mod /tmp/src/go.mod
COPY go.sum /tmp/src/go.sum
RUN cd /tmp/src/ && \
    go mod download

COPY . /tmp/src

# Test
RUN cd /tmp/src && \
    go test

# Compile
RUN cd /tmp/src && \
    go build

FROM alpine

COPY --from=build-env /tmp/src/ohammer /usr/local/bin/ohammer

RUN ohammer --version

CMD ohammer
