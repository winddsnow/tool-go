// Package middleware 实现请求预处理中间件。
// Permission 中间件依赖于 Auth 中间件先执行（先解析 JWT 再校验角色）。

package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Permission 是一个"中间件工厂"——它接收所需角色列表，返回一个中间件函数。
// 这是 Go 中闭包（closure）的典型用法：Permission 的参数 requiredRoles 被
// 内部匿名函数"捕获"（闭包引用），即使 Permission 返回后仍可访问。
//
// 用法：在路由注册时调用，如 middleware.Permission("super_admin", "admin")
// 返回的中间件函数会：
//   1. 从 context 中读取用户角色（由 Auth 中间件提前注入）
//   2. 检查用户是否拥有 requiredRoles 中的任意一个
//   3. 如果没有 → 返回 403 Forbidden 并终止请求链
//   4. 如果有   → 调用 r.Middleware.Next() 继续处理
func Permission(requiredRoles ...string) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		// 从 context 中获取用户角色列表（由 Auth 中间件在之前注入）。
		// HasAnyRole 遍历角色列表，只要有一个匹配就返回 true。
		// requiredRoles 通过闭包从外部函数访问。
		ctx := r.GetCtx()
		if !HasAnyRole(ctx, requiredRoles...) {
			// 角色不匹配，返回 403 Forbidden。
			// WriteStatus 设置状态码，WriteJsonExit 返回 JSON 错误信息。
			// 没有调用 Next()，请求在该中间件终止。
			r.Response.WriteStatus(http.StatusForbidden)
			r.Response.WriteJsonExit(g.Map{
				"code":    gcode.CodeNotAuthorized.Code(),
				"message": "没有权限访问",
			})
			return
		}
		// 角色匹配，继续执行后续中间件或 controller。
		r.Middleware.Next()
	}
}
