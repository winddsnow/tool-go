# Tool Go — AGENTS.md

Monorepo: GoFrame v2 backend (PostgreSQL) + Vue 3 + TypeScript + Element Plus frontend.

## Project structure

| Directory | Purpose |
|-----------|---------|
| `api/v1/` | Go request/response structs (OpenAPI tagged, `g.Meta` declares route) |
| `internal/cmd/cmd.go` | Backend entrypoint — routes, middleware binding |
| `internal/controller/` | HTTP handlers (thin, delegates to `service`) |
| `internal/logic/` | Business logic (registers via `init()` → `service.RegisterXxx`) |
| `internal/middleware/` | Auth (JWT), Permission (role check), CORS |
| `internal/service/` | Interface + singleton registry pattern |
| `internal/dao/` | GoFrame DAO — auto-generated via `gf gen dao` |
| `internal/model/entity/` | DB entity structs — auto-generated |
| `internal/library/password/` | MD5 + Salt password hashing |
| `internal/library/jwt/` | JWT create/parse using `golang-jwt/jwt/v5` |
| `web/src/api/` | Axios request wrappers (proxy `/api` → `:8000`) |
| `web/src/router/index.ts` | Frontend routes + permission guard |
| `web/src/store/modules/user.ts` | Pinia store — token, roles, `hasRole`/`hasAnyRole` |
| `manifest/sql/init.sql` | DB schema + seed data (PostgreSQL) |
| `manifest/config/config.yaml` | App config (DB, JWT, server) |

## Commands

```bash
# Backend (requires PostgreSQL — see init.sql + config.yaml)
go mod tidy
go run main.go              # dev server :8000, swagger at /swagger

# Frontend
cd web
npm install
npm run dev                 # :3000, proxies /api → :8000
npm run build               # vue-tsc --noEmit + vite build
npm run lint                # eslint --fix
npm run type-check          # vue-tsc --noEmit

# Go codegen (after DB schema change)
# go install github.com/gogf/gf/cmd/gf/v2@latest
gf gen dao

# Password hash utility (for seed data)
go run hack/gen_password.go

# Env switching: export GF_GENV=prod  loads config.prod.yaml
```

## Auth & permission model

- **Backend**: `Auth` middleware parses JWT → stores `userId`/`roles` in context. `Permission(roles ...string)` middleware factory gates sub-routes via `HasAnyRole`. All role-gated routes in `cmd.go`.
- **Frontend**: Router guard checks `userStore.hasAnyRole(route.meta.roles)`. Per-view `v-if="userStore.hasAnyRole([...])"` hides elements.
- Both backend `Permission()` calls and frontend `meta.roles` must stay in sync — mismatches cause invisible-but-accessible routes.

## Key gotchas

- **Status=0 filter**: `UserListReq.Status` is `*int` (not `uint`) — Go zero-value `0` prevents filtering disabled users. Always use pointer type for optional numeric query filters.
- **GoFrame query order**: After `Count()`, the model is consumed. Chain `.Page().OrderDesc().Scan()` in a fresh chain.
- **Password hashing**: Uses MD5 + Salt (`internal/library/password/`), not bcrypt. Use `go run hack/gen_password.go` to generate seed-compatible hashes.
- **PostgreSQL**: Driver imported via blank import `_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"` in `main.go`. Not MySQL.
- **npm run build** runs `vue-tsc --noEmit` first — type errors block the build.
- **Default account**: `admin`/`admin123` with `super_admin` role (check init.sql for actual seed data).
- **JWT default**: `config.yaml` has `secret: "dev-jwt-secret-key-123456"`. Change in production.
- **Swagger**: Enabled at `/swagger` and `/api.json` in dev mode.
- **Docker**: `manifest/docker/Dockerfile` — override DB via env vars.
