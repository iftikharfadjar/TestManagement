# 📋 Restructuring Proposal: Layered → Feature-Based

## Status: ✅ COMPLETED

### Final Structure

```
cmd/server/main.go

be/auth/                    # Feature: Authentication
├── domain/user.go          # User, AuthRepository, AuthUseCase
├── usecase/service.go     # Business logic
└── delivery/
    ├── rest/            # HTTP handlers
    └── graphql/       # GraphQL resolvers

shared/                     # Cross-cutting
├── config/config.go      # config loading
├── jwt/jwt.go         # JWT utilities
├── middleware/auth.go  # Fiber JWT middleware
├── db/db.go          # PocketBase init
└── adapter/          # DB Adapters
    ├── pocketbase/
    └── sqlite_adapter/
        └── sqlc/

graph/                     # GraphQL (unchanged)
sql/                       # migrations, schema
```

### Migration Complete
| From | To |
|------|-----|
| `internal/domain/auth.go` | `be/auth/domain/user.go` |
| `internal/usecase/auth/service.go` | `be/auth/usecase/service.go` |
| `internal/rest/*` | `be/auth/delivery/rest/` |
| `internal/graphql/*` | `be/auth/delivery/graphql/` |
| `internal/repository/pocketbase/*` | `shared/adapter/pocketbase/` |
| `internal/repository/sql_adapter/*` | `shared/adapter/sqlite_adapter/` |
| `pkg/config/*` | `shared/config/` |
| `pkg/jwt/*` | `shared/jwt/` |
| `internal/auth/*` | `shared/middleware/` |
| `internal/db/*` | `shared/db/` |

### Verification
- `go build ./...` ✅
- `go test ./...` ✅