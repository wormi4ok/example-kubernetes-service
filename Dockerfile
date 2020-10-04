FROM golang:1.15.2-alpine3.12 AS builder

COPY . /go/src/example-kubernetes-service
WORKDIR /go/src/example-kubernetes-service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o /go/bin/weather-service .

# -----------------------------------------------------------------------------

FROM alpine:3.12 as certificates

RUN apk add -U --no-cache ca-certificates

# -----------------------------------------------------------------------------

FROM scratch

LABEL org.opencontainers.image.title="weather-service"
LABEL org.opencontainers.image.description="Jodel SRE Tech Challenge"
LABEL org.opencontainers.image.authors="Stanislav Petrashov"

COPY --from=builder /go/bin/weather-service /bin/weather-service
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/bin/weather-service"]
