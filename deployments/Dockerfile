FROM golang:latest AS builder

WORKDIR /opt/go_skele

RUN apk add --no-cache bash make git

COPY . .

RUN go build -o go_skele

FROM alpine:3.15.5

WORKDIR /usr/local/bin

RUN apk add --no-cache tzdata

COPY --from=builder /opt/go_skele/go_skele .

COPY db/migrations db/migrations

EXPOSE 8080/tcp
