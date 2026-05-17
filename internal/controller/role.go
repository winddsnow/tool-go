package controller

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/service"
)

var Role = cRole{}

type cRole struct{}

func (c *cRole) Create(ctx context.Context, req *v1.RoleCreateReq) (*v1.RoleCreateRes, error) {
	return service.Role().Create(ctx, req)
}

func (c *cRole) Delete(ctx context.Context, req *v1.RoleDeleteReq) (*v1.RoleDeleteRes, error) {
	err := service.Role().Delete(ctx, req)
	return &v1.RoleDeleteRes{}, err
}

func (c *cRole) Update(ctx context.Context, req *v1.RoleUpdateReq) (*v1.RoleUpdateRes, error) {
	err := service.Role().Update(ctx, req)
	return &v1.RoleUpdateRes{}, err
}

func (c *cRole) GetOne(ctx context.Context, req *v1.RoleGetOneReq) (*v1.RoleGetOneRes, error) {
	return service.Role().GetOne(ctx, req)
}

func (c *cRole) List(ctx context.Context, req *v1.RoleListReq) (*v1.RoleListRes, error) {
	return service.Role().List(ctx, req)
}
