// Package controller 包含所有 HTTP 请求处理器。
// Go 的 package 机制：同一目录下的所有 .go 文件必须属于同一个 package，文件名不影响编译。
// 本文件是用户管理相关的 Controller，负责接收 HTTP 请求、参数校验、
// 调用 service 层处理业务逻辑、返回响应。
//
// controller 层遵循"薄 Controller，厚 Service"原则：
// controller 只做"编排"（解析参数、调用 service、包装响应），业务逻辑在 service 中实现。
package controller

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/middleware"
	"tool-go/internal/service"
)

// var User = cUser{} 是 Go 的 package-level 单例模式。
// var 声明在包初始化时执行，整个进程生命周期内只有一个实例。
// cUser{} 创建空结构体（零内存占用），赋值给可导出的 User 变量，
// 路由注册时通过 controller.User.List 来引用 controller 方法。
var User = cUser{}

// cUser 是用户管理 controller，私有类型（小写 c），外部通过 var User 引用实例。
// 方法接收器 (c *cUser) 使用指针，虽然空结构体用值接收器也一样，
// 但 Go 惯例保持一致性（全部用指针或全部用值）。
type cUser struct{}

// Create 创建新用户，直接委托给 service.User().Create。
// service.User() 返回一个单例的 service 实例（类似本文件的 var User），
// 实现了"controller -> service -> dao"的三层架构。
//
// 三层的职责：
//   - controller：HTTP 语义层（请求解析、响应格式、权限检查）
//   - service：业务逻辑层（事务管理、数据校验、跨表操作）
//   - dao：数据访问层（SQL 生成与执行）
func (c *cUser) Create(ctx context.Context, req *v1.UserCreateReq) (*v1.UserCreateRes, error) {
	return service.User().Create(ctx, req)
}

// Delete 标记删除用户（软删除），调用 service 层执行。
// 注意这里返回的是 &v1.UserDeleteRes{}, err：
//   - err == nil 时框架返回 v1.UserDeleteRes{}（空结构体）作为 JSON body，HTTP 200
//   - err != nil 时框架返回错误响应，忽略第一个返回值
// 这种模式在 Go 中常见：当响应为空时仍返回空结构体以保持 API 一致性。
func (c *cUser) Delete(ctx context.Context, req *v1.UserDeleteReq) (*v1.UserDeleteRes, error) {
	err := service.User().Delete(ctx, req)
	return &v1.UserDeleteRes{}, err
}

// Update 更新用户信息。
func (c *cUser) Update(ctx context.Context, req *v1.UserUpdateReq) (*v1.UserUpdateRes, error) {
	err := service.User().Update(ctx, req)
	return &v1.UserUpdateRes{}, err
}

// GetOne 获取单个用户详情。
func (c *cUser) GetOne(ctx context.Context, req *v1.UserGetOneReq) (*v1.UserGetOneRes, error) {
	return service.User().GetOne(ctx, req)
}

// List 获取用户列表（分页）。在调用 service 之前注入当前请求者的 UserId，
// 用于数据权限过滤（如非管理员只能看到自己）。
// middleware.GetUserId(ctx) 从 context 中提取 JWT 中间件设置的当前用户 ID。
func (c *cUser) List(ctx context.Context, req *v1.UserListReq) (*v1.UserListRes, error) {
	req.UserId = middleware.GetUserId(ctx)
	return service.User().List(ctx, req)
}

// GetRoles 获取指定用户的角色分配。
func (c *cUser) GetRoles(ctx context.Context, req *v1.UserGetRolesReq) (*v1.UserGetRolesRes, error) {
	return service.User().GetRoles(ctx, req)
}

// AssignRoles 为用户分配角色。
func (c *cUser) AssignRoles(ctx context.Context, req *v1.UserAssignRolesReq) (*v1.UserAssignRolesRes, error) {
	err := service.User().AssignRoles(ctx, req)
	return &v1.UserAssignRolesRes{}, err
}
