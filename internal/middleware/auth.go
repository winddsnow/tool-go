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

const (
	CtxUserId   = "userId"
	CtxUsername = "username"
	CtxRoles    = "roles"
)

func Auth(r *ghttp.Request) {
	ctx := r.GetCtx()
	jwtConfig := g.Cfg().MustGet(ctx, "jwt").MapStrStr()
	secret := jwtConfig["secret"]
	if secret == "" {
		secret = "tool-go-jwt-secret-key-change-in-production"
	}

	j := jwt.New(secret, 0, "")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.WriteJsonExit(g.Map{
			"code":    gcode.CodeNotAuthorized.Code(),
			"message": "未登录或登录已过期",
		})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.WriteJsonExit(g.Map{
			"code":    gcode.CodeNotAuthorized.Code(),
			"message": "token格式错误",
		})
		return
	}

	claims, err := j.ParseToken(tokenString)
	if err != nil {
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.WriteJsonExit(g.Map{
			"code":    gcode.CodeNotAuthorized.Code(),
			"message": gstr.TrimLeftStr(err.Error(), "jwt: "),
		})
		return
	}

	ctx = context.WithValue(ctx, CtxUserId, claims.UserId)
	ctx = context.WithValue(ctx, CtxUsername, claims.Username)
	ctx = context.WithValue(ctx, CtxRoles, claims.Roles)
	r.SetCtx(ctx)

	r.Middleware.Next()
}

func GetUserId(ctx context.Context) uint64 {
	if v := ctx.Value(CtxUserId); v != nil {
		return v.(uint64)
	}
	return 0
}

func GetUsername(ctx context.Context) string {
	if v := ctx.Value(CtxUsername); v != nil {
		return v.(string)
	}
	return ""
}

func GetRoles(ctx context.Context) []string {
	if v := ctx.Value(CtxRoles); v != nil {
		return v.([]string)
	}
	return []string{}
}

func HasRole(ctx context.Context, role string) bool {
	roles := GetRoles(ctx)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func HasAnyRole(ctx context.Context, roles ...string) bool {
	for _, role := range roles {
		if HasRole(ctx, role) {
			return true
		}
	}
	return false
}
