package service

import (
	"context"

	v1 "tool-go/api/v1"
)

type IRole interface {
	Create(ctx context.Context, req *v1.RoleCreateReq) (*v1.RoleCreateRes, error)
	Delete(ctx context.Context, req *v1.RoleDeleteReq) error
	Update(ctx context.Context, req *v1.RoleUpdateReq) error
	GetOne(ctx context.Context, req *v1.RoleGetOneReq) (*v1.RoleGetOneRes, error)
	List(ctx context.Context, req *v1.RoleListReq) (*v1.RoleListRes, error)
}
