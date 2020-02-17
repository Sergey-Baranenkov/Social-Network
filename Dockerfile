FROM golang:1.13-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app .
