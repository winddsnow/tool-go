# Tool Go — AGENTS.md

Monorepo: GoFrame v2 backend (PostgreSQL) + Vue 3 + TypeScript + Element Plus frontend.

## Project structure

| Directory | Purpose |
|-----------|---------|
| `api/v1/` | Go request/response structs + route declarations (`g.Meta` tag) |
| `internal/cmd/cmd.go` | Backend entrypoint — routes, middleware binding |
| `internal/controller/` | HTTP handlers (thin, delegates to `service`) |
| `internal/logic/` | Business logic (registers via `init()` → `service.RegisterXxx`) |
| `internal/middleware/` | Auth (JWT), Permission (role check), CORS |
| `internal/service/` | Interface + singleton registry (sync.RWMutex pattern) |
| `internal/dao/` | Hand-written GoFrame DAO wrappers with Column constants |
| `internal/model/entity/` | DB entity structs |
| `internal/model/do/` | Data Operation structs (used in inserts/updates) |
| `internal/library/password/` | MD5 + Salt password hashing |
| `internal/library/jwt/` | JWT create/parse using `golang-jwt/jwt/v5` |
| `manifest/config/config.yaml` | App config (DB, JWT, server) |
| `manifest/sql/init.sql` | DB schema + seed data (PostgreSQL) |
| `web/src/api/` | Axios request wrappers (baseURL `/api/v1`, auto-attaches JWT) |
| `web/src/router/index.ts` | Frontend routes + permission guard |
| `web/src/store/modules/user.ts` | Pinia store — token, roles, `hasRole`/`hasAnyRole` |

## Commands

```bash
# Backend (requires PostgreSQL — see manifest/sql/init.sql + manifest/config/config.yaml)
go mod tidy
go run main.go              # dev server :8000, swagger at /swagger

# Frontend (Vite proxies /api → :8000)
cd web
npm install
npm run dev                 # :3000
npm run build               # vue-tsc --noEmit + vite build
npm run lint                # eslint --fix
npm run type-check          # vue-tsc --noEmit

# Password hash utility (for seed-compatible hashes)
go run hack/gen_password.go

# Env switching: export GF_GENV=prod  loads config.prod.yaml
```

## Auth & permission model

- **Backend**: `Auth` middleware parses JWT from `Authorization: Bearer <token>` → stores `userId`/`username`/`roles` in context. `Permission(roles ...string)` middleware factory gates sub-routes via `HasAnyRole`. All role-gated routes in `cmd.go`.
- **Frontend**: Router guard checks `userStore.hasAnyRole(route.meta.roles)`. Per-view `v-if="userStore.hasAnyRole([...])"` hides elements.
- Backend `Permission()` calls and frontend `meta.roles` must stay in sync.

## Key gotchas

- **Seed account**: `walter`/`walter` with `super_admin` role (from `manifest/sql/init.sql`). Not `admin`/`admin123`.
- **Status=0 filter**: `UserListReq.Status` is `*int` (pointer) — Go zero-value `0` prevents filtering disabled users with plain `uint`. Always use pointer type for optional numeric query filters.
- **GoFrame query order**: After `Count()`, the model is consumed. Chain `.Page().OrderDesc().Scan()` in a fresh chain.
- **Password hashing**: Uses MD5 + Salt (`internal/library/password/`), not bcrypt. Use `go run hack/gen_password.go` to generate seed-compatible hashes.
- **PostgreSQL driver**: Blank import `_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"` in `main.go`.
- **Frontend auto-imports**: `unplugin-auto-import` + `unplugin-vue-components` with `ElementPlusResolver` — components like `ElButton`, `ElMessage` are globally available without explicit imports.
- **Response unwrapping**: Frontend Axios interceptor strips the `{code, message, data}` wrapper — API functions receive `data` directly.
- **No tests**: Neither backend Go tests nor frontend tests exist in this repo.
- **Frontend `VITE_API_BASE_URL`**: Development loads from `web/.env.development` (`http://127.0.0.1:8000/api/v1`); production from `web/.env.production` (`/api/v1`). Vite proxy also forwards `/api` → `:8000` for same-origin dev.
- **Swagger**: Enabled at `/swagger` and `/api.json` in dev mode.
- **JWT default**: `config.yaml` has `secret: "dev-jwt-secret-key-123456"`. Change in production.
- **Docker**: `manifest/docker/Dockerfile` — override DB via env vars.
- **npm run build** runs `vue-tsc` first — type errors block the build.
- **Service registration**: Each `logic/*.go` has `init()` → `service.RegisterXxx(NewXxx())`. Service layer uses `sync.RWMutex` for singleton access.
