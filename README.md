# 瓦特的工具站

基于 GoFrame v2 + Vue3 的开发和后台管理工具箱。

## 技术栈

### 后端
- **GoFrame v2** - 高性能 Go Web 框架
- **PostgreSQL** - 关系型数据库
- **JWT (JSON Web Token)** - 身份认证
- **RESTful API** - 标准接口风格

### 前端
- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - 类型安全
- **Vite** - 快速构建工具
- **Element Plus** - UI 组件库（简体中文）
- **Pinia** - 状态管理
- **Vue Router** - 路由管理

## 项目结构

```
tool-go/
├── api/                      # API 接口定义
│   └── v1/                   # API v1 版本
│       ├── auth.go           # 认证接口
│       ├── user.go           # 用户接口
│       └── role.go           # 角色接口
├── internal/                 # 内部代码
│   ├── cmd/                  # 命令行入口
│   ├── controller/           # 控制器层
│   ├── dao/                  # 数据访问层
│   ├── library/              # 公共库
│   │   └── jwt/              # JWT 工具
│   ├── logic/                # 业务逻辑层
│   ├── middleware/           # 中间件
│   │   ├── auth.go           # 认证中间件
│   │   └── permission.go     # 权限中间件
│   ├── model/                # 数据模型
│   │   ├── do/               # 数据对象
│   │   └── entity/           # 实体对象
│   ├── service/              # 服务接口
│   └── packed/               # 资源打包
├── manifest/                 # 清单文件
│   ├── config/               # 配置文件
│   │   ├── config.yaml       # 默认配置
│   │   ├── config.local.yaml # 本地开发配置
│   │   └── config.prod.yaml  # 生产环境配置
│   ├── sql/                  # 数据库脚本
│   │   └── init.sql
│   └── docker/               # Docker 配置
├── web/                      # Vue3 前端项目
│   ├── src/
│   │   ├── api/              # API 请求
│   │   ├── layouts/          # 布局组件
│   │   ├── router/           # 路由配置（含守卫）
│   │   ├── store/            # 状态管理
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   │       ├── login/        # 登录页
│   │       ├── dashboard/    # 工作台
│   │       ├── user/         # 用户管理
│   │       ├── role/         # 角色管理
│   │       └── error/        # 错误页面
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
# 安装依赖
go mod tidy

# 启动服务（默认端口 8000）
go run main.go
```

### 前端启动

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器（默认端口 3000）
npm run dev
```

### 一键启动（Windows）

```bash
start.bat
```

## 认证与权限

### JWT 认证流程

1. 用户通过 `POST /api/v1/login` 提交用户名密码
2. 后端验证成功后返回 JWT Token
3. 前端存储 Token，后续请求在 Header 中携带 `Authorization: Bearer <token>`
4. 后端通过认证中间件验证 Token 有效性
5. Token 过期或无效时，返回 401，前端自动跳转登录页

### 权限控制

- 只有 `super_admin` 角色可以创建新用户
- 路由级权限控制：前端根据用户角色动态显示菜单
- 接口级权限控制：后端中间件验证角色权限
- 无权限访问时返回 403，前端跳转无权限页面

### API 接口

#### 认证接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /api/v1/login | 用户登录 | 无 |
| GET | /api/v1/user/info | 获取当前用户信息 | 已登录 |
| POST | /api/v1/logout | 退出登录 | 已登录 |

#### 用户管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /api/v1/user | 创建用户 | super_admin |
| GET | /api/v1/user | 获取用户列表 | 已登录 |
| GET | /api/v1/user/{id} | 获取用户详情 | 已登录 |
| PUT | /api/v1/user/{id} | 更新用户 | 已登录 |
| DELETE | /api/v1/user/{id} | 删除用户 | 已登录 |

#### 角色管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /api/v1/role | 创建角色 | 已登录 |
| GET | /api/v1/role | 获取角色列表 | 已登录 |
| GET | /api/v1/role/{id} | 获取角色详情 | 已登录 |
| PUT | /api/v1/role/{id} | 更新角色 | 已登录 |
| DELETE | /api/v1/role/{id} | 删除角色 | 已登录 |

## 配置说明

### JWT 配置

在 `manifest/config/config.yaml` 中配置：

```yaml
jwt:
  secret: "your-secret-key"  # JWT 签名密钥（生产环境务必修改）
  expires: "24h"             # Token 过期时间
  issuer: "tool-go"          # 签发者
```

### 环境切换

通过 `GF_GENV` 环境变量切换：

```bash
# Windows PowerShell
$env:GF_GENV="prod"

# Linux/Mac
export GF_GENV=prod
```

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 超级管理员 (super_admin) |

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

## 代码生成

```bash
# 安装 gf CLI
go install github.com/gogf/gf/cmd/gf/v2@latest

# 生成 DAO 代码
gf gen dao
```

## License

MIT
