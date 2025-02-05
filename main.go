package main

import (
	"github.com/anyproto/anytype-cli/cmd"
	"github.com/anyproto/anytype-cli/internal"
)

func main() {
	defer internal.CloseGRPCConnection()
	cmd.Execute()
}
