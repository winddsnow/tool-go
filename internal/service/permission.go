package service

import (
	"context"
	"sync"

	v1 "tool-go/api/v1"
)

type IPermission interface {
	List(ctx context.Context, req *v1.PermissionListReq) (*v1.PermissionListRes, error)
}

var (
	localPermission IPermission
	permissionMu    sync.RWMutex
)

func Permission() IPermission {
	permissionMu.RLock()
	defer permissionMu.RUnlock()
	return localPermission
}

func RegisterPermission(i IPermission) {
	permissionMu.Lock()
	defer permissionMu.Unlock()
	localPermission = i
}
