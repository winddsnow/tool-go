package main

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/os/gctx"

	"tool-go/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
