#Compile stage
FROM golang:1.16.0-alpine AS build

# Add required packages
RUN apk add  --no-cache --update git curl bash

WORKDIR /app
ADD . .
ENV CGO_ENABLED 0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go get -u github.com/revel/revel
RUN go get -u github.com/revel/cmd/revel
RUN go mod download

RUN revel build apisim apisim dev

# Run stage
FROM alpine:3.13.2
RUN apk update && \
    apk add mailcap tzdata && \
    rm /var/cache/apk/*
WORKDIR /apisim
COPY --from=builder /app/apisim .
ENTRYPOINT /apisim/run.sh
