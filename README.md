# My favorite tree 🌳

This is a simple HTTP service that returns an information about my favorite tree in JSON format.

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
