APP=favorite-tree

all: test build deploy

test:
	go test ./...

build:
	eval $$(minikube docker-env --shell='$(SHELL)'); \
	docker build . -t $(APP):test

deploy:
	kubectl apply -f .k8s/
	echo ""
	echo "Waiting for service availability..."
	kubectl wait --for=condition=available --timeout=60s deployment/$(APP)
	echo ""
	echo "Service is available. You can test in using curl:"
	echo "curl -X GET http://$$(minikube ip):80/tree -H 'Host:local.ecosia.org'"

destroy:
	kubectl delete -f .k8s/

.PHONY: test build deploy destroy
