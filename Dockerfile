FROM golang:1.14 AS builder

COPY . /go/src/example-kubernetes-service
WORKDIR /go/src/example-kubernetes-service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/example-kubernetes-service .

FROM alpine:3.11

RUN addgroup -S wormi4ok && adduser -S wormi4ok -G wormi4ok
USER wormi4ok

COPY --from=builder /go/bin/example-kubernetes-service /bin/example-kubernetes-service

ENV PORT 8080
ENTRYPOINT ["example-kubernetes-service"]
