# Cruder – Simple User CRUD API

A simple user management CRUD API built with Go (Gin).  
Features include:
- JSON structured logging middleware
- API key authentication (`X-Api-Key`)
- Auto-generated Swagger documentation at https://cruder.sytes.net/swagger/index.html

The original task description can be found [here](./TASK.md).


## Prerequisites

- [Go](https://go.dev/learn/)
- [Docker](https://www.docker.com/get-started/)
- [Goose](https://github.com/pressly/goose) (migrations)
- [Gosec](https://github.com/securego/gosec) (security analysis)

## Getting Started

1. Make a copy of .env.example and name it .env

2. Start database

```
## Via Makefile
make db

## Via Docker
docker-compose up -d db
```

3. Run migrations

```
## Via Makefile
make migrate-up

## Via Goose
DB_DRIVER=postgres
DB_STRING="host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
goose -dir ./migrations $(DB_DRIVER) $(DB_STRING) up
```

4. Run application

```
go run cmd/main.go
```

5. Generate API documentation
   
```
make swagger
```

6. Run tests
   
```
make test
```

## Infrastructure

The project also contains terraform scripts for setting up an AKS cluster in Azure. ([Read more](./platform/terraform/README.md)) 

There is also k8s-ingress setup which allows any deployment in the AKS cluster to be accessible through https and with a DNS. ([Read more](./platform/k8s-ingress/README.md)) 

## Deployment

The API itself has been deployed to the AKS cluster with two replicas.

The kubernetes manifests are located in [k8s folder](./k8s/)
- cruder.yaml - deployment, service, configmap and ingress for the API
- postgres.yaml - deployment, service and pvc for postgres
- db-migrate
  - base -> migrate.yaml - base for running goose migrations
  - overlays -> overlays for running either up (apply) or down (rollback) migrations 

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
  - applies kubernetes manifests to deploy API and postgres to AKS
- db-migrate-down
  - manually startable workflow in case a db rollback is required