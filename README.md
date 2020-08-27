# Backend API

This backend API is developed using Gin Web Framework: <https://github.com/gin-gonic/gin>

## Main technical points

- Web framework: gin-gonic - <https://github.com/gin-gonic/gin>
- DB ORM: gorm - <https://github.com/jinzhu/gorm>
- GraphQL server: gqlgen - <https://github.com/99designs/gqlgen>
- Auto build and restart if there are any update: <https://github.com/beego/bee>
  
---

## Quick start

### Install tools and libraries

- Go: https://golang.org/doc/install
- Bee (file watcher): `go get github.com/beego/bee`
- Run `go get` to install all dependencies in `go.mod` file
  
### Run backend

- Build project: run `go run ./scripts/build` to build project
- Auto detect modules and generate graphql: run `go run ./scripts/gqlgen`
- Run with dev mode: run `go run ./scripts/run_dev` (Note: create `config/.dev.env` file to overwrite the default environment file)
- Run with prod mode: run `go run ./scripts/run`
- Run with deploy mode: 
  - Run `go run ./scripts/deploy/service1` (Note: create `config/.dev.service.1.env` file to overwrite the default environment file)
  - Run `go run ./scripts/deploy/service2` (Note: create `config/.dev.service.2.env` file to overwrite the default environment file)

By default in development mode, project start with endpoint `localhost:7777`
We can open playground graphql with url `localhost:7777/graphql`

Example query:

```graphql - example module

query findExamples {
  	examples(input:{
    }) {
    	count,
    	list {
        id,
        code,
        description
      }
    }
}

query findExampleByCode {
  	examples(input:{
      code: "BN001",
    }) {
    	count,
    	list {
        id,
        code,
        description
      }
    }
}

query findExampleByDescription {
  	examples(input:{
      description: "Banana",
    }) {
    	count,
    	list {
        id,
        code,
        description
      }
    }
}

mutation createExampleWithInput{
  createExample(input:{
    code: "BN001",
    description: "Banana",
  }) {
    id,
    code,
    description
  }
}

```

```graphql - user module

query findUsers {
  	users(input:{
      email: ""
    }) {
    	count,
    	list {
        id,
        email,
        authenticationTokens {
          token,
          expiredOn
        }
      }
    }
}

query findUserByEmail {
  	users(input:{
      email: "test@gmail.com"
    }) {
    	count,
    	list {
        id,
        email,
        authenticationTokens {
          token,
          expiredOn
        }
      }
    }
}

mutation createUserWithInput{
  createUser(input:{
    email: "test@gmail.com",
    fullname: "Full name",
    nickname: "Nick name",
    password: "Just a password",
    avatarBase64: "",
    roleID: "ABC",
  }) {
    id,
    email
  }
}

```

```graphql - stock module

query findStocks {
  	stocks(input:{
    }) {
    	count,
    	list {
        id,
        code,
        description
      }
    }
}

query findStockByCode {
  	stocks(input:{
      code: "BN001",
    }) {
    	count,
    	list {
        id,
        code,
        description
      }
    }
}

query findStockByDescription {
  	stocks(input:{
      description: "Banana",
    }) {
    	count,
    	list {
        id,
        code,
        description
      }
    }
}

mutation createStockWithInput{
  createStock(input:{
    code: "BN001",
    description: "Banana",
  }) {
    id,
    code,
    description
  }
}

```

---

## Architect

Project using graphql as query language for APIs, and be layout as below

- config: contains config files for project
- server: play reponsibility for start server, config router
- utils: utility for project, manipulate environment variable, constant variable...
- modules: contain project modules
  - gql: hold all the related files for the graphql server, include `graph model, resolvers and schemas`
  - handlers: hold graphql server middleware for our server
  - orm: init database, contains orm models and migrations jobs for project, query context
- scripts: build scripts
  - build: build for production
  - gqlgen: generate graphql for all modules
  - run_dev: run server in development mode with code watcher using `bee run`
    - Note: create `config/.dev.env` file to overwrite the default environment file
  - run: run server in production mode
  - deploy: run multiple services locally with code watcher using `bee run`
    - service1: create `config/.dev.service.1.env` file to overwrite the default environment file for service 1
    - service2: create `config/.dev.service.2.env` file to overwrite the default environment file for service 2

### When you add a new module
- Remember to register handler in server/server.go->Run function

### Flags to run server
- (Optional)port: Provide a port to run the server on (Example: ./api -port=7788)
- (Optional)envfile: Provide an environment file to overload the server configuration (Example: ./api -envfile="config/.dev.env")

### Coding standard
- Rule number 1: Please observe the coding style and code template carefully and try to keep code as clean as possible, don't make your own style unless has been approved.
- Rule number 2: Write comments as the format tool suggest, don't leave warning on any file.
- Rule number 3: Keep your code GIT friendly by not writing a very long line of code, try to use line break.

### Other

- logger: logger utility for project
- scripts: contains script for run, build, deploy...

---

## Deploy with docker

Prerequisites: Install Docker from Docker home page

### 1. Deploy separate container

#### Postgres container

The first of all, we need initiate a database server, here is Postgres.
To set up a postgres container, folow instruction steps:

- Pull `postgres`[https://hub.docker.com/_/postgres] image from docker hub

> docker pull postgres

Check postgres image is pull to local successful

> docker images

Check postgres with latest tag is exist in list images

- Run postgres image to initiate a new container instance

> docker run --name my-postgres -p 2345:5432 -e POSTGRES_PASSWORD=admin postgres

Initiate a new container with name `my-postgres`, and export port 2345 to docker host map with 5432 of postgres server in container. After running above command, container with name `my-postgres` will be run.

- Access to postgres container to create master database

> docker exec -it my-postgres bash

Using psql to create db

> psql -U postgres

Create database with name `ctbwebmaster`

> create database ctbwebmaster;

Make sure datbase is created successfully by command `\l` of psql

> \l

You can see `ctbwebmaster` in table list

#### Backend API container

- Config postgres database server in environment file

Update following env variable in `.env` file config

```env

DB_SERVER_HOST=0.0.0.0
DB_SERVER_PORT=2345
DB_SERVER_USER=postgres
DB_SERVER_PASS=admin
DB_NAME=ctbwebmaster

```

- Build backend api image.

> docker build -t backend:1.0.0 .

Build `backend` image with tag `1.0.0`

- Start backend container

> docker run -it backend:1.0.0 bash

Now you can work with backend as normally

### 2. Using docker compose

In case you want to start all service for back end environment that you can develop front end part or for testing purpose, you can run backend with docker-compose

- Config `.env` file for docker compose

```env

# Used by pgadmin service
PGADMIN_DEFAULT_EMAIL=ctb@cookingthebooks.com.au
PGADMIN_DEFAULT_PASSWORD=admin

#Postgres config for docker
DB_SERVER_HOST=db
DB_SERVER_PORT=5432
DB_SERVER_USER=postgres
DB_SERVER_PASS=admin
DB_NAME=ctbwebmaster
```

- Start docker-compose

> docker-compose up --build

Docker-compose reads step by step in `docker-compose.yml` file and start container for postgres, pgadmin as well as api_module back-end service.

- Using api server

Go to your browser and open

> http://localhost:7777/auth/graphql

to register new account and working with backend as normally

- Using pdadmin

Go to your browser and open

> http://localhost:5050

Login with account ctb@cookingthebooks.com.au/admin as config from `.env` file

Right Click on the Servers to create a new server. Choose Create then Server

Fill in server information to connect our postgres database

```postgres
- Name: anything name you want
- Host: db
- Port: 5432
- Maintainer database: postgres
- Username: postgres
- Password: admin
```

Click Save to connect to postgres database

- Stops containers and removes containers, networks, volumes, and images created by up

> docker-compose down --remove-orphans --volumes

---

## Deploy with local kubernetes

### Install required app/tools

- [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/): run Kubernetes locally
- [kubectl (Kubernetes command-line tool)](https://kubernetes.io/docs/tasks/tools/install-kubectl/) : allows you to run commands against Kubernetes clusters
- [Docker driver](https://minikube.sigs.k8s.io/docs/drivers/docker/) : allows you to install Kubernetes into an existing Docker install

### Deployment for Postgres database

Our database (postgresql) will need a PersistentVolume(PV) and a PersistentVolumeClaim(PVC). The define for them in `postgres-db-pv.yaml` and `postgres-db-pvc.yaml`


Using environment `.env` for kubernetes:

```env
DB_SERVER_HOST=fullstack-postgres      # service name
DB_SERVER_PORT=5432
DB_DRIVER=postgres
DB_SERVER_USER=postgres
DB_SERVER_PASS=admin
DB_NAME=ctbwebmaster
```

- Start minikube

> minikube start

- Set environment variables for Postgres

> kubectl create -f postgres-secret.yaml

After creating the secret, we confirmed that the creation process was successful by running

> kubectl get secrets

> kubectl describe secrets  postgres-secret

Observe the “DATA” has 10 items, exactly the number of items we have in our secret file.

- Deploy postgres

Apply the commands we have in each `.yaml` file. These should be run one after the other:

``` command
kubectl apply -f postgres-db-pv.yaml
kubectl apply -f postgres-db-pvc.yaml
kubectl apply -f postgres-db-deployment.yaml
kubectl apply -f postgres-db-service.yaml
```

A pod is created. To view the pod run:

> kubectl get pods

A postgres image is pulled from docker, then container will be created so this might take some time

You can check some status of pod using following command

```command
kubectl describe pod <pod_name>
kubectl logs <pod_name
```

Now the service for postgres created successfully!

### Deploy API to Kubernetes

#### Push API to docker hub

- Create docker hub account to publish image. Then create a repository for storing images
Assume that we have dockerhub account `ctb` and repository `go-api`

- Build image with Dockerfile

``` command
docker build -t <app-name> .

<app-name>: naming for image up to you.
```

For this we naming `go-api`

- Add tag for image

```
docker tag go-api ctb/go-api:1.0.0
```

- Push to docker hub

> docker push ctb/go-api:1.0.0

#### Deployment API

- Apply deployment script in `.yaml` file

``` command
kubectl apply -f app-postgres-deployment.yaml
kubectl apply -f app-postgres-service.yaml
```

- Get pods to check deployment successfully

> kubectl get pods

You should see the pod `go-api` running with 1 replica. We can scale number of replicas after that.
Cool! Now we can test app.

- Get the list of services

> kubectl get services

We are going to get the URL that was exposed to us for the `go-api`:

> minikube service go-api --url

We can access this url and do everything with go-api as normally
