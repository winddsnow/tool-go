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
					v1.Bind(
						controller.Auth,
					)

					v1.Group("", func(auth *ghttp.RouterGroup) {
						auth.Middleware(middleware.Auth)

						auth.Bind(
							controller.User,
							controller.Role,
						)

						auth.Group("/user", func(user *ghttp.RouterGroup) {
							user.POST("", middleware.Permission("super_admin"), controller.User.Create)
						})
					})
				})
			})
			s.Run()
			return nil
		},
	}
)
