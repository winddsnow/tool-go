package controller

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/middleware"
	"tool-go/internal/service"
)

var User = cUser{}

type cUser struct{}

func (c *cUser) Create(ctx context.Context, req *v1.UserCreateReq) (*v1.UserCreateRes, error) {
	return service.User().Create(ctx, req)
}

func (c *cUser) Delete(ctx context.Context, req *v1.UserDeleteReq) (*v1.UserDeleteRes, error) {
	err := service.User().Delete(ctx, req)
	return &v1.UserDeleteRes{}, err
}

func (c *cUser) Update(ctx context.Context, req *v1.UserUpdateReq) (*v1.UserUpdateRes, error) {
	err := service.User().Update(ctx, req)
	return &v1.UserUpdateRes{}, err
}

func (c *cUser) GetOne(ctx context.Context, req *v1.UserGetOneReq) (*v1.UserGetOneRes, error) {
	return service.User().GetOne(ctx, req)
}

func (c *cUser) List(ctx context.Context, req *v1.UserListReq) (*v1.UserListRes, error) {
	req.UserId = middleware.GetUserId(ctx)
	return service.User().List(ctx, req)
}

func (c *cUser) GetRoles(ctx context.Context, req *v1.UserGetRolesReq) (*v1.UserGetRolesRes, error) {
	return service.User().GetRoles(ctx, req)
}

func (c *cUser) AssignRoles(ctx context.Context, req *v1.UserAssignRolesReq) (*v1.UserAssignRolesRes, error) {
	err := service.User().AssignRoles(ctx, req)
	return &v1.UserAssignRolesRes{}, err
}
