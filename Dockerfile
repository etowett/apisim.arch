#Compile stage
FROM golang:1.16.2-alpine AS builder

# Add required packages
RUN apk add  --no-cache --update git curl bash

RUN go get -u github.com/revel/revel
RUN go get -u github.com/revel/cmd/revel

WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
ENV CGO_ENABLED 0 \
    GOOS=linux \
    GOARCH=amd64
ADD . .

RUN revel package .

# Run stage
FROM alpine:3.13
RUN apk update && \
    apk add mailcap tzdata && \
    rm /var/cache/apk/*
WORKDIR /app
COPY --from=builder /app/app.tar.gz .
RUN tar -xzvf app.tar.gz && rm app.tar.gz
ENTRYPOINT /app/run.sh
