# Tool Go — AGENTS.md

## Project structure

Monorepo: GoFrame v2 backend + Vue 3 + TypeScript + Element Plus frontend.

| Directory | Purpose |
|-----------|---------|
| `api/v1/` | Go request/response structs (OpenAPI tagged) |
| `internal/cmd/cmd.go` | **Backend entrypoint** — routes, middleware binding |
| `internal/controller/` | HTTP handlers |
| `internal/logic/` | Business logic |
| `internal/middleware/` | Auth (JWT), Permission (role check), CORS |
| `internal/dao/` | GoFrame DAO — auto-generated via `gf gen dao` |
| `internal/model/entity/` | DB entity structs — auto-generated |
| `web/src/api/` | Axios request wrappers |
| `web/src/router/index.ts` | **Frontend routes + permission guard** |
| `web/src/store/modules/user.ts` | Pinia store — token, roles, `hasRole`/`hasAnyRole` |
| `web/src/views/login/index.vue` | Login page — sets token + roles in store |
| `manifest/sql/init.sql` | DB schema + seed data |
| `manifest/config/config.yaml` | App config (DB, JWT, server) |

## Commands

```bash
# Backend
go mod tidy
go run main.go                        # dev server :8000

# Frontend
cd web
npm install
npm run dev                           # dev server :3000, proxies /api -> :8000
npm run build                         # production build -> web/dist/
npm run lint                          # eslint fix
npm run type-check                    # vue-tsc noEmit

# One-click dev startup (Win)
start.bat

# Go codegen (after DB schema change)
# go install github.com/gogf/gf/cmd/gf/v2@latest
gf gen dao

# Env switching
# Windows: $env:GF_GENV="prod"
# Linux:   export GF_GENV=prod
# Loads manifest/config/config.prod.yaml
```

## Auth & permission model

- **Backend**: `internal/middleware/auth.go` parses JWT -> stores `userId`/`roles` in context. `Permission("super_admin", "admin")` in `cmd.go` gates routes. `HasRole`/`HasAnyRole` check context roles.
- **Frontend**: Router guard in `web/src/router/index.ts` checks `userStore.hasAnyRole(route.meta.roles)`. Per-view `v-if="userStore.hasAnyRole([...])"` hides elements.
- Both frontend sidebar visibility AND create-user button use `hasAnyRole(['super_admin', 'admin'])`.

## Known pitfalls

- **Status=0 query param**: The backend `UserListReq.Status` field uses `*int` (not `uint`) because Go zero-value `0` prevents filtering disabled users. Always keep pointer type for optional/numeric query filters.
- **Permission middleware**: Is variadic (`middleware.Permission(roles ...string)`); pass multiple roles to allow any match.
- **Frontend role checks**: All public-facing checks (`v-if`, `hasAnyRole`) must stay in sync with backend `Permission()` calls — mismatches cause invisible-but-accessible routes.
- **JWT default**: `config.yaml` has `secret: "dev-jwt-secret-key-123456"`. Change in production.
- **README default account**: `admin`/`admin123` with `super_admin` role — ensure init.sql creates this user with the correct role assignment.
- **GoFrame query order**: Chain `.Page().OrderDesc().Scan()` after `Count()` — `Count` consumes the model, re-query needed.
