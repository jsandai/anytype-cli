package typecmd

import (
	"github.com/spf13/cobra"

	typeCreateCmd "github.com/anyproto/anytype-cli/cmd/type/create"
	typeListCmd "github.com/anyproto/anytype-cli/cmd/type/list"
)

func NewTypeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "type <command>",
		Short: "Manage object types",
		Long:  "Create and list custom object types",
	}

	cmd.AddCommand(typeCreateCmd.NewCreateCmd())
	cmd.AddCommand(typeListCmd.NewListCmd())

	return cmd
}
