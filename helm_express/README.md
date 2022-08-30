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

## Demo

![Screen Shot 2022-08-30 at 10 40 32 AM](https://user-images.githubusercontent.com/9366404/187336590-a5bc47ad-d584-481f-aa92-b62a22d65eab.png)


1. Video1 - PR Walkthrough: https://www.loom.com/share/cd80c0d019c04e49832b83e275fae8f9
2. Video2 - Helm Walkthrough: https://www.loom.com/share/e681bb14f0104a4e8b91543e0b2ec24e
3. Video3 - API Walkthrough: https://www.loom.com/share/b38e209d61064f61a6419fb725157d97

## Config

Please check [values.yaml](helm-express/values.yaml) for detail!

## CI

Please refer to [.github/workflows/node_ci.yaml](.github/workflows/node_ci.yaml) for details!
