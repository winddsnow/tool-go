package controller

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
	"tool-go/internal/library/jwt"
	"tool-go/internal/middleware"
)

var PageView = cPageView{}

type cPageView struct{}

func (c *cPageView) Track(ctx context.Context, req *v1.PageViewTrackReq) (*v1.PageViewTrackRes, error) {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return &v1.PageViewTrackRes{}, nil
	}

	userId := middleware.GetUserId(ctx)
	username := middleware.GetUsername(ctx)

	// 未登录用户记作游客
	if userId == 0 {
		username = "游客"
		// 尝试从 Authorization 头解析用户信息（前端可能已登录但当前请求未走 Auth 中间件）
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			jwtCfg := g.Cfg().MustGet(ctx, "jwt").MapStrVar()
			secret := jwtCfg["secret"].String()
			if secret == "" {
				secret = "tool-go-jwt-secret-key-change-in-production"
			}
			j := jwt.New(secret, 0, "")
			claims, err := j.ParseToken(tokenStr)
			if err == nil && claims.UserId > 0 {
				userId = claims.UserId
				username = claims.Username
			}
		}
	}

	ip := r.GetClientIp()
	ua := r.Header.Get("User-Agent")
	if len(ua) > 512 {
		ua = ua[:512]
	}

	_, err := dao.PageView.Ctx(ctx).Data(g.Map{
		dao.PageView.Columns.PagePath:  req.PagePath,
		dao.PageView.Columns.UserId:    userId,
		dao.PageView.Columns.Username:  username,
		dao.PageView.Columns.IpAddress: ip,
		dao.PageView.Columns.UserAgent: ua,
	}).Insert()
	if err != nil {
		g.Log().Warning(ctx, "记录页面访问失败:", err)
	}

	return &v1.PageViewTrackRes{}, nil
}

func (c *cPageView) Stats(ctx context.Context, req *v1.PageViewStatsReq) (*v1.PageViewStatsRes, error) {
	totalVisits, err := dao.PageView.Ctx(ctx).Count()
	if err != nil {
		totalVisits = 0
	}

	type visitCount struct {
		Username string `json:"username"`
		Count    int    `json:"count"`
	}
	var userCounts []visitCount
	err = dao.PageView.Ctx(ctx).
		Fields("username, COUNT(*) as count").
		Where("username != ?", "").
		Group("username").
		OrderDesc("count").
		Limit(10).
		Scan(&userCounts)
	if err != nil {
		userCounts = nil
	}

	items := make([]v1.UserVisitItem, 0, len(userCounts))
	for _, v := range userCounts {
		if strings.TrimSpace(v.Username) == "" {
			continue
		}
		items = append(items, v1.UserVisitItem{
			Username: v.Username,
			Count:    v.Count,
		})
	}

	return &v1.PageViewStatsRes{
		TotalVisits: totalVisits,
		UserVisits:  items,
	}, nil
}
