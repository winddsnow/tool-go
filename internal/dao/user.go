package dao

import (
	"tool-go/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gctx"
)

var User = userDao{}

type userDao struct {
	table    string
	group    string
	columns  userColumns
}

type userColumns struct {
	Id        string
	Username  string
	Password  string
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
		table:   "user",
		group:   "default",
		columns: userColumns{
			Id:        "id",
			Username:  "username",
			Password:  "password",
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

func (dao *userDao) DB(ctx context.Context) *gdb.Model {
	return gdb.Ctx(ctx).Safe().Model(dao.table).Safe()
}

func (dao *userDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return gdb.Ctx(ctx).Transaction(ctx, f)
}

func (dao *userDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB(ctx)
}

func (dao *userDao) Data(data *do.User) *gdb.Model {
	return dao.DB(gctx.New()).Data(data)
}

func (dao *userDao) As(as string) *gdb.Model {
	return dao.DB(gctx.New()).As(as)
}
