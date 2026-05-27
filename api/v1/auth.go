// ============================================================
// package v1 — API 版本包
// ------------------------------------------------------------
// Go 语言中，package（包）是代码组织的基本单位。
// 一个包由同一目录下所有 .go 文件组成，package 声明必须在文件第一行。
// v1 表示 API 的第一版本。通过将不同版本的 API 放在不同包中
// （如 v1、v2），可以在同一项目中同时维护多个 API 版本，
// 方便向后兼容。其他包通过 import "tool-go/api/v1" 来使用这里定义的类型。
//
// 本文件定义了"认证（Auth）"相关的 HTTP 请求和响应数据结构。
// 按照 GoFrame 框架的约定，API 层只定义数据结构（struct），
// 不写业务逻辑。每个 API 接口对应一对结构体：
//   • XxxReq  — 请求参数结构体，框架会自动解析 HTTP 请求参数到该结构体
//   • XxxRes  — 响应数据结构体，框架会自动将该结构体序列化为 JSON 返回
// ============================================================
package v1

// import 关键字用于导入其他包。
// "github.com/gogf/gf/v2/frame/g" 是 GoFrame 框架的核心包。
// GoFrame 的所有核心功能（路由、数据库、缓存等）都通过这个包下的接口使用。
// 在这里我们使用 g.Meta 结构体来嵌入路由元数据。
import "github.com/gogf/gf/v2/frame/g"

// ============================================================
// LoginReq — 用户登录请求参数
// ------------------------------------------------------------
// struct 是 Go 中定义结构体的关键字，类似于其他语言的 class，
// 但 struct 只包含字段（数据），不包含方法（行为）。
//
// 字段说明:
//   g.Meta — GoFrame 框架提供的特殊结构体，嵌入到请求结构体中后，
//            框架会解析其 struct tag（反引号中的内容）来获取路由信息。
//            这是 GoFrame 的特色设计：通过 struct tag 声明路由，
//            而不是像传统框架那样在路由注册处写 URL 字符串。
//            这样做的好处是路由和参数定义在一起，代码更内聚。
//
// Struct tag 详解:
//   `path:"/login"`      — 路由路径，HTTP 请求的 URL 路径
//   `method:"post"`      — HTTP 请求方法，这里是 POST
//   `tags:"Auth"`         — OpenAPI/Swagger 文档的分类标签
//   `summary:"用户登录"`   — OpenAPI/Swagger 文档的接口摘要
//
//   Username 字段:
//     string — Go 内置的字符串类型
//     `json:"username"`   — JSON 序列化时的字段名。
//                           当该结构体被转为 JSON 字符串时，字段名变为 username。
//                           Go 的 json 包默认使用字段名作为 JSON 键名，
//                           但 Go 约定导出字段首字母大写，而前端通常用小写，
//                           所以需要用 json tag 来指定 JSON 中的名称。
//
//     `v:"required#用户名不能为空"` — GoFrame 验证规则（validation tag）。
//          required 表示该字段必填。`#` 是分隔符，`#` 后面是验证失败时的
//          自定义错误信息。如果不写 `#`，框架使用默认错误信息。
//          验证在请求到达控制器之前自动执行，开发者无需手动写验证代码。
//
//     `dc:"用户名"`        — description comment，OpenAPI 文档中对字段的描述。
//                           dc 是 description comment 的缩写。
//                           框架会自动根据 dc 生成 Swagger 文档。
//
//   Password 字段: 同上模式，定义密码字段
// ============================================================
type LoginReq struct {
	g.Meta   `path:"/login" method:"post" tags:"Auth" summary:"用户登录"`
	Username string `json:"username" v:"required#用户名不能为空" dc:"用户名"`
	Password string `json:"password" v:"required#密码不能为空" dc:"密码"`
}

// ============================================================
// LoginRes — 用户登录响应数据
// ------------------------------------------------------------
// 响应结构体不需要 g.Meta，因为响应不需要路由信息。
// 框架会自动将此类结构体序列化为 JSON 格式返回给客户端。
// 如果返回的是 nil 或空结构体，框架会返回 204 No Content。
//
// 字段说明:
//   AccessToken  string   — JWT（JSON Web Token）访问令牌。
//                       JWT 是一种用于身份验证的令牌格式，
//                       客户端在后续请求的 Authorization 头中携带此令牌。
//                       []string 表示字符串切片（类似于其他语言的数组）。
//
//   UserId   uint64   — 用户 ID。uint64 是无符号 64 位整数，
//                       范围 0 到 18446744073709551615。
//                       因为用户 ID 不可能是负数，所以使用无符号整数。
//
//   Nickname string   — 用户昵称，用于显示而非用户名。
//
//   Roles    []string — 角色列表，每个元素是角色编码（如 "admin"）。
//                       []string 是"字符串切片"，Go 中的动态数组。
// ============================================================
type LoginRes struct {
	AccessToken  string     `json:"access_token" dc:"访问令牌"`
	UserId   uint64     `json:"user_id" dc:"用户ID"`
	Username string     `json:"username" dc:"用户名"`
	Nickname string     `json:"nickname" dc:"昵称"`
	Roles       []string   `json:"roles" dc:"角色列表"`
	Menus       []MenuTree `json:"menus" dc:"菜单树"`
	Permissions []string   `json:"permissions" dc:"权限码列表"`
}

type RefreshReq struct {
	g.Meta `path:"/refresh" method:"post" tags:"Auth" summary:"刷新访问令牌"`
}

type RefreshRes struct {
	AccessToken string `json:"access_token" dc:"新的访问令牌"`
}

// ============================================================
// GetUserInfoReq — 获取当前用户信息请求参数
// ------------------------------------------------------------
// 这个请求不需要任何参数，因为当前用户的身份由 JWT 令牌确定。
// 令牌在 Authorization 请求头中传递，由中间件解析。
// 所以结构体中除了 g.Meta 没有其他字段。
//
// method:"get" — HTTP GET 请求。
// GET 请求通常用于获取数据，请求参数通过 URL 查询字符串传递。
// ============================================================
type GetUserInfoReq struct {
	g.Meta `path:"/user/info" method:"get" tags:"Auth" summary:"获取当前用户信息"`
}

// GetUserInfoRes — 用户信息响应
// 与 LoginRes 的结构相同，但没有 Token 字段。
// 这在 GoFrame 中很常见：不同的 API 可以返回不同的字段，
// 而不是用一个通用结构体。这样每个接口的响应都是精确的，
// 不会返回多余字段。
type GetUserInfoRes struct {
	UserId      uint64     `json:"user_id" dc:"用户ID"`
	Username    string     `json:"username" dc:"用户名"`
	Nickname    string     `json:"nickname" dc:"昵称"`
	Roles       []string   `json:"roles" dc:"角色列表"`
	Menus       []MenuTree `json:"menus" dc:"菜单树"`
	Permissions []string   `json:"permissions" dc:"权限码列表"`
}

// ============================================================
// LogoutReq — 退出登录请求参数
// ------------------------------------------------------------
// 同样不需要额外参数。JWT 无状态，所以真正的"退出"通常在客户端
// 删除令牌。后端可以选择将令牌加入黑名单。
// ============================================================
type LogoutReq struct {
	g.Meta `path:"/logout" method:"post" tags:"Auth" summary:"退出登录"`
}

// LogoutRes — 退出登录响应
// 空结构体表示成功时返回 204 No Content（无内容）。
// 所有字段都为空，框架只返回 HTTP 状态码。
type LogoutRes struct{}
