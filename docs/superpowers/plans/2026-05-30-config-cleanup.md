# 配置管理整理 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 整理项目配置管理，消除安全隐患，清理冗余配置，统一环境切换机制。

**Architecture:** 只保留 config.yaml（开发默认）+ config.prod.yaml.example（生产模板），移除所有 JWT fallback，Docker 通过 GF_GENV=prod 切换环境。

**Tech Stack:** GoFrame v2, Docker, Vite

---

## File Structure

### New Files
| File | Purpose |
|------|---------|
| `manifest/config/config.prod.yaml.example` | 生产配置模板（提交到 git） |

### Modified Files
| File | Change |
|------|--------|
| `.gitignore` | 忽略 config.prod.yaml，允许 config.prod.yaml.example |
| `manifest/config/config.yaml` | 精简注释 |
| `manifest/docker/Dockerfile` | 添加 ENV GF_GENV=prod，复制 example 模板 |
| `internal/controller/auth.go` | 移除 JWT fallback，config 为空时 panic |
| `internal/middleware/auth.go` | 移除 JWT fallback |
| `internal/controller/pageview.go` | 移除 JWT fallback |
| `main.go` | 精简注释 |

### Deleted Files
| File | Reason |
|------|--------|
| `manifest/config/config.prod.yaml` | 含敏感信息，改为 .example 模板 |
| `.env.example` | 死代码，无任何 Go 代码读取 |

---

## Task 1: 安全修复 — config.prod.yaml + JWT fallback

- [ ] **Step 1: 创建 config.prod.yaml.example 模板**

```yaml
# 生产环境配置模板
# 使用方法：复制此文件为 config.prod.yaml 并修改
# cp manifest/config/config.prod.yaml.example manifest/config/config.prod.yaml
# 启动：GF_GENV=prod ./tool-go

server:
  address: ":8000"

database:
  default:
    link: "pgsql:YOUR_USER:YOUR_PASSWORD@tcp(YOUR_HOST:5432)/YOUR_DB?sslmode=disable"
    debug: false
    maxIdle: 10
    maxOpen: 100
    maxLifetime: "30s"

logger:
  level: "warning"
  stdout: false

jwt:
  secret: "CHANGE_ME_TO_A_RANDOM_STRING"
  expires: "15m"
  issuer: "tool-go"
```

- [ ] **Step 2: 更新 .gitignore**

Change:
```
manifest/config/config.*.yaml
!manifest/config/config.local.yaml
!manifest/config/config.prod.yaml
```
To:
```
manifest/config/config.*.yaml
!manifest/config/config.local.yaml
!manifest/config/config.prod.yaml.example
```

- [ ] **Step 3: 删除 config.prod.yaml from git tracking**

```bash
git rm --cached manifest/config/config.prod.yaml
```

- [ ] **Step 4: 移除 JWT fallback in auth.go**

Read the file. Find the 2 places where `secret = "tool-go-jwt-secret-key-change-in-production"` is used as fallback. Change to:

```go
if secret == "" {
    g.Log().Fatal(ctx, "JWT secret not configured in config.yaml")
}
```

Remove the fallback assignment.

- [ ] **Step 5: 移除 JWT fallback in middleware/auth.go**

Same change — remove fallback, add Fatal.

- [ ] **Step 6: 移除 JWT fallback in pageview.go**

Same change — remove fallback, add Fatal.

- [ ] **Step 7: Commit**

```bash
git add manifest/config/config.prod.yaml.example .gitignore internal/controller/auth.go internal/middleware/auth.go internal/controller/pageview.go
git rm --cached manifest/config/config.prod.yaml
git commit -m "fix: remove prod config from git, remove JWT fallback secrets"
```

---

## Task 2: Docker 修复

- [ ] **Step 1: 更新 Dockerfile**

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/tool-go .

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/tool-go .
COPY --from=builder /app/manifest/config ./manifest/config
COPY --from=builder /app/resource ./resource

ENV GF_GENV=prod

EXPOSE 8000
CMD ["./tool-go"]
```

Changes:
- Add `ENV GF_GENV=prod` after COPY
- Remove `RUN apk add --no-cache git` from builder (only needed if go mod download needs git — actually it does for private repos, keep it)

- [ ] **Step 2: Commit**

```bash
git add manifest/docker/Dockerfile
git commit -m "fix: set GF_GENV=prod in Dockerfile for production config"
```

---

## Task 3: 清理死代码 + 精简注释

- [ ] **Step 1: 删除 .env.example**

```bash
git rm .env.example
```

- [ ] **Step 2: 精简 main.go 注释**

当前 main.go 有 78 行，实际代码 3 行。精简为：

```go
package main

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2" // PostgreSQL 驱动
	"github.com/gogf/gf/v2/os/gctx"
	_ "tool-go/internal/logic" // 注册业务逻辑
	"tool-go/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
```

- [ ] **Step 3: Commit**

```bash
git add main.go
git rm .env.example
git commit -m "chore: remove dead .env.example, slim down main.go comments"
```

---

## Summary

| Task | Description | Files |
|------|-------------|-------|
| 1 | 安全修复: prod config gitignore + JWT fallback | 5 modified, 1 new, 1 deleted |
| 2 | Docker 修复: GF_GENV=prod | 1 modified |
| 3 | 清理: 删除死代码 + 精简注释 | 1 modified, 1 deleted |
