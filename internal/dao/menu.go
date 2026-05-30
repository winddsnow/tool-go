package dao

import (
    "context"
    "tool-go/internal/model/do"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

var Menu = menuDao{}

type menuDao struct {
    Table   string
    Group   string
    Columns menuColumns
}

type menuColumns struct {
    Id        string
    ParentId  string
    Name      string
    Path      string
    Component string
    Icon      string
    Sort      string
    Visible   string
    Status    string
    Type      string
    CreatedAt string
    UpdatedAt string
    DeletedAt string
}

func init() {
    Menu = menuDao{
        Table: "menu",
        Group: "default",
        Columns: menuColumns{
            Id:        "id",
            ParentId:  "parent_id",
            Name:      "name",
            Path:      "path",
            Component: "component",
            Icon:      "icon",
            Sort:      "sort",
            Visible:   "visible",
            Status:    "status",
            Type:      "type",
            CreatedAt: "created_at",
            UpdatedAt: "updated_at",
            DeletedAt: "deleted_at",
        },
    }
}

func (d *menuDao) Ctx(ctx context.Context) *gdb.Model {
    return g.Model(d.Table).Safe().Ctx(ctx)
}

func (d *menuDao) Data(data *do.Menu) *gdb.Model {
    return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}
