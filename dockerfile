FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

ADD cmd ./cmd

ADD pkg ./pkg

RUN go build -o ./main ./cmd/main.go

COPY config.json ./

ENTRYPOINT [ "/app/main" ]