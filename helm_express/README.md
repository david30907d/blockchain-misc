# Crypto Market Cap API Server (Helm + Express + Mongo + Redis) [![.github/workflows/node_ci.yml](https://github.com/david30907d/blockchain-misc/actions/workflows/node_ci.yml/badge.svg)](https://github.com/david30907d/blockchain-misc/actions/workflows/node_ci.yml)

## Prerequisites

1. minikube
2. docker-compose
3. node

## Install

Npm dependencies, for linter, formatter and commit linter (optional):

1. `brew install npm`
2. `npm ci`

## Run & Deploy

1. `minikube start`
2. Build image: `docker build -t davidtnfsh/helm_express:<tag> .`
3. Load those images you need into minikube:
   1. `minikube image load davidtnfsh/helm_express:<tag>`
   2. `minikube image load mongo:5.0.9`
   3. `minikube image load redis:6.2-alpine`
4. helm install express helm-express
5. get service url: `minikube service express-helm-express --url`
6. (optional): `helm upgrade express helm-express/`

## Config

Please check [values.yaml](helm-express/values.yaml) for detail!

## CI

Please refer to [.github/workflows/node_ci.yaml](.github/workflows/node_ci.yaml) for details!
