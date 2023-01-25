FROM golang:1.18.10-alpine3.17 AS builder

WORKDIR /usr/src/app/

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
RUN apk add pkgconfig zeromq-dev zeromq build-base
COPY handler/go.mod handler/go.sum ./
RUN go mod download && go mod verify
COPY handler ./
RUN go build -v

FROM alpine:3.17

WORKDIR /usr/src/app/
RUN apk add zeromq
COPY --from=builder /usr/src/app/loadHandler ./
