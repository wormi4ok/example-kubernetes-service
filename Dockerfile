FROM golang:1.15 AS builder

COPY . /go/src/example-kubernetes-service
WORKDIR /go/src/example-kubernetes-service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/weather-service .

FROM alpine:3.12

RUN addgroup -S wormi4ok && adduser -S wormi4ok -G wormi4ok
USER wormi4ok

COPY --from=builder /go/bin/weather-service /bin/weather-service

ENV PORT 8080
ENTRYPOINT ["weather-service"]
