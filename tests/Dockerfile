FROM golang:1.13-alpine as build

WORKDIR /app_tests
COPY main_test.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -c main_test.go -o testik

FROM alpine
COPY --from=build /app_tests/testik .

CMD ["./testik","-test.v"]
