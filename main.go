// ============================================================
// package main — Go 程序的入口包
// ------------------------------------------------------------
// Go 语言规定，可执行程序必须包含 package main 和 func main()。
// main 包是特殊的：它定义了一个可独立运行的程序。
// 其他包（如 package v1、package cmd）都是被 main 包导入使用的。
//
// 当运行 go build 时，编译器从 main 包的 main() 函数开始执行。
// ============================================================
package main

import (
	// _ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	// ------------------------------------------------------------
	// _ 是 Go 的"空白标识符"（blank identifier）。
	// 当导入一个包但没有在代码中直接使用时，Go 编译器会报错。
	// 使用 _ 前缀告诉编译器："我知道这个包没有直接使用，
	// 但我需要它的 init() 函数在程序启动时自动执行。"
	//
	// 这个包是 PostgreSQL 数据库驱动。它通过 init() 函数
	// 向 database/sql 注册 PostgreSQL 驱动，这样 GoFrame 的
	// ORM 才能连接和操作 PostgreSQL 数据库。
	//
	// init() 函数是 Go 的特殊函数：
	//   • 每个包可以有多个 init() 函数
	//   • init() 在包被导入时自动执行
	//   • 执行顺序：先执行依赖包的 init()，后执行当前包的 init()
	//   • 最终执行 main() 函数
	//
	// 如果不导入这个驱动，程序连接 PostgreSQL 时会报错：
	// "sql: unknown driver 'pgsql' (forgotten import?)"
	// ============================================================
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	// gctx — GoFrame 上下文工具包
	// gctx.GetInitCtx() 创建一个框架级初始化上下文，
	// 用于在 main 函数中初始化框架组件（数据库、缓存等）。
	"github.com/gogf/gf/v2/os/gctx"

	// _ "tool-go/internal/logic"
	// 同理，使用空白标识符导入 logic 包。
	// logic 包中的 init() 函数会注册所有业务逻辑处理器，
	// 这样控制器（controller）才能调用它们。
	// 这种"通过 init() 注册"的模式在 GoFrame 中很常见，
	// 避免了显式的注册代码，使 main 函数保持简洁。
	_ "tool-go/internal/logic"

	// cmd — 命令行命令定义包
	// 导入它的 Main 命令对象进行启动。
	"tool-go/internal/cmd"
)

// ============================================================
// func main() — 程序入口函数
// ------------------------------------------------------------
// func 是 Go 定义函数的关键字。
// main 是函数名，没有参数，没有返回值。
// 这是 Go 程序的唯一入口点。
//
// 执行流程：
//   1. Go runtime 初始化（分配堆栈、启动垃圾回收器）
//   2. 按依赖顺序执行所有导入包的 init() 函数
//      → pgsql 驱动注册自己
//      → logic 包注册业务逻辑
//   3. 执行 main() 函数
//   4. main() 返回，程序退出
// ============================================================
func main() {
	// cmd.Main — 我们在 cmd 包中定义的 gcmd.Command 变量
	// .Run() — 执行这个命令，启动 HTTP 服务器
	// gctx.GetInitCtx() — 创建 GoFrame 初始化上下文，
	//                     框架会在上下文中设置数据库连接池、
	//                     缓存、日志等基础服务的初始化参数。
	//
	// 整个程序的核心就这一行：
	// 解析路由 → 启动服务器 → 监听请求 → 处理请求 → 返回响应
	cmd.Main.Run(gctx.GetInitCtx())
}
