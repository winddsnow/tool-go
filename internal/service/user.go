package service

import (
	"context"
	"sync"

	v1 "tool-go/api/v1"
)

type IUser interface {
	Create(ctx context.Context, req *v1.UserCreateReq) (*v1.UserCreateRes, error)
	Delete(ctx context.Context, req *v1.UserDeleteReq) error
	Update(ctx context.Context, req *v1.UserUpdateReq) error
	GetOne(ctx context.Context, req *v1.UserGetOneReq) (*v1.UserGetOneRes, error)
	List(ctx context.Context, req *v1.UserListReq) (*v1.UserListRes, error)
}

var (
	localUser IUser
	userMu    sync.RWMutex
)

func User() IUser {
	userMu.RLock()
	defer userMu.RUnlock()
	return localUser
}

func RegisterUser(i IUser) {
	userMu.Lock()
	defer userMu.Unlock()
	localUser = i
}
