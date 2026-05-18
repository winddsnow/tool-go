# Tool Go — Developer Toolbox

A development and admin toolbox built with GoFrame v2 + Vue 3 + TypeScript + Element Plus.

## Tech Stack

### Backend

- **GoFrame v2** — High-performance Go web framework
- **PostgreSQL** — Relational database
- **JWT (golang-jwt/jwt/v5)** — Authentication
- **RESTful API** — Standard API style
- **Swagger** — API docs (dev mode)

### Frontend

- **Vue 3** — Progressive JavaScript framework
- **TypeScript** — Type safety
- **Vite** — Fast build tool
- **Element Plus** — UI component library
- **Pinia** — State management
- **Vue Router** — Route management

## Features

### Toolbox (11 tools, client-side)

| Category | Tools |
|----------|-------|
| Text | JSON Formatter, Text Diff, Regex Tester, Case Converter |
| Encoding | Base64 Encoder/Decoder, Hash (MD5/SHA1/SHA256) |
| Generation | Password Generator, Mock Data Generator (9 types), UUID Generator, QR Code Generator |
| Conversion | Timestamp Converter (16 timezones) |

### Admin Panel

- **Dashboard** — System stats (users, roles, visits)
- **User Management** — CRUD, role assignment, paginated search
- **Role Management** — CRUD, permission control, paginated search

### Auth & Permissions

- JWT login authentication
- RBAC role-based access control (`super_admin` / `admin`)
- Dual-layer guard: backend middleware + frontend route guard

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
