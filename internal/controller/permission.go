package controller

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/service"
)

var Permission = cPermission{}

type cPermission struct{}

func (c *cPermission) List(ctx context.Context, req *v1.PermissionListReq) (*v1.PermissionListRes, error) {
	return service.Permission().List(ctx, req)
}
