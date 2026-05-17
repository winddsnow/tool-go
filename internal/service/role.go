// Package service 定义业务逻辑接口层。
// 通过接口抽象，controller 层只需依赖接口，无需关注具体实现细节。
// 这种模式在 Go 中称为"面向接口编程"。

package service

import (
	"context"
	"sync"

	v1 "tool-go/api/v1"
)

// IRole 角色管理的业务逻辑接口。
// Go 接口是隐式满足的——logic/sRole 只要实现了这些方法就自动实现了 IRole。
type IRole interface {
	Create(ctx context.Context, req *v1.RoleCreateReq) (*v1.RoleCreateRes, error)
	Delete(ctx context.Context, req *v1.RoleDeleteReq) error
	Update(ctx context.Context, req *v1.RoleUpdateReq) error
	GetOne(ctx context.Context, req *v1.RoleGetOneReq) (*v1.RoleGetOneRes, error)
	List(ctx context.Context, req *v1.RoleListReq) (*v1.RoleListRes, error)
}

var (
	// localRole 保存 Role 业务逻辑单例。
	localRole IRole
	// roleMu 读写锁，保护 localRole 的并发访问。
	// sync.RWMutex 比 sync.Mutex 更适合读多写少的单例场景。
	roleMu sync.RWMutex
)

// Role 获取 Role 业务逻辑单例（读安全）。
func Role() IRole {
	roleMu.RLock()
	defer roleMu.RUnlock()
	return localRole
}

// RegisterRole 注册 Role 业务逻辑实现。
// 在 logic/role.go 的 init() 中调用，通过 Go 的 init 机制在 main() 前自动完成注册。
func RegisterRole(i IRole) {
	roleMu.Lock()
	defer roleMu.Unlock()
	localRole = i
}
