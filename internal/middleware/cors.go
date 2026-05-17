// Package middleware 提供 HTTP 中间件。
// CORS 中间件通常在请求链的最外层注册，确保所有请求都能正确处理跨域。

package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// CORS 跨域资源共享（Cross-Origin Resource Sharing）中间件。
// 当浏览器端（如 Vue 前端运行在 localhost:3000）向不同源的后端（:8000）发起请求时，
// 浏览器会先发送 OPTIONS 预检请求（Preflight Request），
// 服务器必须返回适当的 CORS 响应头，浏览器才会允许实际的 HTTP 请求。
//
// r.Response.CORSDefault() 是 GoFrame 的快捷方法，自动设置以下响应头：
//   Access-Control-Allow-Origin:  *
//   Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS,HEAD,PATCH
//   Access-Control-Allow-Headers: Origin,Content-Type,Accept,Authorization
//   Access-Control-Max-Age:       86400
// 如果前端有自定义需求，可以改用 r.Response.CORS() 传入自定义 CORS 配置。
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	// 调用 Next() 将请求传递到链中的下一个处理器。
	// CORS 中间件不拦截任何请求（对所有请求都放行），只做响应头设置。
	r.Middleware.Next()
}
