// Package controller 包含所有 HTTP 请求处理器。
// 本文件是页面访问追踪（PageView）相关的 Controller，用于记录和统计用户浏览行为。
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

// var PageView = cPageView{} 是 package-level 单例，全局只有一个实例。
var PageView = cPageView{}

// cPageView 是页面访问追踪 controller。
type cPageView struct{}

// Track 记录一次页面访问。这个接口被前端路由守卫（router.beforeEach）自动调用，
// 每次页面切换都会发送请求到此接口记录访问日志。
//
// 客人检测机制：
//   - 通常 API 请求会经过 Auth 中间件，中间件解析 JWT 后将 userId 和 username 写入 context。
//   - 但页面访问追踪在登录页、公开页面也需记录，这些请求不会经过 Auth 中间件。
//   - 所以：先尝试从 context 获取用户（如果没走中间件则为 0），
//     再手动解析 Authorization 头中的 JWT（如果前端已登录但当前路由未走中间件）。
func (c *cPageView) Track(ctx context.Context, req *v1.PageViewTrackReq) (*v1.PageViewTrackRes, error) {
	// g.RequestFromCtx(ctx) 从 context 中恢复出 GoFrame 的 Request 对象。
	// 正常请求一定能取到，但某些测试场景或异步任务可能没有关联的请求。
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return &v1.PageViewTrackRes{}, nil
	}

	userId := middleware.GetUserId(ctx)
	username := middleware.GetUsername(ctx)

	// 未登录用户记作游客
	if userId == 0 {
		username = "游客"
		// 尝试从 Authorization 头解析用户信息（前端可能已登录但当前请求未走 Auth 中间件）。
		// Auth 中间件只挂载在需要认证的路由上，但 pageview 追踪接口对所有用户开放，
		// 所以需要在 controller 层自行判断。
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			// strings.TrimPrefix 去掉 "Bearer " 前缀得到 token 字符串。
			// Go 的 strings.HasPrefix / TrimPrefix 是高效的字符串操作（不分配新内存）。
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			jwtCfg := g.Cfg().MustGet(ctx, "jwt").MapStrVar()
			secret := jwtCfg["secret"].String()
			if secret == "" {
				secret = "tool-go-jwt-secret-key-change-in-production"
			}
			// jwt.New(secret, 0, "")：第二个参数 expires=0 表示"不需要过期验证"，
			// 因为我们只需要解析已有的 JWT 获取用户信息，而不是生成新 token。
			j := jwt.New(secret, 0, "")
			claims, err := j.ParseToken(tokenStr)
			if err == nil && claims.UserId > 0 {
				userId = claims.UserId
				username = claims.Username
			}
		}
	}

	// r.GetClientIp() 获取客户端 IP（支持 X-Forwarded-For 等代理头）。
	// r.Header.Get("User-Agent") 获取 User-Agent 请求头。
	// Go 的 map/slice 操作：len(ua) 返回字符串字节长度（非字符数）。
	// ua[:512] 是切片操作，截取前 512 个字节。
	// 为什么要限制 512 字节？User-Agent 可能非常长（某些浏览器发送的 UA 可达几千字节），
	// 数据库中该列一般设为 VARCHAR(512)，截断可避免数据库插入错误。
	ip := r.GetClientIp()
	ua := r.Header.Get("User-Agent")
	if len(ua) > 512 {
		ua = ua[:512]
	}

	// dao.PageView.Ctx(ctx) 获取 ORM 模型。
	// Data(g.Map{...}) 设置插入字段，g.Map 是 GoFrame 的泛型 map（map[string]interface{}）。
	// dao.PageView.Columns.PagePath 是 ORM 自动生成的列名常量（类型安全，拼写错误在编译期即可发现）。
	// Insert() 执行 INSERT INTO 并返回结果。
	_, err := dao.PageView.Ctx(ctx).Data(g.Map{
		dao.PageView.Columns.PagePath:  req.PagePath,
		dao.PageView.Columns.UserId:    userId,
		dao.PageView.Columns.Username:  username,
		dao.PageView.Columns.IpAddress: ip,
		dao.PageView.Columns.UserAgent: ua,
	}).Insert()
	if err != nil {
		// g.Log().Warning 是 GoFrame 的日志组件，级别为 Warning（低于 Error）。
		// 使用 Warning 而非 Error，因为页面访问记录是一个辅助功能，
		// 偶尔失败不应该影响用户体验。前端不需要知道记录失败。
		// 日志会自动包含请求追踪 ID、时间戳等信息。
		g.Log().Warning(ctx, "记录页面访问失败:", err)
	}

	return &v1.PageViewTrackRes{}, nil
}

// Stats 获取页面访问统计：总访问数和用户访问排行（Top 10）。
// 逻辑和 dashboard.go 的 GetStats 类似，但这里作为一个独立的 API 暴露，
// 供"访问统计"页面使用（而不是放在仪表盘里）。
//
// 错误处理说明：
//   - err = dao.PageView.Ctx(ctx).Count()：与 dashboard.go 不同，
//     这里显式检查了错误（而非用 _ 忽略当"零值"），但返回值相同：出错时 totalVisits 保持 0。
//   - 两种风格在 Go 中都常见："显式检查"更规范，"忽略"更简洁。
func (c *cPageView) Stats(ctx context.Context, req *v1.PageViewStatsReq) (*v1.PageViewStatsRes, error) {
	totalVisits, err := dao.PageView.Ctx(ctx).Count()
	if err != nil {
		totalVisits = 0
	}

	// type visitCount struct 是函数内局部类型。在 Go 中，可以在任意代码块中定义类型。
	// 当某类型只在当前函数使用，就定义在函数内部，避免污染包命名空间。
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

	// make([]v1.UserVisitItem, 0, len(userCounts)) 预分配切片，
	// 第三个参数是容量（capacity），避免 append 时动态扩容的性能开销。
	// Go 的切片扩容策略：当容量不足时，一般扩容为原容量的 2 倍（元素数 < 1024）或 1.25 倍。
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
