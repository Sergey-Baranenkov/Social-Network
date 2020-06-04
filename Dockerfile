FROM golang:1.13-alpine as build

WORKDIR /app


COPY backend/go.mod .
COPY backend/go.sum .

RUN go mod download

COPY backend backend
COPY frontend frontend

WORKDIR ./backend

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app

WORKDIR ../


FROM scratch
COPY --from=build app/backend/app .
COPY --from=build app/frontend frontend