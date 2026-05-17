// ============================================================
// package v1 — API v1 版本的用户管理接口数据结构
// ------------------------------------------------------------
// 本文件定义了"用户 CRUD（增删改查）"相关的 HTTP 请求/响应结构体。
// CRUD 是 Create（创建）、Read（查询）、Update（更新）、Delete（删除）
// 的缩写，是 RESTful API 设计的基本操作模式。
//
// 在 GoFrame 中，API 结构体同时承担了"请求参数定义"、"路由声明"、
// "请求验证规则"、"OpenAPI 文档生成"四个职责。这是通过 struct tag
// （结构体标签）实现的，是 GoFrame 的核心设计理念之一。
// ============================================================
package v1

// 导入 GoFrame 框架核心包，用于 g.Meta 路由元数据嵌入。
import "github.com/gogf/gf/v2/frame/g"

// ============================================================
// UserCreateReq — 创建用户请求参数
// ------------------------------------------------------------
// POST 请求通常用于创建资源。
// 注意这里 Status 字段的类型是 uint（无符号整数），
// 因为创建时我们指定了默认值 d:"1"（d 是 default 的缩写），
// 即默认状态为启用，所以即使用户不传这个字段，框架也会设为 1。
// 没有零值问题，所以这里可以用 uint。
//
// v:"email#邮箱格式不正确" — 内置 email 验证规则，自动检查邮箱格式。
//                            如果字符串不符合 email 格式，返回自定义错误。
// ============================================================
type UserCreateReq struct {
	g.Meta   `path:"/user" method:"post" tags:"User" summary:"创建用户"`
	Username string `json:"username" v:"required#用户名不能为空" dc:"用户名"`
	Password string `json:"password" v:"required#密码不能为空" dc:"密码"`
	Nickname string `json:"nickname" dc:"昵称"`
	Email    string `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`
	Phone    string `json:"phone" dc:"手机号"`
	Status   uint   `json:"status" d:"1" dc:"状态: 1=启用, 0=禁用"`
}

// UserCreateRes — 创建用户成功后的响应，返回新用户的 ID
type UserCreateRes struct {
	Id uint64 `json:"id" dc:"用户ID"`
}

// ============================================================
// UserDeleteReq — 删除用户请求参数
// ------------------------------------------------------------
// DELETE 请求用于删除资源。
// path:"/user/{id}" — {id} 是路径参数，在 URL 中以 /user/123 形式出现。
// `path:"id"` 标签告诉框架：Id 字段的值来自 URL 路径中的 {id} 参数，
// 而不是来自 JSON 请求体或 URL 查询参数。
//
// 相比 json tag、path tag 从 URL 路径中提取参数。
// ============================================================
type UserDeleteReq struct {
	g.Meta `path:"/user/{id}" method:"delete" tags:"User" summary:"删除用户"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"用户ID"`
}

// UserDeleteRes — 删除成功响应（空，返回 204）
type UserDeleteRes struct{}

// ============================================================
// UserUpdateReq — 更新用户请求参数
// ------------------------------------------------------------
// PUT 请求通常用于完整更新资源。
// 注意这里 Id 字段使用 `path:"id"` 从 URL 中获取，
// 而其他字段使用 `json:"..."` 从请求体中获取。
// 同一个结构体可以从不同位置提取参数。
// ============================================================
type UserUpdateReq struct {
	g.Meta   `path:"/user/{id}" method:"put" tags:"User" summary:"更新用户"`
	Id       uint64 `path:"id" v:"required#ID不能为空" dc:"用户ID"`
	Username string `json:"username" dc:"用户名"`
	Nickname string `json:"nickname" dc:"昵称"`
	Email    string `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`
	Phone    string `json:"phone" dc:"手机号"`
	Status   uint   `json:"status" dc:"状态: 1=启用, 0=禁用"`
}

// UserUpdateRes — 更新成功响应（空）
type UserUpdateRes struct{}

// ============================================================
// UserGetOneReq — 获取单个用户详情请求参数
// ------------------------------------------------------------
// GET 请求用于获取数据，这里通过 URL 路径参数 {id} 指定用户。
// ============================================================
type UserGetOneReq struct {
	g.Meta `path:"/user/{id}" method:"get" tags:"User" summary:"获取用户详情"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"用户ID"`
}

// UserGetOneRes — 用户详情响应
// datetime 格式的字符串直接返回给前端，由前端进行格式化显示。
// CreatedAt/UpdatedAt 是数据库的约定字段名，通常由 ORM 自动维护。
type UserGetOneRes struct {
	Id        uint64 `json:"id" dc:"用户ID"`
	Username  string `json:"username" dc:"用户名"`
	Nickname  string `json:"nickname" dc:"昵称"`
	Email     string `json:"email" dc:"邮箱"`
	Phone     string `json:"phone" dc:"手机号"`
	Status    uint   `json:"status" dc:"状态"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
	UpdatedAt string `json:"updated_at" dc:"更新时间"`
}

// ============================================================
// UserListReq — 获取用户列表请求参数（分页 + 筛选）
// ------------------------------------------------------------
// 这是需要特别关注的结构体，因为它展示了 Go Frame 的重要设计模式。
//
// 分页参数：
//   Page     int    — 当前页码，默认值 1（d:"1"）
//   PageSize int    — 每页显示数量，默认值 10（d:"10"）
//
// 筛选参数：
//   Username string — 按用户名模糊搜索
//
// ★ 关键设计：Status 字段使用 *int 而非 uint 的原因
// ------------------------------------------------------------
// *int 是"指向 int 的指针"，指针类型在没有赋值时为 nil（空指针）。
// 而 int/uint 的零值（默认值）是 0。
//
// 问题场景：如果用 uint Status，前端传 status=0（禁用状态）时，
// Go 无法区分"用户传了 0"和"用户没传这个字段"（都是零值 0）。
// 使用 *int 后：
//   • 用户没传 status → Status 为 nil（空指针，不触发筛选）
//   • 用户传 status=0 → Status 指向 0（触发筛选，查禁用的用户）
//   • 用户传 status=1 → Status 指向 1（触发筛选，查启用的用户）
//
// 这是一个常见的 Go 指针用法惯用模式，用于区分"传了零值"和"没传"。
//
// UserId   uint64  — JSON 标签为 json:"-"。
//                    "-" 表示该字段不参与 JSON 序列化/反序列化。
//                    这个值由认证中间件自动注入（从 JWT 中解析出的用户 ID），
//                    用户无法通过 HTTP 请求修改它。
//                    这是一种安全设计，防止用户非法查询其他用户的数据。
// ============================================================
type UserListReq struct {
	g.Meta   `path:"/user" method:"get" tags:"User" summary:"获取用户列表"`
	Page     int    `json:"page" d:"1" dc:"页码"`
	PageSize int    `json:"page_size" d:"10" dc:"每页数量"`
	Username string `json:"username" dc:"按用户名筛选"`
	Status   *int   `json:"status" dc:"按状态筛选"`
	UserId   uint64 `json:"-" dc:"当前用户ID"`
}

// ============================================================
// UserListRes — 用户列表响应
// ------------------------------------------------------------
// 响应包含三个部分：
//   List  []UserItem — 当前页的用户数据列表（切片/动态数组）
//   Total int        — 符合条件的用户总数（用于前端计算总页数）
//   Page  int        — 当前页码（前端用于回显）
//
// 这种 List + Total + Page 的分页模式是 RESTful API 的最佳实践之一。
// ============================================================
type UserListRes struct {
	List  []UserItem `json:"list" dc:"用户列表"`
	Total int        `json:"total" dc:"总数"`
	Page  int        `json:"page" dc:"当前页"`
}

// ============================================================
// UserItem — 列表中的单个用户项
// ------------------------------------------------------------
// 注意 UserItem 与 UserGetOneRes 不同：
//   UserItem 用于列表展示，字段较少
//   UserGetOneRes 用于详情页，字段更多（有 UpdatedAt）
//
// 这是一种设计模式：为不同场景定义不同的结构体，
// 而不是复用同一个。好处是：
//   1. 列表查询只需查较少字段（性能更好）
//   2. 避免返回不需要的数据（安全性更好）
//   3. 列表和详情可以独立演进而不会相互影响
// ============================================================
type UserItem struct {
	Id        uint64 `json:"id" dc:"用户ID"`
	Username  string `json:"username" dc:"用户名"`
	Nickname  string `json:"nickname" dc:"昵称"`
	Email     string `json:"email" dc:"邮箱"`
	Phone     string `json:"phone" dc:"手机号"`
	Status    uint   `json:"status" dc:"状态"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
}

// ============================================================
// UserGetRolesReq — 获取用户当前拥有的角色
// ------------------------------------------------------------
// 这是多对多关系查询的典型场景。
// 路径 /user/{id}/roles 表示"获取某个用户的角色列表"，
// 这是 RESTful API 的常见设计模式。
// ============================================================
type UserGetRolesReq struct {
	g.Meta `path:"/user/{id}/roles" method:"get" tags:"User" summary:"获取用户角色"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"用户ID"`
}

// UserGetRolesRes — 返回用户已分配的角色 ID 列表
type UserGetRolesRes struct {
	RoleIds []uint64 `json:"role_ids" dc:"角色ID列表"`
}

// ============================================================
// UserAssignRolesReq — 分配角色给用户请求
// ------------------------------------------------------------
// PUT /user/{id}/roles — 替换用户的所有角色（全量更新）。
// RoleIds 是 uint64 切片（数组），v:"required" 确保至少传一个角色。
// ============================================================
type UserAssignRolesReq struct {
	g.Meta  `path:"/user/{id}/roles" method:"put" tags:"User" summary:"分配用户角色"`
	Id      uint64   `path:"id" v:"required#ID不能为空" dc:"用户ID"`
	RoleIds []uint64 `json:"role_ids" v:"required#角色列表不能为空" dc:"角色ID列表"`
}

// UserAssignRolesRes — 分配角色成功响应（空）
type UserAssignRolesRes struct{}
