# Simple CRUD Interface

Rewrite the README according to the application.

The task itself can be found [here](/TASK.md)

## Prerequisites

- [Docker](https://www.docker.com/get-started/)
- [Goose](https://github.com/pressly/goose)
- [Gosec](https://github.com/securego/gosec)

## Getting Started

1. Start database

```
## Via Makefile
make db

## Via Docker
docker-compose up -d db
```

2. Run migrations

```
## Via Makefile
make migrate-up

## Via Goose
DB_DRIVER=postgres
DB_STRING="host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
goose -dir ./migrations $(DB_DRIVER) $(DB_STRING) up
```

3. Run application

```
go run cmd/main.go
```

## API

The project features a simple CRUD API for users. It has a JSON logger middleware and an X-Api-Key middleware for authentication.

The project also contains terraform scripts for setting up AKS cluster in Azure. (Read more from readme at ./platform/terraform)

There is also k8s-ingress setup which allows any deployment in the AKS cluster to be accessible through https and with a DNS. (Read more from readme at ./platform/terraform)

The API itself has been deployed to the AKS cluster, manifests are in k8s folder.

## Pipeline

There are three workflows
- go - simple CI pipeline that runs 
  - build
  - test
  - vet
  - fmt
  - lint
  - gosec
- build-and-deploy
  - builds a docker image
  - pushes image to docker hub
  - runs database migrations with goose
  - applies kubernetes manifests to deploy API and postgres
- db-migrate-down
  - manually startable workflow in case a db rollback is required