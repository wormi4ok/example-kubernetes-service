# Example kubernetes service

This is a simple HTTP service that returns a hello message in JSON format.
It comes with a sample configuration to deploy to Kubernetes cluster.

## API

`GET /hello` - returns hello message in the format: `{"helloMsg":<TEXT>}`

## Requirements

To build an application:

* [golang](https://golang.org/dl/) v1 (tested with 1.14)

To deploy to minikube:

* [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/) with ingress add-on enabled
* [docker](https://docs.docker.com/install/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Deployment

Run the given command in the project root directory and follow the instructions:
```bash
make
```

To see all available options, use `make help`
