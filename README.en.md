# Tool Go — Developer Toolbox

A development and admin toolbox built with GoFrame v2 + Vue 3 + TypeScript + Element Plus.

## Tech Stack

### Backend

- **GoFrame v2** — High-performance Go web framework
- **PostgreSQL** — Relational database
- **JWT (golang-jwt/jwt/v5)** — Authentication (Access Token + Refresh Token)
- **RESTful API** — Standard API style
- **Swagger** — API docs (dev mode)

### Frontend

- **Vue 3** — Progressive JavaScript framework
- **TypeScript** — Type safety
- **Vite** — Fast build tool
- **Element Plus** — UI component library
- **Pinia** — State management
- **Vue Router** — Dynamic routes (menu-driven)

## Features

### Toolbox (12 tools, client-side)

| Category | Tools |
|----------|-------|
| Text | JSON Formatter, Text Diff, Regex Tester, Case Converter |
| Encoding | Base64 Encoder/Decoder, Hash (MD5/SHA1/SHA256) |
| Generation | Password Generator, Mock Data Generator (9 types), UUID Generator, QR Code Generator |
| Conversion | Timestamp Converter (16 timezones) |

### Admin Panel

- **Dashboard** — System stats (users, roles, visits)
- **User Management** — CRUD, role assignment, paginated search
- **Role Management** — CRUD, menu assignment, permission assignment, paginated search
- **Menu Management** — Dynamic menu tree CRUD (directory/menu/button)

### Auth & Permissions (3-layer RBAC)

- **JWT Auth** — Access Token (15min) + Refresh Token (7 days, HttpOnly Cookie)
- **Role Control** — `super_admin` / `admin` / `user`, customizable
- **Menu Permissions** — Different roles see different sidebar menus
- **Button Permissions** — Fine-grained permission codes (e.g. `user:create`, `role:delete`)
- **New users auto-assigned `user` role**

### Security

- **Login Rate Limiting** — IP-based sliding window, 5 req/min
- **CSP Headers** — Content-Security-Policy, X-Frame-Options, X-XSS-Protection, etc.
- **Token Refresh** — 401 auto-refresh via Axios interceptor
- **Password Hashing** — MD5 + Salt

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- PostgreSQL 14+

### Database Setup

```sql
CREATE DATABASE tool_go_dev;
```

```bash
psql -U postgres -d tool_go_dev -f manifest/sql/init.sql
```

### Backend

```bash
go mod tidy
go run main.go        # Dev server :8000, Swagger at /swagger
```

### Frontend

```bash
cd web
npm install
npm run dev           # Dev server :3000, Vite proxies /api → :8000
npm run build         # Production build (vue-tsc typecheck + vite build)
```

## Default Account

| Username | Password | Role |
|----------|----------|------|
| walter | walter | Super Admin (super_admin) |

> Seed data from `manifest/sql/init.sql`. Password hashing uses MD5 + Salt.

## RBAC Permission Model

### Roles

| Role | Description | Visible Menus | Button Permissions |
|------|-------------|----------------|-------------------|
| super_admin | Super Admin | All | All 7 permissions |
| admin | Admin | Toolbox, User Mgmt, Role Mgmt | user:create/delete/assign-roles, role:create/delete |
| user | Regular User | Toolbox | None |

### Permission Codes

| Code | Description |
|------|-------------|
| user:create | Create user |
| user:delete | Delete user |
| user:assign-roles | Assign roles |
| role:create | Create role |
| role:delete | Delete role |
| menu:create | Create menu |
| menu:delete | Delete menu |

## Configuration

### JWT

```yaml
# manifest/config/config.yaml
jwt:
  secret: "dev-jwt-secret-key-123456"  # Change in production
  expires: "24h"
  issuer: "tool-go"
```

### Environment Switch

```bash
export GF_GENV=prod   # Loads config.prod.yaml
```

## Docker

```bash
docker build -f manifest/docker/Dockerfile -t tool-go .
docker run -d -p 8000:8000 \
  -e DB_HOST=your-db-host \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your-password \
  -e DB_NAME=tool_go_prod \
  tool-go
```

## License

MIT
