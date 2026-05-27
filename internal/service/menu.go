package service

import (
	"context"
	"sync"

	v1 "tool-go/api/v1"
)

type IMenu interface {
	Create(ctx context.Context, req *v1.MenuCreateReq) (*v1.MenuCreateRes, error)
	Delete(ctx context.Context, req *v1.MenuDeleteReq) error
	Update(ctx context.Context, req *v1.MenuUpdateReq) error
	GetOne(ctx context.Context, req *v1.MenuGetOneReq) (*v1.MenuGetOneRes, error)
	List(ctx context.Context, req *v1.MenuListReq) (*v1.MenuListRes, error)
	GetUserMenus(ctx context.Context, userId uint64) (*v1.MenuGetUserMenusRes, error)
}

var (
	localMenu IMenu
	menuMu    sync.RWMutex
)

func Menu() IMenu {
	menuMu.RLock()
	defer menuMu.RUnlock()
	return localMenu
}

func RegisterMenu(i IMenu) {
	menuMu.Lock()
	defer menuMu.Unlock()
	localMenu = i
}
