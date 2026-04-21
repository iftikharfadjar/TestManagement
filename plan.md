# Backend Boilerplate Implementation Plan

## 1. Overview
A Go-based backend boilerplate utilizing **Fiber v3** as the primary web framework, integrating a highly decoupled **Clean Architecture / Ports & Adapters** footprint. The architecture gracefully supports multiple database adapters dynamically at run-time (**PocketBase**, **SQLite**, **PostgreSQL**) and exposes API connectivity via both **REST** and **GraphQL** endpoints.

## 2. Architecture & Design Principles
- **Clean Architecture:** The system strictly separates business logic from technical frameworks. The pure entities (`services/auth/domain/`) and application business logic (`services/auth/usecase/`) have absolutely no knowledge of Fiber, PocketBase, SQL, or GraphQL.
- **Dynamic Database Swapping:** By enforcing the `domain.AuthRepository` interface, the boilerplate supports reading a `config.json` file at boot to gracefully inject either PocketBase, a raw `database/sql` SQLite instance, or a PostgreSQL instance.
- **sqlc Integration:** Raw SQL queries are written natively in `sql/` and strongly typed via `sqlc generate`. The auto-generated code serves as the core mapping foundation for both the SQLite and PostgreSQL backend adapters.
- **Golang Migrations:** The boilerplate launches automated schema verifications targeting `sql/migrations` upon start-up, shielding structural inconsistencies across environments.
- **Custom Authentication (`jwx/v3`):** To avoid locking into PocketBase's proprietary token mechanisms, the boilerplate employs a standardized enterprise `jwx` implementation.

## 3. Project Structure
```text
cmd/server/main.go                 # App wireup / Config Loading

services/                        # Feature-based modules
├── auth/
│   ├── domain/                  # Pure Go structs (User) and abstraction interfaces
│   ├── usecase/                 # Pure business orchestrator (DB agnostic)
│   └── delivery/
│       ├── rest/               # Fiber delivery adapter & Endpoints
│       └── graphql/            # gqlgen delivery adapter

shared/                          # Cross-cutting concerns
├── config/                      # Dynamic config.json + os.Getenv loader
├── jwt/                         # lestrrat-go/jwx/v3 native token signer and verifier
├── middleware/                 # Fiber middleware parsing HTTP Bearer tokens
├── db/                         # PocketBase initialization core
└── adapter/                    # DB Adapters
    ├── pocketbase/            # PocketBase SDK adapter
    └── sqlite_adapter/
        └── sqlc/              # Type-safe auto-generated SQL mappings (DO NOT EDIT)

graph/                          # GraphQL schema + generated code
sql/                           # migrations, schema.sql, query.sql
go.mod
sqlc.yaml                       # SQL code-generation configurations
config.json                     # Environment definitions
```

## 4. Key Components

### A. Configuration Management
- Loaded dynamically using `shared/config`. It prioritizes `config.json` natively before falling over to standard Environment Variables.

### B. Database Integrations
**SQL-based Adapters:**
- Rely absolutely on `database/sql`, executing raw queries managed tightly via `sqlc`.
- `main.go` calls into `golang-migrate/migrate` instantly verifying local schemas before allocating ports.

**PocketBase Integration:**
- Optionally boots into a headless mode proxied flawlessly through Fiber using `app.All("/_/*")` to preserve Administration UI benefits securely.

### C. Web Interfaces
**Fiber REST API Endpoints:**
- Standard payloads are unmarshalled inside `services/auth/delivery/rest` controllers mapping gracefully to UseCases.
- Endpoints: `POST /api/v1/signup`, `POST /api/v1/login`, `POST /api/v1/logout`

**GraphQL:**
- Leverages `github.com/99designs/gqlgen`. Resolvers natively utilize `domain.AuthUseCase` independently of Fiber architectures.

### D. Security & Lifecycle
- Handled primarily by `shared/jwt`.
- Intercepted by Fiber's `middleware.Middleware(...)` parsing `Authorization` headers. If validated securely by the `jwx` engine, the inner `*domain.User` entity resolves successfully and is tightly bound to Context Locals.

## 5. Current Implementation State
- **[COMPLETED]** Defined pure structural Domain logic.
- **[COMPLETED]** Established generic Dependency Injection wire-ups within `cmd/server/main.go`.
- **[COMPLETED]** Interfaced PocketBase SDK safely via independent Repositories.
- **[COMPLETED]** Interfaced full PostgreSQL + SQLite run-times seamlessly utilizing `sqlc`.
- **[COMPLETED]** Assembled Custom Auth using enterprise `jwx/v3` mechanics.
- **[COMPLETED]** Automated run-time schema synchronization via `golang-migrate`.
- **[COMPLETED]** Validated logic through standalone decoupled Unit Tests.