FROM golang:1.12.6-alpine3.9 AS build-env

RUN apk add --no-cache build-base git

ENV GO111MODULE=on

COPY . /tmp/src

# Test
RUN cd /tmp/src && \
    go test -mod=vendor -short ./...

# Compile
RUN cd /tmp/src && \
    go -mod=vendor build


FROM alpine:3.9.4

COPY --from=build-env /tmp/src/ohammer /usr/local/bin/ohammer

RUN ohammer --version

CMD ohammer
