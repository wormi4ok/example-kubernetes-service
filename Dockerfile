FROM golang:1.13 AS builder

COPY . /go/src/favorite-tree
WORKDIR /go/src/favorite-tree

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/favorite-tree .

FROM alpine:3.11

RUN addgroup -S ecosia && adduser -S ecosia -G ecosia
USER ecosia

COPY --from=builder /go/bin/favorite-tree /bin/favorite-tree

ENV PORT 8080
ENTRYPOINT ["favorite-tree"]
