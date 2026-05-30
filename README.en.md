# Tool Go — Developer Toolbox

<p align="center">
  <strong>The Swiss Army Knife for Developers — Toolbox + Admin Panel + RBAC Permission System</strong>
</p>

<p align="center">
  <a href="https://github.com/winddsnow/tool-go/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/vue-3-4FC08D?logo=vue.js" alt="Vue 3">
  <img src="https://img.shields.io/badge/typescript-5-3178C6?logo=typescript" alt="TypeScript">
  <img src="https://img.shields.io/badge/postgreSQL-14-336791?logo=postgresql" alt="PostgreSQL">
</p>

---

**Tool Go** is a full-stack developer toolbox and admin management system. It integrates 12 commonly-used development tools (JSON Formatter, Regex Tester, Base64, Hash Encryption, etc.) with a complete RBAC permission management system, serving both as an everyday productivity tool and an extensible admin platform.

Built with **GoFrame v2 + Vue 3 + TypeScript + Element Plus**, the backend follows GoFrame's RESTful API design, the frontend uses Vite for blazing-fast builds, and supports enterprise-grade features like dynamic menus, button-level permission control, and JWT dual-token authentication.

> This project was developed collaboratively by developer **winddsnow** and AI large language models (**MiMo v2.5 Pro**, **DeepSeek V4**). The AI models actively participated in architecture design, code implementation, debugging, and feature iteration throughout the entire development process — a practical demonstration of human-AI collaborative development.

## Highlights

- **12 Developer Tools** — JSON Formatter, Text Diff, Regex Tester, Base64, Hash, Password Generator, UUID, QR Code, and more — all client-side
- **Dynamic Menu System** — Data-driven menus, different roles see different sidebars, 3-level support (Directory / Menu / Button)
- **Fine-grained RBAC** — Role → Menu + Role → Permission code, button-level control (e.g. `user:create`, `role:delete`)
- **JWT Dual-Token Auth** — Access Token (15min) + Refresh Token (7 days, HttpOnly Cookie), auto-refresh on 401
- **Security Hardening** — Login rate limiting (5 req/min/IP), CSP headers, XSS protection, password MD5+Salt hashing
- **Admin Panel** — User, Role, and Menu management with visual permission assignment
- **Instant Usability** — New users auto-assigned basic role, login and use the toolbox immediately

## Tech Stack

### Backend

| Technology | Description |
|------------|-------------|
| **GoFrame v2** | High-performance Go web framework, RESTful API, ORM, middleware |
| **PostgreSQL** | Relational database, 7 core tables |
| **JWT** | golang-jwt/jwt/v5, dual-token authentication |
| **Swagger** | Auto-generated API docs (dev mode) |

### Frontend

| Technology | Description |
|------------|-------------|
| **Vue 3** | Composition API, reactive, component-based |
| **TypeScript** | Type safety, compile-time error checking |
| **Vite** | Lightning-fast builds, HMR in dev |
| **Element Plus** | Enterprise-grade UI component library |
| **Pinia** | Lightweight state management with persistence |
| **Vue Router** | Dynamic routes, menu-driven registration |

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

## Collaborators

This project was developed collaboratively by **winddsnow** and AI large language models. The AI models actively participated in architecture design, code implementation, debugging, and feature iteration throughout the entire development process.

| Collaborator | Model | Contribution |
|--------------|-------|--------------|
| **winddsnow** | — | Project initiation, requirements, code review, testing |
| **MiMo v2.5 Pro** | Xiaomi MiMo | Code implementation, feature development, debugging, documentation |
| **DeepSeek V4** | DeepSeek | Architecture design, solution review, early feature development |

> This project is a practice of human-AI collaborative development. AI models have significantly improved efficiency in daily coding, debugging, and documentation, while human developers play an irreplaceable role in requirements control, architectural decisions, and quality assurance.

## License

MIT
