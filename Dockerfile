#Compile stage
FROM golang:1.18.2-alpine AS builder

# Add required packages
RUN apk add  --no-cache --update git curl bash

WORKDIR /app

ADD go.mod go.sum ./
RUN go mod download

RUN go install github.com/revel/cmd/revel@v1.0.0

ADD . .

RUN revel build . /build

# Run stage
FROM alpine:3.15

RUN apk update && \
    apk add mailcap tzdata && \
    rm /var/cache/apk/*

WORKDIR /app

COPY --from=builder /build/ .

ENTRYPOINT ["/app/run.sh"]
