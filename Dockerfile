FROM golang:1.13-alpine as build

WORKDIR /app


COPY backend/go.mod .
COPY backend/go.sum .

RUN go mod download

COPY backend .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app

FROM scratch

COPY --from=build app .
