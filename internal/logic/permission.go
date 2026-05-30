package logic

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
	"tool-go/internal/model/entity"
	"tool-go/internal/service"
)

func init() {
	service.RegisterPermission(NewPermission())
}

func NewPermission() service.IPermission {
	return &sPermission{}
}

type sPermission struct{}

func (s *sPermission) List(ctx context.Context, req *v1.PermissionListReq) (*v1.PermissionListRes, error) {
	m := dao.Permission.Ctx(ctx)

	if req.Code != "" {
		m = m.WhereLike(dao.Permission.Columns.Code, "%"+req.Code+"%")
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var list []*entity.Permission
	err = m.Page(req.Page, req.PageSize).
		OrderDesc(dao.Permission.Columns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}

	items := make([]v1.PermissionItem, 0, len(list))
	for _, p := range list {
		items = append(items, v1.PermissionItem{
			Id:     p.Id,
			Code:   p.Code,
			Name:   p.Name,
			MenuId: p.MenuId,
		})
	}

	return &v1.PermissionListRes{List: items, Total: total, Page: req.Page}, nil
}
