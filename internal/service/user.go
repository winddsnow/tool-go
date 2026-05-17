// Package service 定义业务逻辑接口层。
// Go 中接口（interface）是方法的集合，定义行为契约而不关心具体实现。
// 设计模式：依赖倒置原则——高层模块（controller）不依赖低层模块（logic），
// 两者都依赖抽象（接口）。这样便于测试（可替换 mock）和解耦。

package service

import (
	"context"
	"sync"

	v1 "tool-go/api/v1"
)

// IUser 用户管理的业务逻辑接口。
// Go 的接口是隐式实现（duck typing）——不需要 implements 关键字，
// 只要类型实现了接口中所有方法，就自动满足该接口。
type IUser interface {
	Create(ctx context.Context, req *v1.UserCreateReq) (*v1.UserCreateRes, error)
	Delete(ctx context.Context, req *v1.UserDeleteReq) error
	Update(ctx context.Context, req *v1.UserUpdateReq) error
	GetOne(ctx context.Context, req *v1.UserGetOneReq) (*v1.UserGetOneRes, error)
	List(ctx context.Context, req *v1.UserListReq) (*v1.UserListRes, error)
	GetRoles(ctx context.Context, req *v1.UserGetRolesReq) (*v1.UserGetRolesRes, error)
	AssignRoles(ctx context.Context, req *v1.UserAssignRolesReq) error
}

var (
	// localUser 保存 User 业务逻辑的单例实例。
	// 通过 RegisterUser() 在 init() 中注入，通过 User() 读取。
	localUser IUser
	// userMu 读写锁，保护 localUser 的并发读写安全。
	// sync.RWMutex 是 Go 标准库的读写互斥锁：
	//   RLock() 允许多个 goroutine 同时读，
	//   Lock()  只允许一个 goroutine 写，写时禁止读。
	// 场景：User() 读取频繁（RLock），RegisterUser() 只初始化一次（Lock）。
	userMu sync.RWMutex
)

// User 获取 User 业务逻辑单例。
// 使用 RLock（读锁）而非 Lock（写锁），因为读取操作不修改共享变量。
// RLock 可以被多个 goroutine 同时持有，比 Lock 有更好的并发性能。
func User() IUser {
	userMu.RLock()
	defer userMu.RUnlock() // defer 确保在函数返回前释放锁
	return localUser
}

// RegisterUser 注册 User 业务逻辑实现。
// 一般在 logic 包的 init() 函数中调用。使用 Lock（写锁）确保写入安全。
// 注册逻辑：logic 包的 init() 在程序启动时自动执行 → 调用 RegisterUser(NewUser())
// → 将实现存入 localUser → controller 通过 User() 获取。
func RegisterUser(i IUser) {
	userMu.Lock()
	defer userMu.Unlock()
	localUser = i
}
