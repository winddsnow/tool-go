// Package dao 是数据访问层（Data Access Object），封装对数据库表的基本操作。
// DAO 模式的核心思想是将数据访问逻辑与业务逻辑分离：
//   - 业务层（logic）只调用 DAO 方法，不需要知道 SQL 细节
//   - DAO 层提供类型安全的列名引用（通过 Columns 结构体），避免拼写错误
//   - 当数据库表结构变化时，只需修改 DAO 层，业务层不受影响
package dao

import (
	"context"

	"tool-go/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// User 是全局导出的 DAO 实例，包内通过 init() 初始化。
// 外部使用时直接通过 dao.User.Ctx(ctx)、dao.User.Data(data) 调用。
var User = userDao{}

// userDao 是 user 表的 DAO 结构体。
//   - Table：表名 "user"，用于 g.Model(Table) 获取 ORM 操作对象
//   - Group：数据库分组 "default"，用于区分多数据源
//   - Columns：类型安全的列名集合，通过 Columns.Id、Columns.Username 引用列名，
//     避免在代码中直接硬编码字符串，防止拼写错误，同时支持 IDE 代码补全
type userDao struct {
	Table   string
	Group   string
	Columns userColumns
}

// userColumns 定义 user 表所有列的字符串常量，提供类型安全的列名访问。
// 每个结构体字段对应数据库表的一列，字段值即为数据库列名的字符串。
// 例如 Columns.CreatedAt 的值是 "created_at"，在 WHERE 条件中可写为：
//   dao.User.Ctx(ctx).Where(dao.User.Columns.CreatedAt, "2024-01-01")
type userColumns struct {
	Id        string
	Username  string
	Password  string
	Salt      string
	Nickname  string
	Email     string
	Phone     string
	Status    string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

func init() {
	User = userDao{
		Table: "user",
		Group: "default",
		Columns: userColumns{
			Id:        "id",
			Username:  "username",
			Password:  "password",
			Salt:      "salt",
			Nickname:  "nickname",
			Email:     "email",
			Phone:     "phone",
			Status:    "status",
			CreatedAt: "created_at",
			UpdatedAt: "updated_at",
			DeletedAt: "deleted_at",
		},
	}
}

// Ctx 创建带上下文的 ORM 查询模型（query chain），适用于 SELECT 查询。
// 链式调用解析：
//   g.Model(d.Table)  — 获取 user 表对应的 ORM 操作对象
//   .Safe()           — 安全模式：创建 ORM 模型的副本而非直接修改原对象。
//                        GoFrame 的 ORM 链式调用会累积状态（如 Where、OrderBy），
//                        不加 Safe() 时多个并发请求会相互污染链式条件。
//                        Safe() 每次返回新对象，保证并发安全。
//   .Ctx(ctx)         — 绑定上下文，用于传递请求级别的生命周期控制（超时、取消）、
//                        链路追踪（trace ID）和日志上下文。
func (d *userDao) Ctx(ctx context.Context) *gdb.Model {
	return g.Model(d.Table).Safe().Ctx(ctx)
}

// Data 创建带插入/更新数据的 ORM 模型，适用于 INSERT 和 UPDATE 操作。
// 接收 *do.User（Data Object / 数据对象），其中的字段类型为 any（interface{}），
// 与 entity 层的区别：
//   - entity：类型严格，包含所有字段，主要用于数据库读取后的结果映射
//   - do：字段类型为 any，只包含需要写入的字段，零值字段（如空字符串、0）不参与 SQL，
//     实现精确的部分更新（partial update）
// gctx.New() 创建一个背景上下文，用于没有请求上下文（如定时任务）的场景。
func (d *userDao) Data(data *do.User) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}

// As 创建带表别名的 ORM 模型，用于联表查询（JOIN）时简化 SQL 语句。
// 例如 .As("u") 可在后续查询中使用 u.id、u.username 等引用。
func (d *userDao) As(as string) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).As(as)
}
