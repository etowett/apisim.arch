#Compile stage
FROM golang:1.16.2-alpine AS build

# Add required packages
RUN apk add  --no-cache --update git curl bash

RUN go get -u github.com/revel/revel
RUN go get -u github.com/revel/cmd/revel

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

ENV CGO_ENABLED 0 \
    GOOS=linux \
    GOARCH=amd64

ADD . .

RUN revel build apisim apisim dev

# Run stage
FROM alpine:3.13.2
RUN apk update && \
    apk add mailcap tzdata && \
    rm /var/cache/apk/*
WORKDIR /apisim
COPY --from=builder /app/apisim .
ENTRYPOINT /apisim/run.sh
