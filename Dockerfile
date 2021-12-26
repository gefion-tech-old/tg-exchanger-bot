FROM golang:1.17-alpine

WORKDIR /bot

COPY . /bot/

RUN go mod download

RUN go build -o bot -v ./cmd/main

CMD mkdir logs

CMD ./bot -prod true