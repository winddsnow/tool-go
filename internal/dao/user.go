package dao

import (
	"context"

	"tool-go/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var User = userDao{}

type userDao struct {
	Table   string
	Group   string
	Columns userColumns
}

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

func (d *userDao) Ctx(ctx context.Context) *gdb.Model {
	return g.Model(d.Table).Safe().Ctx(ctx)
}

func (d *userDao) Data(data *do.User) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}

func (d *userDao) As(as string) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).As(as)
}
