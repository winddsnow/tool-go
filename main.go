package main

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2" // PostgreSQL 驱动
	"github.com/gogf/gf/v2/os/gctx"
	_ "tool-go/internal/logic" // 注册业务逻辑
	"tool-go/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
