package main

import (
	"github.com/anyproto/anytype-cli/cmd"
	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/autoupdate"
)

func main() {
	autoupdate.CheckAndUpdate()

	defer func() {
		core.CloseEventReceiver()
		core.CloseGRPCConnection()
	}()
	cmd.Execute()
}
