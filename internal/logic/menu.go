package logic

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
	"tool-go/internal/model/do"
	"tool-go/internal/model/entity"
	"tool-go/internal/service"
)

func init() {
	service.RegisterMenu(NewMenu())
}

func NewMenu() service.IMenu {
	return &sMenu{}
}

type sMenu struct{}

func (s *sMenu) Create(ctx context.Context, req *v1.MenuCreateReq) (*v1.MenuCreateRes, error) {
	result, err := dao.Menu.Data(&do.Menu{
		ParentId:  req.ParentId,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Visible:   req.Visible,
		Status:    req.Status,
		Type:      req.Type,
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}).Insert()
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.MenuCreateRes{Id: uint64(id)}, nil
}

func (s *sMenu) Delete(ctx context.Context, req *v1.MenuDeleteReq) error {
	_, err := dao.Menu.Ctx(ctx).Where(dao.Menu.Columns.Id, req.Id).Update(gdb.Map{
		dao.Menu.Columns.DeletedAt: gtime.Now(),
	})
	return err
}

func (s *sMenu) Update(ctx context.Context, req *v1.MenuUpdateReq) error {
	data := do.Menu{}
	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Path != "" {
		data.Path = req.Path
	}
	if req.Component != "" {
		data.Component = req.Component
	}
	if req.Icon != "" {
		data.Icon = req.Icon
	}
	if req.ParentId != nil {
		data.ParentId = *req.ParentId
	}
	if req.Sort != nil {
		data.Sort = *req.Sort
	}
	if req.Visible != nil {
		data.Visible = *req.Visible
	}
	if req.Status != nil {
		data.Status = *req.Status
	}
	if req.Type != nil {
		data.Type = *req.Type
	}
	data.UpdatedAt = gtime.Now()

	_, err := dao.Menu.Ctx(ctx).Where(dao.Menu.Columns.Id, req.Id).Data(data).Update()
	return err
}

func (s *sMenu) GetOne(ctx context.Context, req *v1.MenuGetOneReq) (*v1.MenuGetOneRes, error) {
	var menu *entity.Menu
	err := dao.Menu.Ctx(ctx).
		Where(dao.Menu.Columns.Id, req.Id).
		WhereNull(dao.Menu.Columns.DeletedAt).
		Scan(&menu)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, gerror.New("菜单不存在")
	}
	return &v1.MenuGetOneRes{
		Id:        menu.Id,
		ParentId:  menu.ParentId,
		Name:      menu.Name,
		Path:      menu.Path,
		Component: menu.Component,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
		Visible:   menu.Visible,
		Status:    menu.Status,
		Type:      menu.Type,
		CreatedAt: menu.CreatedAt.String(),
		UpdatedAt: menu.UpdatedAt.String(),
	}, nil
}

func (s *sMenu) List(ctx context.Context, req *v1.MenuListReq) (*v1.MenuListRes, error) {
	m := dao.Menu.Ctx(ctx).WhereNull(dao.Menu.Columns.DeletedAt)
	if req.Name != "" {
		m = m.WhereLike(dao.Menu.Columns.Name, "%"+req.Name+"%")
	}
	if req.Status != nil {
		m = m.Where(dao.Menu.Columns.Status, *req.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var list []*entity.Menu
	err = m.Page(req.Page, req.PageSize).
		OrderAsc(dao.Menu.Columns.Sort).
		OrderDesc(dao.Menu.Columns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	items := make([]v1.MenuItem, 0, len(list))
	for _, item := range list {
		items = append(items, v1.MenuItem{
			Id:        item.Id,
			ParentId:  item.ParentId,
			Name:      item.Name,
			Path:      item.Path,
			Component: item.Component,
			Icon:      item.Icon,
			Sort:      item.Sort,
			Visible:   item.Visible,
			Status:    item.Status,
			Type:      item.Type,
			CreatedAt: item.CreatedAt.String(),
		})
	}
	return &v1.MenuListRes{List: items, Total: total, Page: req.Page}, nil
}

func (s *sMenu) GetUserMenus(ctx context.Context, userId uint64) (*v1.MenuGetUserMenusRes, error) {
	var menus []*entity.Menu
	err := dao.Menu.Ctx(ctx).
		LeftJoin("role_menu", "role_menu.menu_id=menu.id").
		LeftJoin("user_role", "user_role.role_id=role_menu.role_id").
		Where("user_role.user_id", userId).
		WhereNull(dao.Menu.Columns.DeletedAt).
		Where(dao.Menu.Columns.Status, 1).
		Fields("menu.*").
		Distinct().
		OrderAsc(dao.Menu.Columns.Sort).
		Scan(&menus)
	if err != nil {
		return nil, err
	}

	tree := buildMenuTree(menus, 0)
	return &v1.MenuGetUserMenusRes{Menus: tree}, nil
}

func buildMenuTree(menus []*entity.Menu, parentId uint64) []v1.MenuTree {
	var tree []v1.MenuTree
	for _, m := range menus {
		if m.ParentId == parentId {
			node := v1.MenuTree{
				Id:        m.Id,
				ParentId:  m.ParentId,
				Name:      m.Name,
				Path:      m.Path,
				Component: m.Component,
				Icon:      m.Icon,
				Sort:      m.Sort,
				Visible:   m.Visible,
				Type:      m.Type,
				Children:  buildMenuTree(menus, m.Id),
			}
			tree = append(tree, node)
		}
	}
	sort.Slice(tree, func(i, j int) bool {
		return tree[i].Sort < tree[j].Sort
	})
	return tree
}
