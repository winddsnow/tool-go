package logic

import (
	"context"

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
	service.RegisterRole(NewRole())
}

func NewRole() service.IRole {
	return &sRole{}
}

type sRole struct{}

func (s *sRole) Create(ctx context.Context, req *v1.RoleCreateReq) (*v1.RoleCreateRes, error) {
	count, err := dao.Role.Ctx(ctx).Where(dao.Role.Columns.Code, req.Code).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, gerror.New("角色编码已存在")
	}

	result, err := dao.Role.Data(&do.Role{
		Name:      req.Name,
		Code:      req.Code,
		Sort:      req.Sort,
		Status:    req.Status,
		Desc:      req.Desc,
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}).Insert()
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &v1.RoleCreateRes{Id: uint64(id)}, nil
}

func (s *sRole) Delete(ctx context.Context, req *v1.RoleDeleteReq) error {
	_, err := dao.Role.Ctx(ctx).Where(dao.Role.Columns.Id, req.Id).Update(gdb.Map{
		dao.Role.Columns.DeletedAt: gtime.Now(),
	})
	return err
}

func (s *sRole) Update(ctx context.Context, req *v1.RoleUpdateReq) error {
	data := do.Role{}
	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Code != "" {
		data.Code = req.Code
	}
	data.Sort = req.Sort
	data.Status = req.Status
	data.Desc = req.Desc
	data.UpdatedAt = gtime.Now()

	_, err := dao.Role.Ctx(ctx).Where(dao.Role.Columns.Id, req.Id).Data(data).Update()
	return err
}

func (s *sRole) GetOne(ctx context.Context, req *v1.RoleGetOneReq) (*v1.RoleGetOneRes, error) {
	var role *entity.Role
	err := dao.Role.Ctx(ctx).Where(dao.Role.Columns.Id, req.Id).WhereNull(dao.Role.Columns.DeletedAt).Scan(&role)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, gerror.New("role not found")
	}

	return &v1.RoleGetOneRes{
		Id:        role.Id,
		Name:      role.Name,
		Code:      role.Code,
		Sort:      role.Sort,
		Status:    role.Status,
		Desc:      role.Desc,
		CreatedAt: role.CreatedAt.String(),
		UpdatedAt: role.UpdatedAt.String(),
	}, nil
}

func (s *sRole) List(ctx context.Context, req *v1.RoleListReq) (*v1.RoleListRes, error) {
	m := dao.Role.Ctx(ctx).WhereNull(dao.Role.Columns.DeletedAt)

	if req.Name != "" {
		m = m.WhereLike(dao.Role.Columns.Name, "%"+req.Name+"%")
	}
	if req.Status != nil {
		m = m.Where(dao.Role.Columns.Status, *req.Status)
	}

	var total int
	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var list []*entity.Role
	err = m.Page(req.Page, req.PageSize).OrderAsc(dao.Role.Columns.Sort).OrderDesc(dao.Role.Columns.Id).Scan(&list)
	if err != nil {
		return nil, err
	}

	items := make([]v1.RoleItem, 0, len(list))
	for _, r := range list {
		items = append(items, v1.RoleItem{
			Id:        r.Id,
			Name:      r.Name,
			Code:      r.Code,
			Sort:      r.Sort,
			Status:    r.Status,
			Desc:      r.Desc,
			CreatedAt: r.CreatedAt.String(),
		})
	}

	return &v1.RoleListRes{
		List:  items,
		Total: total,
		Page:  req.Page,
	}, nil
}
