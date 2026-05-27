// ============================================================
// package cmd — 命令行入口模块
// ------------------------------------------------------------
// Go 语言中，internal 目录是一个特殊的包名约定。
// internal 包只能被其父目录下的代码导入，不能被外部项目导入。
// 这是一种 Go 语言级别的访问控制机制，用于保护内部实现细节。
// 比如 tool-go/internal/xxx 只能被 tool-go/ 下的代码使用。
//
// cmd 包通常用于存放程序的"命令"定义。
// 在这里定义了 HTTP 服务启动命令的路由、中间件等。
// ============================================================
package cmd

import (
	// context — Go 标准库，用于管理协程的生命周期和传递请求范围的值。
	// context.Context 是 Go 中处理超时、取消信号和传递请求级别数据的标准方式。
	"context"

	// g — GoFrame 框架核心包，提供 g.Server() 等服务入口
	"github.com/gogf/gf/v2/frame/g"

	// ghttp — GoFrame 的 HTTP 服务包，提供 RouterGroup（路由组）、
	// MiddlewareHandlerResponse（标准响应处理中间件）等
	"github.com/gogf/gf/v2/net/ghttp"

	// gcmd — GoFrame 的命令行管理包，用于定义和管理子命令。
	// 在这里定义了一个 "main" 命令来启动 HTTP 服务。
	"github.com/gogf/gf/v2/os/gcmd"

	// controller — 项目的控制器包，每个控制器包含一组相关的接口处理函数。
	// 如 controller.Auth 包含 Login、Logout、GetUserInfo 等方法。
	"tool-go/internal/controller"

	// middleware — 项目的中间件包，包含 CORS（跨域）、Auth（JWT认证）、
	// Permission（角色权限）等中间件。
	"tool-go/internal/middleware"
)

// ============================================================
// var — Go 中声明变量的关键字
// Main — 程序的主命令，类型为 gcmd.Command（GoFrame 命令行命令结构体）
// ============================================================
var (
	// gcmd.Command — GoFrame 的命令行命令结构体
	// 在这里定义了 HTTP 服务器的启动命令
	Main = gcmd.Command{
		// Name — 命令名称，命令行中通过该名称调用
		Name: "main",
		// Usage — 使用说明
		Usage: "main",
		// Brief — 简短描述
		Brief: "start http server",

		// Func — 命令的执行函数，当命令被调用时执行
		// func 是 Go 的关键字，用于定义函数
		// 参数:
		//   ctx context.Context — Go 的上下文对象，用于传递超时、取消信号等
		//   parser *gcmd.Parser  — 命令行参数解析器，可获取用户传入的命令行参数
		// 返回值:
		//   err error — Go 的错误类型。如果函数执行出错，返回非 nil 的 error
		//               框架会根据 error 决定返回给客户端的 HTTP 状态码
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// ============================================================
			// g.Server() — GoFrame 框架的服务器入口
			// 返回一个 HTTP 服务器实例，默认监听 :8000 端口（可在 config.yaml 中配置）
			// ============================================================
			s := g.Server()

			// ============================================================
			// s.Group("/", ...) — 创建根路由组
			// Group 方法用于创建路由分组，第一个参数是 URL 路径前缀。
			// "/" 表示根路径，所有路由都以 "/" 开头。
			// 第二个参数是回调函数，在其中注册子路由和中间件。
			//
			// GoFrame 路由组的特点：
			//   1. 路由组可以嵌套（group 里面再建 group）
			//   2. 中间件在路由组上设置，影响该组下所有路由
			//   3. 中间件按注册顺序执行（洋葱模型）
			// ============================================================
			s.Group("/", func(group *ghttp.RouterGroup) {
				// ============================================================
				// group.Middleware(...) — 注册全局中间件
				// 中间件是在请求到达实际处理函数之前/之后执行的函数。
				// 这里注册了两个中间件，按顺序执行：
				//
				// 1. middleware.CORS
				//    CORS（Cross-Origin Resource Sharing，跨域资源共享）中间件。
				//    由于前端（:3000）和后端（:8000）端口不同，
				//    浏览器会阻止跨域请求，CORS 中间件设置响应头来允许跨域。
				//
				// 2. ghttp.MiddlewareHandlerResponse
				//    GoFrame 标准响应处理中间件。
				//    它统一处理 API 响应格式：
				//    { "code": 0, "message": "success", "data": {...} }
				//    这个中间件会自动将控制器的返回值包装成标准格式返回给客户端。
				//    如果控制器返回 nil（空），中间件设置 204 状态码。
				//
				// 中间件执行顺序（洋葱模型）：
				//   请求到达 →
				//     CORS 中间件（处理跨域头）→
				//       MiddlewareHandlerResponse（准备响应包装）→
				//         实际路由处理 →
				//       MiddlewareHandlerResponse（包装响应）→
				//     CORS 中间件（添加响应头）→
				//   响应返回客户端
				// ============================================================
				group.Middleware(middleware.CORS, ghttp.MiddlewareHandlerResponse)

				// ============================================================
				// group.Group("/api/v1", ...) — 创建 /api/v1 子路由组
				// 所有 API 接口都挂在这个组下，URL 前缀为 /api/v1。
				// 这是常见的 API 版本控制方式。
				// ============================================================
				group.Group("/api/v1", func(v1 *ghttp.RouterGroup) {
					// ============================================================
					// 公共接口 — 不需要认证即可访问
					// ------------------------------------------------------------
					// 这些接口放在认证中间件之外，所以不需要 JWT 令牌。
					//
					// v1.POST("/login", controller.Auth, "Login")
					//   注册了一个 POST /api/v1/login 路由。
					//   参数说明：
					//     "/login"      — 相对于当前路由组前缀的路由路径
					//     controller.Auth — 控制器实例（在 controller 包中初始化）
					//     "Login"       — 控制器上的方法名（字符串形式）
					//   这种注册方式叫做"方法名映射"，GoFrame 通过反射调用对应方法。
					//   等价于 controller.Auth.Login(ctx, req) -> res, err
					//
					//   最终 URL: POST /api/v1/login
					// ============================================================
					v1.POST("/login", controller.Auth, "Login")
					v1.POST("/pageview/track", controller.PageView, "Track")
					v1.POST("/tools/mock-data", controller.Tools, "MockData")

					// ============================================================
					// v1.Group("", ...) — 需要认证的子路由组
					// 路径前缀为空字符串 ""，表示继承父路由组的 /api/v1 路径。
					// 这个组内的所有路由都需要经过认证中间件检查。
					// ============================================================
					v1.Group("", func(auth *ghttp.RouterGroup) {
						// auth.Middleware(middleware.Auth)
						// 在该路由组上添加认证中间件。
						// middleware.Auth 会：
						//   1. 从请求头中提取 JWT 令牌（Authorization: Bearer <token>）
						//   2. 验证令牌签名和有效期
						//   3. 解析出用户 ID 和角色列表
						//   4. 将用户信息注入到请求上下文中
						//   5. 如果令牌无效，直接返回 401 未授权错误
						//
						// 这个中间件对该组下所有路由生效（包括子路由组）。
						auth.Middleware(middleware.Auth)

						// 认证后的 Auth 接口
						auth.GET("/user/info", controller.Auth, "GetUserInfo")
						auth.POST("/logout", controller.Auth, "Logout")

						// 认证后的页面统计接口
						auth.GET("/pageview/stats", controller.PageView, "Stats")

						// ============================================================
						// auth.Bind(...) — 批量绑定控制器
						// Bind 方法是 GoFrame 提供的便捷注册方式。
						// 它会自动扫描控制器中所有带有 g.Meta 标签的结构体，
						// 并根据 meta 中的 path/method 自动注册路由。
						//
						// 这比手动写 GET/POST 更简洁，因为路由信息已经定义在结构体 tag 中。
						// 支持 Bind 的控制器：controller.User, controller.Role, controller.Dashboard
						//
						// Bind 与手动注册的区别：
						//   • 手动注册（POST GET PUT DELETE）— 显式写方法名
						//   • Bind 注册 — 反射自动扫描，代码更简洁
						//
						// 注意：Bind 注册的路由会继承当前路由组的中间件和前缀。
						// 所以 controller.User 的 /user 路由实际是 /api/v1/user。
						// ============================================================
					auth.Bind(
						controller.User,
						controller.Role,
						controller.Dashboard,
						controller.Menu,
					)

						// ============================================================
						// auth.Group("/user", ...) — 用户管理子路由组
						// 这个组额外添加了权限中间件：
						//   middleware.Permission("super_admin", "admin")
						// 表示只有 super_admin 或 admin 角色的用户才能访问该组下的接口。
						//
						// Permission 中间件是变参函数（variadic function），
						// 可以传多个角色名，匹配任一角色即可通过。
						// 中间件的执行顺序：
						//   请求 → CORS → MiddlewareHandlerResponse → Auth → Permission → 实际处理函数
						// ============================================================
						auth.Group("/user", func(user *ghttp.RouterGroup) {
							user.Middleware(middleware.Permission("super_admin", "admin"))
							// 手动注册需要额外权限控制的接口
							user.POST("", controller.User, "Create")
							user.PUT("/{id}/roles", controller.User, "AssignRoles")
						})

						// ============================================================
						// auth.Group("/role", ...) — 角色管理子路由组
						// 同样需要 super_admin 或 admin 权限。
						// 注意这里只添加了中间件，没有手动注册路由。
						// 路由由 controller.Role 的 Bind 自动注册（见上面的 auth.Bind）。
						//
						// 权限层级总结：
						//   /api/v1/login 等          — 无需认证，完全公开
						//   /api/v1/user/info 等     — 需要 JWT 认证（Auth 中间件）
						//   /api/v1/user 等           — 需要认证 + super_admin/admin 角色
						//   /api/v1/role 等           — 需要认证 + super_admin/admin 角色
						// ============================================================
						auth.Group("/role", func(role *ghttp.RouterGroup) {
							role.Middleware(middleware.Permission("super_admin", "admin"))
						})

						// Menu management (requires super_admin or admin)
						auth.Group("/menu", func(menu *ghttp.RouterGroup) {
							menu.Middleware(middleware.Permission("super_admin", "admin"))
						})

						// Current user menus (any authenticated user)
						auth.GET("/menu/user", controller.Menu, "GetUserMenus")
					})
				})
			})

			// s.Run() — 启动 HTTP 服务器，开始监听端口并处理请求
			// 这是一个阻塞调用，会一直运行直到收到中断信号。
			s.Run()

			// return nil — 函数正常结束，没有错误
			// nil 是 Go 的空值，表示"没有错误"
			return nil
		},
	}
)
