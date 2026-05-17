// Package controller 包含所有 HTTP 请求处理器。
// 本文件是角色管理相关的 Controller。
// 对比 user.go，role.go 的模式完全一致，展示了 GoFrame CRUD Controller 的标准样板代码。
package controller

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/service"
)

// var Role = cRole{} 是 package-level 单例。
// 在 Go 中，var 声明的初始化顺序由包的依赖关系决定：导入顺序决定初始化顺序。
// cRole{} 创建零值空结构体，Go 保证包变量初始化在 main 函数执行前完成。
var Role = cRole{}

// cRole 是角色管理 controller。
// 空结构体（struct{}）在 Go 中不占用内存空间（zero size），
// 常用于只挂载方法、不存储状态的场景。
// 指针接收器 (c *cRole) 虽然不是必须的（空结构体值/指针无区别），
// 但 Go 社区惯例所有方法接收器使用一致的类型。
type cRole struct{}

// Create 创建新角色，委托给 service 层。
// service.Role() 返回 service 单例，Create 方法内部进行参数校验和数据库写入。
// Go 支持多返回值（return a, b），这是 Go 错误处理的核心机制：
// 每个可能失败的调用都应检查 err。
func (c *cRole) Create(ctx context.Context, req *v1.RoleCreateReq) (*v1.RoleCreateRes, error) {
	return service.Role().Create(ctx, req)
}

// Delete 删除角色（软删除）。返回空响应 + err 的模式。
// 当 service.Role().Delete 返回 nil 时，框架序列化 v1.RoleDeleteRes{} 为 JSON {}。
func (c *cRole) Delete(ctx context.Context, req *v1.RoleDeleteReq) (*v1.RoleDeleteRes, error) {
	err := service.Role().Delete(ctx, req)
	return &v1.RoleDeleteRes{}, err
}

// Update 更新角色信息。
func (c *cRole) Update(ctx context.Context, req *v1.RoleUpdateReq) (*v1.RoleUpdateRes, error) {
	err := service.Role().Update(ctx, req)
	return &v1.RoleUpdateRes{}, err
}

// GetOne 获取单个角色详情。
func (c *cRole) GetOne(ctx context.Context, req *v1.RoleGetOneReq) (*v1.RoleGetOneRes, error) {
	return service.Role().GetOne(ctx, req)
}

// List 获取角色列表（分页）。
// 注意与 user.go 不同，这里不需要注入 UserId，因为角色列表通常是公开的管理功能。
func (c *cRole) List(ctx context.Context, req *v1.RoleListReq) (*v1.RoleListRes, error) {
	return service.Role().List(ctx, req)
}
