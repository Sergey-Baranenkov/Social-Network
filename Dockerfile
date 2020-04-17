FROM golang:1.13-alpine as build

WORKDIR /app

COPY go.mod . 
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app

FROM scratch

COPY --from=build app .
