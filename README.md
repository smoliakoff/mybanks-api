# MyBanks API

**GraphQL Bank Directory Backend (Go + Ent + Relay)**

A production-style GraphQL API for a bank directory, built with **Go**, **Ent ORM**, and **Relay-compliant pagination**.

This project demonstrates:

* Schema-driven development using **Ent** as the single source of truth
* Auto-generated GraphQL schema via **entgql**
* Relay-style cursor pagination
* Strongly typed filtering (`WhereInput`)
* Translation support per entity (`BankTranslation`)
* Clean and reproducible generation pipeline (`make generate`)
* Strict separation between generated and manual code

---

## ✨ Features

* 🔁 Relay-compliant pagination (`Connection`, `Edge`, `Cursor`)
* 🔎 Strongly typed filtering (`BankWhereInput`, `OfferWhereInput`, etc.)
* 🌍 Localized bank translations (`translation(locale: String!)`)
* ⚙️ Fully reproducible code generation
* 🧱 Ent as the single source of truth
* 🧩 Minimal and explicit resolver layer
* 📜 Optional OpenAPI spec generation

---

## 🏗 Architecture

The project follows a layered generation model:

### 1️⃣ Ent Schema (Source of Truth)

All entities are defined in:

```
ent/schema/*
```

From this, Ent generates:

* Database client
* Query builders
* Pagination helpers
* GraphQL schema via `entgql`

---

### 2️⃣ GraphQL Schema (Auto-generated)

`entgql` generates:

```
graph/schema.graphqls
```

This includes:

* Object types
* Relay connections
* Edge types
* `WhereInput` filters
* Node interface
* Optional CRUD mutations

---

### 3️⃣ gqlgen Runtime Layer

`gqlgen` generates:

```
graph/generated.go
graph/model/models_gen.go
```

Manual resolvers live in:

```
graph/resolver.go
```

Generated files are never edited manually.

---

## 🛠 Tech Stack

* Go 1.26+
* Ent (ORM)
* entgql (GraphQL schema generator)
* gqlgen (GraphQL runtime)
* PostgreSQL (default database)

---

## 🚀 Getting Started

### Requirements

* Go 1.26+
* PostgreSQL
* make
* (optional) `yq` for OpenAPI enrichment

---

### Clone & Install

```bash
git clone https://github.com/smoliakoff/mybanks-api.git
cd mybanks-api
go mod download
```

---


## ⚙️ Configuration

The application is fully environment-driven and does not use hardcoded credentials.

All configuration is loaded from environment variables.

### Required Variables

| Variable       | Description                  | Example                                                     |
| -------------- | ---------------------------- | ----------------------------------------------------------- |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://app:app@localhost:5432/mybanks?sslmode=disable` |

### Optional Variables

| Variable  | Description             | Default       |
| --------- | ----------------------- | ------------- |
| `PORT`    | HTTP server port        | `8080`        |
| `APP_ENV` | Application environment | `development` |

---

## 🧪 Local Development (.env)

For local development, create a `.env` file in the project root:

```env
DATABASE_URL=postgres://app:app@localhost:5432/mybanks?sslmode=disable
PORT=8080
APP_ENV=development
```

The application can load this file automatically (e.g., using `godotenv`).

---

## 🏗 Configuration Structure

Configuration is centralized in a dedicated package:

```
internal/config/config.go
```

Example structure:

```go
type Config struct {
    DatabaseURL string
    Port        string
    Environment string
}
```

All environment variables are loaded once at application startup and injected into the application.

---

## 🚀 Production Deployment

In production environments (Docker, Kubernetes, CI/CD), environment variables should be provided by the runtime environment:

```bash
export DATABASE_URL="postgres://user:pass@db:5432/mybanks?sslmode=require"
export PORT=8080
export APP_ENV=production
```

No configuration changes are required in code.

---

## 🔧 Code Generation

Generate Ent client + GraphQL schema + gqlgen runtime:

```bash
make generate
```

If OpenAPI enrichment is enabled:

```bash
make openapi
```

The generation pipeline is deterministic and safe to run at any time.

---

## ▶ Run Server

```bash
go run .
```

Example endpoint:

```
http://localhost:8080/query
```

If GraphQL Playground is enabled:

```
http://localhost:8080/
```

---

## 📚 GraphQL Examples

### Relay Pagination

```graphql
query Banks($first: Int, $after: Cursor) {
  banks(first: $first, after: $after) {
    totalCount
    pageInfo {
      hasNextPage
      endCursor
    }
    edges {
      cursor
      node {
        id
        name
        country
      }
    }
  }
}
```

---

### Filtering

```graphql
query {
  banks(where: { country: "Georgia" }) {
    edges {
      node {
        id
        name
      }
    }
  }
}
```

---

### Translation by Locale

```graphql
query {
  banks(first: 10) {
    edges {
      node {
        id
        name
        translation(locale: "en") {
          name
          description
        }
      }
    }
  }
}
```

If a translation does not exist, `translation` returns `null`.

---

## 📂 Project Structure

```
.
├─ ent/
│  ├─ schema/              # Entity definitions (source of truth)
│  ├─ generated client
│  └─ openapi.json
│
├─ graph/
│  ├─ schema.graphqls         # Auto-generated GraphQL schema
│  ├─ generated.go         # gqlgen runtime (generated)
│  ├─ resolver.go          # Manual resolvers
│  └─ model/
│     └─ models_gen.go
│
├─ main.go
├─ Makefile
└─ go.mod
```

---

## 🧠 Design Principles

* Ent is the single source of truth
* No manual edits in generated files
* Relay-style pagination only
* Explicit filtering via typed inputs
* Minimal and transparent resolver layer
* Production-style reproducible generation

---

## 🔄 Generation Philosophy

The project intentionally separates:

* Schema definition
* Generated code
* Manual logic
* Runtime wiring

You can safely run:

```bash
make generate
```

at any time without breaking manual logic.

---

## 📌 Purpose

This repository serves as:

* A reference architecture for Ent + GraphQL + Relay
* A portfolio-grade backend example
* A foundation for building a production-ready directory service

---


---

## 🐳 Docker Support

The application can be containerized using Docker.

### Example `Dockerfile`

```dockerfile
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app ./app

EXPOSE 8080

CMD ["./app"]
```

---

### Example `docker-compose.yml`

```yaml
version: "3.9"

services:
  db:
    image: postgres:16
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: app
      POSTGRES_DB: mybanks
    ports:
      - "5432:5432"

  api:
    build: .
    depends_on:
      - db
    environment:
      DATABASE_URL: postgres://app:app@db:5432/mybanks?sslmode=disable
      PORT: 8080
      APP_ENV: development
    ports:
      - "8080:8080"
```

Run locally with:

```bash
docker-compose up --build
```

Build with BUILD_SHA tag
```bash
 BUILD_SHA=$(git rev-parse --short HEAD) docker compose up -d --build --force-recreate api
```


## 📜 License

MIT
