package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Permission(requiredRole string) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		ctx := r.GetContext()
		if !HasRole(ctx, requiredRole) {
			r.Response.WriteStatus(http.StatusForbidden)
			r.Response.WriteJsonExit(g.Map{
				"code":    gcode.CodeNotAuthorized.Code(),
				"message": "没有权限访问",
			})
			return
		}
		r.Middleware.Next()
	}
}
