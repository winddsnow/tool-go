# 瓦特的工具站

<p align="center">
  <strong>开发者的瑞士军刀 — 工具箱 + 后台管理 + RBAC 权限系统</strong>
</p>

<p align="center">
  <a href="https://github.com/winddsnow/tool-go/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/vue-3-4FC08D?logo=vue.js" alt="Vue 3">
  <img src="https://img.shields.io/badge/typescript-5-3178C6?logo=typescript" alt="TypeScript">
  <img src="https://img.shields.io/badge/postgreSQL-14-336791?logo=postgresql" alt="PostgreSQL">
</p>

---

**瓦特的工具站** 是一个面向开发者的全栈工具箱和后台管理系统。它将 12 个常用开发工具（JSON 格式化、正则测试、Base64、哈希加密等）与一个完整的 RBAC 权限管理系统整合在一起，既是日常开发的效率工具，也是一个可扩展的后台管理平台。

项目采用 **GoFrame v2 + Vue 3 + TypeScript + Element Plus** 的现代化技术栈，后端基于 GoFrame 的 RESTful API 设计，前端使用 Vite 构建，支持动态菜单、按钮级权限控制、JWT 双 Token 认证等企业级特性。

> 本项目由开发者 **winddsnow** 与 AI 大模型（**MiMo v2.5 Pro**、**DeepSeek V4**）协作开发完成。大模型全程参与了架构设计、代码编写、问题排查和功能迭代，是人机协作开发模式的一次实践。

## 核心特性

- **12 个开发工具** — JSON 格式化、文本对比、正则测试、Base64、哈希加密、密码生成、UUID、二维码等，纯前端实现
- **动态菜单系统** — 菜单数据驱动，不同角色看到不同侧边栏，支持三级菜单（目录/菜单/按钮）
- **细粒度 RBAC 权限** — 角色 → 菜单 + 角色 → 权限码，按钮级权限控制（如 `user:create`、`role:delete`）
- **JWT 双 Token 认证** — Access Token (15分钟) + Refresh Token (7天, HttpOnly Cookie)，401 自动刷新
- **安全防护** — 登录限流 (5次/分钟/IP)、CSP 安全头、XSS 防护、密码 MD5+Salt 哈希
- **管理后台** — 用户管理、角色管理、菜单管理，可视化配置权限分配
- **新建用户即用** — 自动分配基础角色，登录即可使用工具箱

## 技术栈

### 后端

| 技术 | 说明 |
|------|------|
| **GoFrame v2** | 高性能 Go Web 框架，RESTful API、ORM、中间件 |
| **PostgreSQL** | 关系型数据库，7 张核心表 |
| **JWT** | golang-jwt/jwt/v5，双 Token 认证 |
| **Swagger** | 自动生成 API 文档（开发模式） |

### 前端

| 技术 | 说明 |
|------|------|
| **Vue 3** | Composition API、响应式、组件化 |
| **TypeScript** | 类型安全，编译期错误检查 |
| **Vite** | 极速构建，开发热更新 |
| **Element Plus** | 企业级 UI 组件库 |
| **Pinia** | 轻量状态管理，持久化存储 |
| **Vue Router** | 动态路由，菜单驱动注册 |

## 功能概览

### 工具箱（12 个工具，纯前端）

| 分类 | 工具 |
|------|------|
| 文本处理 | JSON 格式化、文本对比(Diff)、正则表达式测试、大小写/Naming Case 转换 |
| 编码加密 | Base64 编解码、哈希加密(MD5/SHA1/SHA256) |
| 生成类 | 密码生成器、随机数据生成器(9种数据)、UUID 生成器、二维码生成 |
| 转换类 | 时间戳转换(16个时区) |

### 管理后台

- **仪表盘** — 系统概览统计（用户数、角色数、访问量）
- **用户管理** — 用户 CRUD、角色分配、分页搜索
- **角色管理** — 角色 CRUD、菜单分配、权限分配、分页搜索
- **菜单管理** — 动态菜单树 CRUD（目录/菜单/按钮三级）

### 认证与权限（RBAC 三层模型）

- **JWT 认证** — Access Token (15分钟) + Refresh Token (7天, HttpOnly Cookie)
- **角色控制** — `super_admin` / `admin` / `user`，可自定义扩展
- **菜单权限** — 不同角色看到不同侧边栏菜单
- **按钮权限** — 细粒度权限码控制（如 `user:create`、`role:delete`）
- **新建用户自动分配 `user` 角色**

### 安全特性

- **登录限流** — IP 级别滑动窗口，5次/分钟
- **CSP 安全头** — Content-Security-Policy、X-Frame-Options、X-XSS-Protection 等
- **Token 刷新** — 401 时自动调用 refresh 接口，无感续期
- **密码哈希** — MD5 + Salt

## 项目结构

```
tool-go/
├── api/v1/                   # API 请求/响应结构体 + 路由声明
├── internal/
│   ├── cmd/                  # 应用入口、路由绑定
│   ├── controller/           # HTTP 控制器层
│   ├── dao/                  # 数据访问层（手写 DAO + 列名常量）
│   ├── library/
│   │   ├── jwt/              # JWT 创建/解析（Access + Refresh）
│   │   └── password/         # MD5 + Salt 密码哈希
│   ├── logic/                # 业务逻辑层
│   ├── middleware/           # CORS、Auth、Permission、RateLimit、SecurityHeaders
│   ├── model/
│   │   ├── do/               # 数据操作对象
│   │   └── entity/           # 数据库实体
│   └── service/              # 服务接口层
├── manifest/
│   ├── config/               # 应用配置
│   ├── sql/                  # 数据库初始化脚本 + 迁移脚本
│   └── docker/               # Dockerfile
├── web/
│   ├── src/
│   │   ├── api/              # Axios 请求封装
│   │   ├── layouts/          # 布局组件（动态侧边栏）
│   │   ├── router/           # 动态路由（菜单驱动）
│   │   ├── store/            # Pinia 状态管理（含权限检查）
│   │   ├── utils/            # 工具函数（菜单转换、请求拦截）
│   │   └── views/            # 页面视图
│   └── ...
└── ...
```

## 数据库设计

| 表名 | 说明 |
|------|------|
| `user` | 用户表 |
| `role` | 角色表 |
| `user_role` | 用户-角色关联 |
| `menu` | 动态菜单表（目录/菜单/按钮三级） |
| `role_menu` | 角色-菜单关联 |
| `permission` | 权限码表（如 `user:create`） |
| `role_permission` | 角色-权限关联 |
| `page_view` | 页面访问埋点 |

## 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+
- PostgreSQL 14+

### 数据库初始化

```sql
CREATE DATABASE tool_go_dev;
```

```bash
psql -U postgres -d tool_go_dev -f manifest/sql/init.sql
```

### 后端启动

```bash
go mod tidy
go run main.go        # 开发服务器 :8000，Swagger 文档 /swagger
```

### 前端启动

```bash
cd web
npm install
npm run dev           # 开发服务器 :3000，Vite 代理 /api → :8000
npm run build         # 生产构建（vue-tsc 类型检查 + vite 打包）
```

## API 接口

### 认证

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /api/v1/login | 用户登录（限流 5次/分钟） | 无 |
| POST | /api/v1/refresh | 刷新 Access Token | 无（需 Refresh Token Cookie） |
| POST | /api/v1/logout | 退出登录（清除 Cookie） | 已登录 |
| GET | /api/v1/user/info | 获取当前用户信息（含菜单+权限） | 已登录 |

### 用户管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /api/v1/user | 创建用户（自动分配 user 角色） | user:create |
| GET | /api/v1/user | 用户列表(分页/搜索) | 已登录 |
| GET | /api/v1/user/{id} | 获取用户详情 | 已登录 |
| PUT | /api/v1/user/{id} | 更新用户 | 已登录 |
| DELETE | /api/v1/user/{id} | 删除用户 | user:delete |
| GET | /api/v1/user/{id}/roles | 获取用户角色 | 已登录 |
| PUT | /api/v1/user/{id}/roles | 分配角色 | user:assign-roles |

### 角色管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /api/v1/role | 创建角色 | role:create |
| GET | /api/v1/role | 角色列表(分页/搜索) | 已登录 |
| GET | /api/v1/role/{id} | 获取角色详情 | 已登录 |
| PUT | /api/v1/role/{id} | 更新角色 | 已登录 |
| DELETE | /api/v1/role/{id} | 删除角色 | role:delete |
| GET | /api/v1/role/{id}/menus | 获取角色菜单 | 已登录 |
| PUT | /api/v1/role/{id}/menus | 分配菜单 | 已登录 |
| GET | /api/v1/role/{id}/permissions | 获取角色权限 | 已登录 |
| PUT | /api/v1/role/{id}/permissions | 分配权限 | 已登录 |

### 菜单管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /api/v1/menu | 创建菜单 | menu:create |
| GET | /api/v1/menu | 菜单列表(分页/搜索) | 已登录 |
| GET | /api/v1/menu/{id} | 获取菜单详情 | 已登录 |
| PUT | /api/v1/menu/{id} | 更新菜单 | 已登录 |
| DELETE | /api/v1/menu/{id} | 删除菜单 | menu:delete |
| GET | /api/v1/menu/user | 获取当前用户菜单树 | 已登录 |

### 权限管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/permission | 权限列表 | 已登录 |

### 其他

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/dashboard/stats | 系统统计数据 | 已登录 |
| POST | /api/v1/pageview/track | 记录页面访问 | 无 |
| GET | /api/v1/pageview/stats | 访问统计 | 已登录 |
| POST | /api/v1/tools/mock-data | 生成模拟数据 | 无 |

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| walter | walter | 超级管理员 (super_admin) |

> 种子数据来自 `manifest/sql/init.sql`，密码哈希使用 MD5 + Salt。

## RBAC 权限模型

### 角色

| 角色 | 说明 | 可见菜单 | 按钮权限 |
|------|------|----------|----------|
| super_admin | 超级管理员 | 全部 | 全部 7 个 |
| admin | 管理员 | 工具箱、用户管理、角色管理 | user:create/delete/assign-roles, role:create/delete |
| user | 普通用户 | 工具箱 | 无 |

### 权限码

| 权限码 | 说明 |
|--------|------|
| user:create | 创建用户 |
| user:delete | 删除用户 |
| user:assign-roles | 分配角色 |
| role:create | 创建角色 |
| role:delete | 删除角色 |
| menu:create | 创建菜单 |
| menu:delete | 删除菜单 |

### 新增页面/按钮配置流程

```sql
-- 1. 插入菜单
INSERT INTO menu (parent_id, name, path, component, icon, sort, type)
VALUES (0, '公告管理', '/notice', 'views/notice/index.vue', 'Bell', 3, 2);

-- 2. 给角色分配菜单
INSERT INTO role_menu (role_id, menu_id) VALUES (2, 最新menu_id);

-- 3. 插入按钮权限
INSERT INTO permission (code, name, menu_id) VALUES
('notice:create', '创建公告', 最新menu_id),
('notice:delete', '删除公告', 最新menu_id);

-- 4. 给角色分配权限
INSERT INTO role_permission (role_id, permission_id) VALUES (2, 最新permission_id);
```

```vue
<!-- 5. 前端按钮使用权限检查 -->
<el-button v-if="userStore.hasPermission('notice:create')">创建公告</el-button>
```

```go
// 6. 后端路由保护
auth.Group("/notice", func(notice *ghttp.RouterGroup) {
    notice.Middleware(middleware.PermissionCode("notice:create", "notice:delete"))
})
```

## 配置说明

### JWT

```yaml
# manifest/config/config.yaml
jwt:
  secret: "dev-jwt-secret-key-123456"  # 生产环境务必修改
  expires: "24h"
  issuer: "tool-go"
```

### 环境切换

```bash
export GF_GENV=prod   # 加载 config.prod.yaml
```

## Docker 部署

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

## 协作者

本项目由开发者 **winddsnow** 与 AI 大模型协作开发完成。大模型全程参与了架构设计、代码编写、问题排查和功能迭代。

| 协作者 | 模型 | 贡献 |
|--------|------|------|
| **winddsnow** | — | 项目发起、需求定义、代码审查、测试验收 |
| **MiMo v2.5 Pro** | 小米 MiMo | 代码实现、功能开发、问题排查、文档编写 |
| **DeepSeek V4** | DeepSeek | 架构设计、方案评审、早期功能开发 |

> 本项目是人机协作开发模式的一次实践。AI 大模型在日常编码、调试、文档编写等环节提供了显著的效率提升，而人类开发者在需求把控、架构决策和质量验收方面发挥着不可替代的作用。

## License

MIT
