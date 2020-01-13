APP=favorite-tree

all: test build deploy

## test: run tests
test:
	go test ./...

## build: build docker-image in minikube
build:
	eval $$(minikube docker-env --shell='$(SHELL)'); \
	docker build . -t $(APP):test

## deploy: apply kubernetes manifests files in the current context
deploy:
	kubectl apply -f .k8s/
	echo ""
	echo "Waiting for service availability..."
	kubectl wait --for=condition=available --timeout=60s deployment/$(APP)
	echo ""
	echo "Service is available. You can test in using curl:"
	echo "curl -X GET http://$$(minikube ip):80/tree -H 'Host:local.ecosia.org'"

## destroy: delete entities deployed to kubernetes cluster
destroy:
	kubectl delete -f .k8s/

## help: print this information
help: Makefile
	echo ' Choose a command to run in ecosia-test:'
	sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'

.PHONY: test build deploy destroy
.SILENT: test build deploy destroy help
