
test:
	go test ./...

build:
	go build -v -o ./favorite-tree

.PHONY: test build
