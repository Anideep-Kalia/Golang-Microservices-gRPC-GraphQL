# Go-GraphQL-Microservice

This project is a Go-based microservice architecture using gRPC, GraphQL, and various other technologies like PostgreSQL, ElasticSearch, and Protocol Buffers.

---

## Table of Contents

- [Technologies Used](#technologies-used)
- [Setup Instructions](#setup-instructions)
  - [GraphQL](#graphql)
  - [gRPC](#grpc)
  - [Protocol Buffers](#protocol-buffers)
- [Microservices Overview](#microservices-overview)
  - [Account](#account)
  - [Catalog](#catalog)
  - [Order](#order)
- [Database Operations](#database-operations)
- [Example Queries and Mutations](#example-queries-and-mutations)
- [Syntax Guide](#syntax-guide)

---

## Technologies Used

- **GraphQL** (gqlgen)
- **gRPC** with Protocol Buffers
- **ElasticSearch** for catalog microservice
- **PostgreSQL** for account and order microservices
- **HTTP/2** for faster transport
- **Docker Compose** for project orchestration

---

## Setup Instructions

### GraphQL
1. Install gqlgen:
```bash
    go run github.com/99designs/gqlgen generate
```
2. Define your schema in .graphql files and configure gqlgen.yaml.
3. Generate resolver interfaces by running the above command.


### gRPC
1. Install required tools:
```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
2. Create a pb directory for Protocol Buffers:  mkdir pb
3. Compile .proto files:
``` bash
    protoc ./account.proto --go_out=./pb --go-grpc_out=./pb
```
4. Update imports to pb in both client and server files.


### Protocol Buffers
- Used to define schemas for gRPC.
- Provides compact, efficient, and strongly-typed serialization.

---

## Microservices Overview
### Account
- Database: PostgreSQL
- Layers:
    - Repository: Interacts with the database for SQL queries.
    - Service: Contains business logic, triggers repository functions.
    - Server: Handles gRPC requests and communicates with the service layer.
### Catalog
- Database: ElasticSearch
- Features:
    - Full-text search and data analytics.
    - High-speed search capabilities.
    - Horizontal scaling with node sniffing.
### Order
- Database: PostgreSQL
- Features:
    - Supports multiple products in an order.
    - Implements transactional operations using BeginTx

---
## Database Operations
### PostgreSQL
- QueryRowContext: Fetch a single row.
- ExecContext: Execute an SQL command.
- QueryContext: Execute multiple-row queries.
- PrepareContext: Prepare reusable SQL commands.
- CopyIn: Insert multiple rows in one query.
---
## Example Queries and Mutations
### GraphQL Queries
#### Fetch all accounts:
``` graphql
query {
  accounts {
    id
    name
  }
}
```
#### Fetch products with pagination:
``` graphql
query {
  products(pagination: { skip: 0, take: 5 }, query: "?") {
    id
    name
    price
  }
}
```
#### Fetch orders for an account:
``` graphql
query {
  accounts(id: "?") {
    name
    orders {
      totalPrice
    }
  }
}
```
### GraphQL Mutations
#### Create a new account:
``` graphql
mutation {
  createAccount(account: { name: "New Account" }) {
    id
    name
  }
}
```
#### Create a new product:
``` graphql
mutation {
  createProduct(product: { name: "New Product", description: "A new product", price: 199 }) {
    id
    name
    price
  }
}

```
#### Create a new order:
``` graphql
mutation {
  createOrder(order: { accountId: "??", products: [{ id: "?", quantity: ? }] }) {
    id
    totalPrice
    products {
      name
      quantity
    }
  }
}

```
---
## Syntax Guide
Retry with exponential backoff:
``` graphql
retry.ForeverSleep(2 * time.Second, func(attempt int) (err error) {
    if attempt >= maxRetries {
        fmt.Println("Max retries reached. Stopping.")
        return nil
    }

    r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
    if err != nil {
        log.Printf("Attempt %d failed: %v", attempt+1, err)
        return err
    }

    fmt.Println("Successfully connected!")
    return nil
})

```
---

## Running the Project
To start all microservices using Docker Compose:
``` graphql
docker-compose up --build
```
This command builds and runs all services in the project.

---

## Features of gRPC
- HTTP/2 for faster transport.
- Streaming: Supports unary, server, client, and bidirectional streaming.
- Authentication: Supports SSL/TLS and token-based mechanisms.
- Protobuf Serialization: Compact and efficient binary format.
- Interoperability: Compatible with multiple programming languages.

---

## Features of ElasticSearch
Full-text search and analytics.
Sniffing for node discovery and scaling.
High-speed search and indexing.

---


