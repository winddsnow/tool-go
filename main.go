package main

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/os/gctx"
	_ "tool-go/internal/logic"

	"tool-go/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
