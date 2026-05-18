# 瓦特的工具站

基于 GoFrame v2 + Vue 3 + TypeScript + Element Plus 的开发和后台管理工具箱。

## 技术栈

### 后端

- **GoFrame v2** — 高性能 Go Web 框架
- **PostgreSQL** — 关系型数据库
- **JWT (golang-jwt/jwt/v5)** — 身份认证
- **RESTful API** — 标准接口风格
- **Swagger** — 接口文档（开发模式）

### 前端

- **Vue 3** — 渐进式 JavaScript 框架
- **TypeScript** — 类型安全
- **Vite** — 快速构建工具
- **Element Plus** — UI 组件库
- **Pinia** — 状态管理
- **Vue Router** — 路由管理

## 功能概览

### 工具箱（11 个工具，纯前端）

| 分类 | 工具 |
|------|------|
| 文本处理 | JSON 格式化、文本对比(Diff)、正则表达式测试、大小写/Naming Case 转换 |
| 编码加密 | Base64 编解码、哈希加密(MD5/SHA1/SHA256) |
| 生成类 | 密码生成器、随机数据生成器(9种数据)、UUID 生成器、二维码生成 |
| 转换类 | 时间戳转换(16个时区) |

### 管理后台

- **仪表盘** — 系统概览统计（用户数、角色数、访问量）
- **用户管理** — 用户 CRUD、角色分配、分页搜索
- **角色管理** — 角色 CRUD、权限控制、分页搜索

### 认证与权限

- JWT 登录认证
- RBAC 角色权限控制（`super_admin` / `admin`）
- 后端中间件 + 前端路由守卫双层校验

## 项目结构

```
tool-go/
├── api/v1/                   # API 请求/响应结构体 + 路由声明
├── internal/
│   ├── cmd/                  # 应用入口、路由绑定
│   ├── controller/           # HTTP 控制器层
│   ├── dao/                  # 数据访问层（手写 DAO + 列名常量）
│   ├── library/
│   │   ├── jwt/              # JWT 创建/解析
│   │   └── password/         # MD5 + Salt 密码哈希
│   ├── logic/                # 业务逻辑层
│   ├── middleware/           # CORS、Auth、Permission 中间件
│   ├── model/
│   │   ├── do/               # 数据操作对象
│   │   └── entity/           # 数据库实体
│   └── service/              # 服务接口层
├── manifest/
│   ├── config/               # 应用配置
│   ├── sql/                  # 数据库初始化脚本
│   └── docker/               # Dockerfile
├── web/
│   ├── src/
│   │   ├── api/              # Axios 请求封装
│   │   ├── layouts/          # 布局组件
│   │   ├── router/           # 路由 + 导航守卫
│   │   ├── store/            # Pinia 状态管理
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   └── ...
└── ...
```

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

| 模块 | 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|------|
| 认证 | POST | /api/v1/login | 用户登录 | 无 |
| | GET | /api/v1/user/info | 获取当前用户信息 | 已登录 |
| | POST | /api/v1/logout | 退出登录 | 已登录 |
| 用户 | POST | /api/v1/user | 创建用户 | super_admin, admin |
| | DELETE | /api/v1/user/{id} | 删除用户 | super_admin, admin |
| | PUT | /api/v1/user/{id} | 更新用户 | super_admin, admin |
| | GET | /api/v1/user/{id} | 获取用户详情 | 已登录 |
| | GET | /api/v1/user | 用户列表(分页/搜索) | 已登录 |
| | GET | /api/v1/user/{id}/roles | 获取用户角色 | 已登录 |
| | PUT | /api/v1/user/{id}/roles | 分配角色 | super_admin, admin |
| 角色 | POST | /api/v1/role | 创建角色 | super_admin, admin |
| | DELETE | /api/v1/role/{id} | 删除角色 | super_admin, admin |
| | PUT | /api/v1/role/{id} | 更新角色 | super_admin, admin |
| | GET | /api/v1/role/{id} | 获取角色详情 | 已登录 |
| | GET | /api/v1/role | 角色列表(分页/搜索) | 已登录 |
| 仪表盘 | GET | /api/v1/dashboard/stats | 系统统计数据 | 已登录 |
| 页面访问 | POST | /api/v1/pageview/track | 记录页面访问 | 无 |
| | GET | /api/v1/pageview/stats | 访问统计 | 已登录 |
| 工具 | POST | /api/v1/tools/mock-data | 生成模拟数据 | 无 |

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| walter | walter | 超级管理员 (super_admin) |

> 种子数据来自 `manifest/sql/init.sql`，密码哈希使用 MD5 + Salt。

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

## License

MIT
