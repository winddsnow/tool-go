// Package service 定义业务逻辑接口层。
// 每个业务模块（User/Role/Tools）都遵循相同的模式：
//   Interface 接口定义 → localXxx 单例变量 → xxMu 读写锁 → Xxx() 获取 → RegisterXxx() 注册

package service

import (
	"context"
	"sync"

	v1 "tool-go/api/v1"
)

// ITools 工具箱业务逻辑接口（模拟数据生成等工具类功能）。
// 接口只定义了一个 MockData 方法，sTools 结构体还额外提供了各种 genXxx 内部方法。
// 对外暴露接口方法，对内保留具体类型的能力——这也是接口设计的常见权衡。
type ITools interface {
	MockData(ctx context.Context, req *v1.MockDataReq) (*v1.MockDataRes, error)
}

var (
	// localTools 保存 Tools 业务逻辑单例。
	localTools ITools
	// toolsMu 读写锁，保护 localTools 字段。
	toolsMu sync.RWMutex
)

// Tools 获取 Tools 业务逻辑单例（读安全）。
// 单例模式的优点：全局只有一个实例，避免重复创建开销，也便于统一管理生命周期。
func Tools() ITools {
	toolsMu.RLock()
	defer toolsMu.RUnlock()
	return localTools
}

// RegisterTools 注册 Tools 业务逻辑实现。
// 在 logic/tools.go 的 init() 中调用。
// init() → RegisterTools(New()) → 注册完成 → controller 可调用 service.Tools()
func RegisterTools(i ITools) {
	toolsMu.Lock()
	defer toolsMu.Unlock()
	localTools = i
}
