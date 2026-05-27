package controller

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/middleware"
	"tool-go/internal/service"
)

var Menu = cMenu{}

type cMenu struct{}

func (c *cMenu) Create(ctx context.Context, req *v1.MenuCreateReq) (*v1.MenuCreateRes, error) {
	return service.Menu().Create(ctx, req)
}

func (c *cMenu) Delete(ctx context.Context, req *v1.MenuDeleteReq) (*v1.MenuDeleteRes, error) {
	err := service.Menu().Delete(ctx, req)
	return &v1.MenuDeleteRes{}, err
}

func (c *cMenu) Update(ctx context.Context, req *v1.MenuUpdateReq) (*v1.MenuUpdateRes, error) {
	err := service.Menu().Update(ctx, req)
	return &v1.MenuUpdateRes{}, err
}

func (c *cMenu) GetOne(ctx context.Context, req *v1.MenuGetOneReq) (*v1.MenuGetOneRes, error) {
	return service.Menu().GetOne(ctx, req)
}

func (c *cMenu) List(ctx context.Context, req *v1.MenuListReq) (*v1.MenuListRes, error) {
	return service.Menu().List(ctx, req)
}

func (c *cMenu) GetUserMenus(ctx context.Context, req *v1.MenuGetUserMenusReq) (*v1.MenuGetUserMenusRes, error) {
	userId := middleware.GetUserId(ctx)
	return service.Menu().GetUserMenus(ctx, userId)
}
