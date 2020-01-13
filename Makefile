APP=favorite-tree

test:
	go test ./...

build:
	docker build . -t $(APP):test

.PHONY: test build
