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

// Role 是全局导出的 DAO 实例，包内通过 init() 初始化。
// 外部使用时直接通过 dao.Role.Ctx(ctx)、dao.Role.Data(data) 调用。
var Role = roleDao{}

// roleDao 是 role 表的 DAO 结构体。
//   - Table：表名 "role"
//   - Group：数据库分组 "default"
//   - Columns：类型安全的列名集合，避免硬编码字符串
type roleDao struct {
	Table   string
	Group   string
	Columns roleColumns
}

// roleColumns 定义 role 表所有列的字符串常量，提供类型安全的列名访问。
type roleColumns struct {
	Id        string
	Name      string
	Code      string
	Sort      string
	Status    string
	Desc      string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

func init() {
	Role = roleDao{
		Table: "role",
		Group: "default",
		Columns: roleColumns{
			Id:        "id",
			Name:      "name",
			Code:      "code",
			Sort:      "sort",
			Status:    "status",
			Desc:      "desc",
			CreatedAt: "created_at",
			UpdatedAt: "updated_at",
			DeletedAt: "deleted_at",
		},
	}
}

// Ctx 创建带上下文的 ORM 查询模型，适用于 SELECT 查询。
// Safe() 保证并发安全，Ctx(ctx) 传递生命周期和链路追踪上下文。
func (d *roleDao) Ctx(ctx context.Context) *gdb.Model {
	return g.Model(d.Table).Safe().Ctx(ctx)
}

// Data 创建带插入/更新数据的 ORM 模型，适用于 INSERT 和 UPDATE。
// 接收 *do.Role 数据对象，零值字段不参与 SQL（部分更新）。
func (d *roleDao) Data(data *do.Role) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}

// As 创建带表别名的 ORM 模型，用于联表查询（JOIN）。
func (d *roleDao) As(as string) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).As(as)
}
