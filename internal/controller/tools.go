// Package controller 包含所有 HTTP 请求处理器。
// 本文件是工具类（Tools）相关的 Controller，提供开发辅助功能（如数据填充）。
package controller

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/service"
)

// var Tools = cTools{} 是 package-level 单例。
// 在 Go 中，全局变量（package-level var）在程序启动时初始化，
// 且 Go 保证所有 import 的包初始化完成后再执行 main 函数。
var Tools = cTools{}

// cTools 是工具类 controller，目前只包含 MockData（生成模拟数据）接口。
type cTools struct{}

// MockData 生成模拟测试数据，供开发和测试使用。
// 这个接口通常需要 super_admin 权限（在路由注册时通过 middleware.Permission 控制）。
// 它完全委托给 service.Tools().MockData，controller 层不做任何额外处理，
// 是"薄 Controller"原则的典型示例。
func (c *cTools) MockData(ctx context.Context, req *v1.MockDataReq) (*v1.MockDataRes, error) {
	return service.Tools().MockData(ctx, req)
}
