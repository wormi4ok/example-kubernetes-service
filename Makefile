APP=favorite-tree

all: test build deploy

test:
	go test ./...

build:
	eval $$(minikube docker-env --shell='$(SHELL)'); \
	docker build . -t $(APP):test

deploy:
	kubectl apply -f .k8s/

destroy:
	kubectl delete -f .k8s/

.PHONY: test build deploy destroy
