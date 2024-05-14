# Setup

Basic setup
- Install Make
- Install Docker
- Install Golang
- $>make build
- $>make dev
- $>API is at http://localhost:8080/order-service/apidocs


Make is used to build/lint/test see makefile for all options the most useful ones are:

To set up the environment you may need to sign in to docker to be able to pull docker images 
```bash
$ make build
```

```bash
# to run linter
$ make lint
```

```bash
# to run all tests
$ make test
```

```bash
# could use "docker-compose up" instead of
$ make dev 
```

```bash
# will remove the volumes and remake the environment
$ make cleandev
```

Once the `make dev` command completes the order service will be hosted at http://localhost:8080/order-service/apidocs

You can then start using it. Here is an example order that can be passed to the create API
```json
{
  "customerId": "d3ba8849-ad3c-4a27-bb7d-2f3731841fbb",
  "orderId": "ea3d06cf-c4bf-4be7-92de-430b952a1111",
  "orderItems": [
    {
      "id": "eb6a5bca-36e7-471b-b3a9-e0c588b70111",
      "orderId": "ea3d06cf-c4bf-4be7-92de-430b952a1111",
      "price": 222,
      "productId": "550e8400-e29b-41d4-a716-446655440000",
      "quantity": 1
    }
  ]
}
```

To start the payment processing service you can run it manually from the project dir with GO or use the output executable from `make build` and run it in alpine linux. I left it separate so that you can inspect the queue before running the payment service. 

You can inspect rabbitmq at http://localhost:15672/#/queues with the credentials being `user` and `password` respectively.

## Things that need completion or would make the service(s) better
As this is a demo running on a localhost there are a few things missing
- more tests -- some are there, however, there are some missing. Generally how I test is demonstrated.
- strong passwords -- not in plaintext and not submitted to the repo
- needs proper password management -- e.g. hashicorp vault, kubernetes env vars
- database tables for customers and payment processing
- table indexing for faster SQL queries
- the service could benefit from tracing, metrics, log aggregation etc. -- e.g. otel, prometheus, loki, grafana
- needs CI/CD -- github, gitlab
- needs deployment configuration to an orchestrtor -- e.g. kubernetes
- needs a better README.md with diagrams -- e.g. sequence diagram with mermaid
- could always benefit from code review
- Passing the full 'order' object around with the queue isn't ideal and contracts should be made to specify what inputs and outputs are expected. 

Things to note:
The order service is written using the hexagonal architecture pattern
The payment service is intentionally bare-bones to show the strong contrast