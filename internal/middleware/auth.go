// Package middleware 提供 HTTP 中间件（Middleware）。
// 中间件是 GoFrame 框架中的一种请求拦截机制，在请求到达 controller 之前
// 或响应返回客户端之前执行预处理/后处理逻辑。
// 核心模式：r.Middleware.Next() 调用下一个中间件或最终的 controller 处理函数。

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"

	"tool-go/internal/library/jwt"
)

// 上下文键名常量，用于通过 context.WithValue 在中间件链中传递数据。
// Go 的 context 是携带请求范围值的容器，贯穿整个请求生命周期。
// context.WithValue(ctx, key, value) 将值附加到 context，
// 后续中间件和 controller 通过 ctx.Value(key) 读取。
const (
	CtxUserId   = "userId"   // 当前登录用户的 ID
	CtxUsername = "username" // 当前登录用户的用户名
	CtxRoles       = "roles"       // 当前登录用户的角色列表
	CtxPermissions = "permissions" // 当前登录用户的权限列表
)

// Auth JWT 认证中间件。
// 执行顺序（GoFrame 中间件链）：
//   请求到达 → Auth 中间件 → [处理 JWT] → r.Middleware.Next() → controller
// Auth 的职责：
//   1. 从请求头提取 Bearer Token
//   2. 使用 JWT 库解析和验证 Token
//   3. 将解析出的用户信息存入 context
//   4. 调用 Next() 继续请求处理
//   如果验证失败，直接返回 401 并终止请求链（不调用 Next()）。
func Auth(r *ghttp.Request) {
	// 从配置中读取 JWT 密钥。g.Cfg().MustGet(ctx, "jwt").MapStrStr()
	// 通过 GoFrame 配置中心获取配置（支持 yaml/json/toml 等格式）。
	// MustGet 在 key 不存在时会 panic，因此需要兜底默认值。
	ctx := r.GetCtx()
	jwtConfig := g.Cfg().MustGet(ctx, "jwt").MapStrStr()
	secret := jwtConfig["secret"]
	if secret == "" {
		g.Log().Fatal(r.GetCtx(), "JWT secret not configured in config.yaml")
	}

	// 创建 JWT 工具实例，传入密钥。jwt 包是本项目自定义的 JWT 封装。
	j := jwt.New(secret, 0, "")

	// 从 HTTP 请求头中提取 Authorization 字段。
	// 标准格式：Authorization: Bearer <token>
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// WriteStatus 设置 HTTP 响应状态码（401 Unauthorized）。
		// WriteJsonExit 写入 JSON 格式的响应体并立即终止当前请求处理。
		// 注意：没有调用 r.Middleware.Next()——这表示请求被拦截，不会继续向后传递。
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.WriteJsonExit(g.Map{
			"code":    gcode.CodeNotAuthorized.Code(),
			"message": "未登录或登录已过期",
		})
		return
	}

	// 去掉 "Bearer " 前缀，解析出实际的 Token 字符串。
	// strings.TrimPrefix 如果字符串以 prefix 开头则去除，否则返回原字符串。
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		// 如果没有 "Bearer " 前缀，说明格式不正确。
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.WriteJsonExit(g.Map{
			"code":    gcode.CodeNotAuthorized.Code(),
			"message": "token格式错误",
		})
		return
	}

	// 使用 JWT 工具解析 Token，验证签名和有效期。
	// 解析成功返回 claims（包含用户 ID、用户名、角色列表等）。
	// ParseToken 内部做了签名校验、过期检查等安全验证。
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.WriteJsonExit(g.Map{
			"code":    gcode.CodeNotAuthorized.Code(),
			"message": gstr.TrimLeftStr(err.Error(), "jwt: "),
		})
		return
	}

	// context.WithValue 将用户信息注入到 context 中。
	// context 是 Go 标准库中携带请求范围数据的容器。
	// 特点：
	//   - context.WithValue 返回新的 context（不可变，基于原 context 派生）。
	//   - 新 context 包含原 context 的所有值 + 新增的键值对。
	//   - context 是协程安全的，多个 goroutine 可以同时读取。
	// 这里注入 userId/username/roles，后续中间件和 controller 通过
	// middleware.GetUserId(ctx) 等辅助函数读取。
	ctx = context.WithValue(ctx, CtxUserId, claims.UserId)
	ctx = context.WithValue(ctx, CtxUsername, claims.Username)
	ctx = context.WithValue(ctx, CtxRoles, claims.Roles)
	ctx = context.WithValue(ctx, CtxPermissions, claims.Permissions)
	r.SetCtx(ctx)

	// r.Middleware.Next() 是 GoFrame 中间件链的核心机制：
	// 调用 Next() 表示当前中间件的处理已完成，将控制权交给链中的下一个处理器。
	// 如果不调用 Next()，请求链在此终止（如验证失败时）。
	// 所有中间件形成一个"洋葱模型"：
	//   请求 → Auth → [Permission → Controller → Permission 后] → Auth 后 → 响应
	r.Middleware.Next()
}

// GetUserId 从 context 中提取当前登录用户的 ID。
// context.WithValue 存入的值通过 ctx.Value(key) 读取，返回 interface{}，
// 因此需要类型断言 .(uint64) 转换为具体类型。如果 key 不存在，返回 0。
func GetUserId(ctx context.Context) uint64 {
	if v := ctx.Value(CtxUserId); v != nil {
		return v.(uint64)
	}
	return 0
}

// GetUsername 从 context 中提取当前登录用户的用户名。
func GetUsername(ctx context.Context) string {
	if v := ctx.Value(CtxUsername); v != nil {
		return v.(string)
	}
	return ""
}

// GetRoles 从 context 中提取当前登录用户的角色列表。
func GetRoles(ctx context.Context) []string {
	if v := ctx.Value(CtxRoles); v != nil {
		return v.([]string)
	}
	return []string{}
}

// HasRole 检查当前用户是否拥有指定的单个角色。
// 遍历角色列表逐一比对，找到即返回 true。
func HasRole(ctx context.Context, role string) bool {
	roles := GetRoles(ctx)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查当前用户是否拥有指定角色列表中的任意一个。
// ...string 是 Go 的可变参数（variadic），调用时可以传任意多个参数。
// 这个方法被 Permission 中间件使用，用于灵活的角色权限控制。
func HasAnyRole(ctx context.Context, roles ...string) bool {
	for _, role := range roles {
		if HasRole(ctx, role) {
			return true
		}
	}
	return false
}

func GetPermissions(ctx context.Context) []string {
	permissions, _ := ctx.Value(CtxPermissions).([]string)
	return permissions
}

func HasPermission(ctx context.Context, permission string) bool {
	for _, p := range GetPermissions(ctx) {
		if p == permission {
			return true
		}
	}
	return false
}

func HasAnyPermission(ctx context.Context, permissions ...string) bool {
	userPerms := GetPermissions(ctx)
	for _, required := range permissions {
		for _, p := range userPerms {
			if p == required {
				return true
			}
		}
	}
	return false
}
