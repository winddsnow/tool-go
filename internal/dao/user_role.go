// Package dao 是数据访问层（Data Access Object），封装对数据库表的基本操作。
package dao

import (
	"context"

	"tool-go/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// UserRole 是全局导出的 DAO 实例。
var UserRole = userRoleDao{}

// userRoleDao 是 user_role（用户-角色关联表）的 DAO 结构体。
// 多对多关系通过关联表实现：一个用户可有多个角色，一个角色可被多个用户拥有。
type userRoleDao struct {
	Table   string
	Group   string
	Columns userRoleColumns
}

// userRoleColumns 定义 user_role 关联表的列名常量。
type userRoleColumns struct {
	Id     string
	UserId string
	RoleId string
}

func init() {
	UserRole = userRoleDao{
		Table: "user_role",
		Group: "default",
		Columns: userRoleColumns{
			Id:     "id",
			UserId: "user_id",
			RoleId: "role_id",
		},
	}
}

// Ctx 创建带上下文的 ORM 查询模型，适用于 SELECT 查询。
func (d *userRoleDao) Ctx(ctx context.Context) *gdb.Model {
	return g.Model(d.Table).Safe().Ctx(ctx)
}

// Data 创建带插入/更新数据的 ORM 模型，适用于 INSERT 和 UPDATE。
func (d *userRoleDao) Data(data *do.UserRole) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}
