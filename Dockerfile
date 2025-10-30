FROM golang:1.23.2-alpine3.20 AS build

RUN apk --update add --no-cache \
    gcc \
    musl-dev \
    git

RUN mkdir /build

COPY . /build

WORKDIR /build

RUN go build -ldflags "-s -w"

FROM alpine:3.20

RUN apk --update add --no-cache \
    curl \
    iproute2-ss

COPY --from=build /build/bitmagnet /usr/bin/bitmagnet

ENTRYPOINT ["bitmagnet"]
