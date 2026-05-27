// Package controller 包含所有 HTTP 请求处理器。
// 在 Go 中，一个 package（包）是代码组织的基本单位，由同一目录下的 .go 文件组成（且必须使用相同的 package 声明）。
// 外部通过 "包名.导出名" 来调用（如 controller.Auth.Login()）。
// 本包实现 GoFrame 框架的标准 Controller 模式：每个方法接收 context + 请求对象，返回响应对象 + error。
package controller

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
	"tool-go/internal/library/jwt"
	"tool-go/internal/library/password"
	"tool-go/internal/middleware"
	"tool-go/internal/model/entity"
	"tool-go/internal/service"
)

// var Auth = cAuth{} 是 Go 中典型的 package-level 单例模式。
// Go 的 var 声明在包初始化时执行一次，cAuth{} 创建了一个空结构体实例赋值给 Auth。
// 由于 cAuth 没有字段（空结构体），所有实例共享同一内存地址，且协程安全（无共享状态需要保护）。
// 这种方式避免了每次请求都 new 一个 controller，也无需依赖注入框架。
var Auth = cAuth{}

// cAuth 是 auth 相关的 controller 结构体。
// 在 Go 中，结构体（struct）是字段和方法的集合，类似其他语言的 class。
// 这里定义为私有（小写字母开头 cAuth），外部只能通过 var Auth 访问。
type cAuth struct{}

// Login 处理用户登录请求，是 GoFrame Controller 的标准方法签名：
//   - ctx context.Context：Go 的上下文接口，用于传递请求范围的值（如用户ID、超时控制、链路追踪等）。
//     每个请求都会创建一个 context，贯穿整个请求生命周期。
//   - req *v1.LoginReq：指向请求结构体的指针（*T 表示指向 T 类型的指针）。
//     按 GoFrame 规范，请求/响应定义在 api/v1 中，包含请求参数、验证规则等。
//   - *v1.LoginRes：返回指向响应结构体的指针，框架会自动序列化为 JSON。
//   - error：Go 内置错误接口。nil 表示成功，非 nil 表示出错并自动返回错误响应。
//
// 流程：查用户 -> 验证密码 -> 检查状态 -> 查角色 -> 生成JWT -> 返回token和用户信息。
func (c *cAuth) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRes, error) {
	// var user *entity.User 声明一个指向 entity.User 的指针变量，初始值为 nil。
	// 指针（*Type）在 Go 中表示"指向某类型值的引用"，允许修改原值且避免拷贝。
	var user *entity.User
	// dao 是 GoFrame 的 Data Access Object 层，封装了数据库操作。
	// Ctx(ctx) 将 context 注入 ORM 操作（用于超时控制、日志追踪）。
	// Where(列名, 值) 生成 WHERE 子句。
	// WhereNull(列名) 生成 WHERE 列名 IS NULL，确保 soft-delete 的用户不被查出。
	// Scan(&user) 将查询结果映射到 user 指针指向的变量（解引用）。
	err := dao.User.Ctx(ctx).
		Where(dao.User.Columns.Username, req.Username).
		WhereNull(dao.User.Columns.DeletedAt).
		Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// gerror.New 创建 GoFrame 的错误实例，支持堆栈跟踪和国际化。
		// 错误返回时框架自动将 message 写入 HTTP 响应体。
		return nil, gerror.New("用户名或密码错误")
	}

	// 验证密码：password.VerifyPassword 用存储的 salt+hash 验证明文密码
	if !password.VerifyPassword(req.Password, user.Salt, user.Password) {
		return nil, gerror.New("用户名或密码错误")
	}

	// Status == 0 表示用户被禁用
	if user.Status == 0 {
		return nil, gerror.New("账号已被禁用")
	}

	// 调用辅助函数获取该用户的角色列表
	roles := getUserRoles(ctx, user.Id)
	permissions := getUserPermissions(ctx, user.Id)

	// g.Cfg().MustGet(ctx, "jwt") 从 config.yaml 读取 jwt 配置段。
	// MustGet 在配置不存在时 panic（自动输出错误信息），在开发阶段便于排查。
	// MapStrVar() 转为 map[string]g.Var，方便用 .String() / .Duration() 取值。
	// GoFrame 的配置管理支持多数据源（文件、etcd、数据库等），默认为 YAML 文件。
	jwtConfig := g.Cfg().MustGet(ctx, "jwt").MapStrVar()
	secret := jwtConfig["secret"].String()
	if secret == "" {
		secret = "tool-go-jwt-secret-key-change-in-production"
	}
	expires := jwtConfig["expires"].Duration()
	if expires == 0 {
		expires = 24 * time.Hour
	}
	issuer := jwtConfig["issuer"].String()
	if issuer == "" {
		issuer = "tool-go"
	}

	// jwt.New 创建 JWT 工具实例，GenerateToken 生成 token 字符串。
	j := jwt.New(secret, expires, issuer)
	token, err := j.GenerateToken(user.Id, user.Username, roles, permissions)
	if err != nil {
		return nil, gerror.New("生成token失败")
	}

	refreshToken, err := j.GenerateRefreshToken(user.Id)
	if err != nil {
		return nil, gerror.New("生成refresh token失败")
	}

	r := ghttp.RequestFromCtx(ctx)
	r.Response.Cookie().Set(ghttp.CookieOption{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1",
		HttpOnly: true,
		MaxAge:   604800,
	})

	// &v1.LoginRes{...} 创建结构体指针并初始化字段（冒号分隔字段名和值），最后返回（nil 错误表示成功）
	menuRes, _ := service.Menu().GetUserMenus(ctx, user.Id)
	var menus []v1.MenuTree
	if menuRes != nil {
		menus = menuRes.Menus
	}
	return &v1.LoginRes{
		AccessToken: token,
		UserId:      user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Roles:       roles,
		Menus:       menus,
		Permissions: permissions,
	}, nil
}

// GetUserInfo 获取当前登录用户的个人信息。
// middleware.GetUserId(ctx) 从 context 中提取 JWT 中间件存入的 userId。
// 这利用了 Go 的 context.WithValue：middleware 在认证时将 userId 写入 context，
// 后续 controller 用对应的 GetUserId 取出，无需重复解析 JWT。
func (c *cAuth) GetUserInfo(ctx context.Context, req *v1.GetUserInfoReq) (*v1.GetUserInfoRes, error) {
	userId := middleware.GetUserId(ctx)
	if userId == 0 {
		return nil, gerror.New("未登录")
	}

	var user *entity.User
	err := dao.User.Ctx(ctx).
		Where(dao.User.Columns.Id, userId).
		WhereNull(dao.User.Columns.DeletedAt).
		Scan(&user)
	if err != nil || user == nil {
		return nil, gerror.New("用户不存在")
	}

	roles := getUserRoles(ctx, user.Id)
	permissions := getUserPermissions(ctx, user.Id)

	menuRes, _ := service.Menu().GetUserMenus(ctx, user.Id)
	var menus []v1.MenuTree
	if menuRes != nil {
		menus = menuRes.Menus
	}
	return &v1.GetUserInfoRes{
		UserId:      user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Roles:       roles,
		Menus:       menus,
		Permissions: permissions,
	}, nil
}

// Logout 处理登出请求，清除 refresh token cookie。
func (c *cAuth) Logout(ctx context.Context, req *v1.LogoutReq) (*v1.LogoutRes, error) {
	r := ghttp.RequestFromCtx(ctx)
	r.Response.Cookie().Set(ghttp.CookieOption{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1",
		HttpOnly: true,
		MaxAge:   -1,
	})
	return &v1.LogoutRes{}, nil
}

func (c *cAuth) Refresh(ctx context.Context, req *v1.RefreshReq) (*v1.RefreshRes, error) {
	r := ghttp.RequestFromCtx(ctx)
	refreshToken := r.Cookie("refresh_token")
	if refreshToken == "" {
		return nil, gerror.New("refresh token不存在")
	}

	jwtConfig := g.Cfg().MustGet(ctx, "jwt").MapStrVar()
	secret := jwtConfig["secret"].String()
	if secret == "" {
		secret = "tool-go-jwt-secret-key-change-in-production"
	}
	expires := jwtConfig["expires"].Duration()
	if expires == 0 {
		expires = 15 * time.Minute
	}
	issuer := jwtConfig["issuer"].String()
	if issuer == "" {
		issuer = "tool-go"
	}

	j := jwt.New(secret, expires, issuer)
	claims, err := j.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, gerror.New("refresh token无效或已过期")
	}

	var user *entity.User
	err = dao.User.Ctx(ctx).
		Where(dao.User.Columns.Id, claims.UserId).
		WhereNull(dao.User.Columns.DeletedAt).
		Scan(&user)
	if err != nil || user == nil {
		return nil, gerror.New("用户不存在")
	}

	roles := getUserRoles(ctx, user.Id)
	permissions := getUserPermissions(ctx, user.Id)

	newAccessToken, err := j.GenerateToken(user.Id, user.Username, roles, permissions)
	if err != nil {
		return nil, gerror.New("生成token失败")
	}

	newRefreshToken, err := j.GenerateRefreshToken(user.Id)
	if err != nil {
		return nil, gerror.New("生成refresh token失败")
	}

	r.Response.Cookie().Set(ghttp.CookieOption{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/api/v1",
		HttpOnly: true,
		MaxAge:   604800,
	})

	return &v1.RefreshRes{
		AccessToken: newAccessToken,
	}, nil
}

// getUserRoles 是包内私有的辅助函数（小写字母开头，外部不可见）。
// uint64 是 Go 的无符号 64 位整数类型，用于存储用户 ID。
// 返回值 []string 是字符串切片（动态数组），Go 中 slice 是引用类型，传递时复制引用头而非全部数据。
//
// LEFT JOIN + Array() 说明：
//   - LeftJoin("user_role", "user_role.role_id=role.id")：SQL LEFT JOIN，
//     将 role 表与关联表 user_role 连接，即使 user_role 没有匹配行，role 的行也会保留。
//     这里实际上希望的是 INNER JOIN 效果，但 LEFT JOIN 也能正确工作。
//   - Fields(列名).Array()：只查询单列，Array() 将结果集第一列返回为 []g.Var 切片。
//     比 Scan 到结构体数组更轻量，适合只需要一列值的场景。
//   - g.Log().Error()：GoFrame 的日志组件，自动记录错误栈和请求追踪 ID。
func getUserRoles(ctx context.Context, userId uint64) []string {
	result, err := dao.Role.Ctx(ctx).
		LeftJoin("user_role", "user_role.role_id=role.id").
		Where("user_role.user_id", userId).
		WhereNull(dao.Role.Columns.DeletedAt).
		Fields(dao.Role.Columns.Code).
		Array()
	if err != nil {
		g.Log().Error(ctx, "获取用户角色失败:", err)
		// return []string{"user"} 创建并返回一个包含字符串 "user" 的切片字面量。
		// 即使查询出错也返回默认角色，使系统不会完全不可用（降级策略）。
		return []string{"user"}
	}
	if len(result) == 0 {
		return []string{"user"}
	}
	// make([]string, len(result)) 预分配切片（长度等于 result 数量），避免 append 动态扩容。
	// v.String() 将 g.Var 类型的数据库值转为 Go string。
	// 循环用 for i, v := range result { ... }，i 是索引（0-based），v 是元素值副本。
	codes := make([]string, len(result))
	for i, v := range result {
		codes[i] = v.String()
	}
	return codes
}

func getUserPermissions(ctx context.Context, userId uint64) []string {
	result, err := dao.Permission.Ctx(ctx).
		LeftJoin("role_permission", "role_permission.permission_id=permission.id").
		LeftJoin("user_role", "user_role.role_id=role_permission.role_id").
		Where("user_role.user_id", userId).
		Fields(dao.Permission.Columns.Code).
		Array()
	if err != nil {
		g.Log().Error(ctx, "获取用户权限失败:", err)
		return []string{}
	}
	codes := make([]string, len(result))
	for i, v := range result {
		codes[i] = v.String()
	}
	return codes
}
