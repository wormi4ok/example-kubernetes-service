# My favorite tree 🌳

This is a simple HTTP service that returns an information about my favorite tree in JSON format.
It comes with a sample configuration to deploy to Kubernetes cluster 

## API

`GET /tree` - returns my favorite tree in the format: `{"myFavouriteTree":<NAME>}`

## Requirements

To build an application

* golang v1 (tested with 1.13)

No dependency management needed. I decided not to use any external libraries, due to the simplicity of the API.

To deploy to minikube

* minikube with ingress add-on enabled
* docker
* kubectl

## Deployment

Run the given command in the project root directory and follow the instructions:
```bash
make
```

To see all available option, use `make help`

--------------

This repository is created as a solution for the technical test assignment for [Ecosia](https://ecosia.org)

![Birch tree](birch-tree.jpg)
Photo by [Peng Chen](https://unsplash.com/@austincppc) on Unsplash
