package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"tool-go/internal/controller"
	"tool-go/internal/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.CORS, ghttp.MiddlewareHandlerResponse)

				group.Group("/api/v1", func(v1 *ghttp.RouterGroup) {
					// 公共接口: 登录不需要认证
					v1.POST("/login", controller.Auth, "Login")

					v1.Group("", func(auth *ghttp.RouterGroup) {
						auth.Middleware(middleware.Auth)

						// 需要认证的 Auth 接口
						auth.GET("/user/info", controller.Auth, "GetUserInfo")
						auth.POST("/logout", controller.Auth, "Logout")

						auth.POST("/pageview/track", controller.PageView, "Track")
						auth.GET("/pageview/stats", controller.PageView, "Stats")

						auth.Bind(
							controller.User,
							controller.Role,
							controller.Dashboard,
						)

						auth.Group("/user", func(user *ghttp.RouterGroup) {
							user.Middleware(middleware.Permission("super_admin", "admin"))
							user.POST("", controller.User, "Create")
							user.PUT("/{id}/roles", controller.User, "AssignRoles")
						})

						auth.Group("/role", func(role *ghttp.RouterGroup) {
							role.Middleware(middleware.Permission("super_admin", "admin"))
						})
					})
				})
			})
			s.Run()
			return nil
		},
	}
)
