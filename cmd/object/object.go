package object

import (
	"github.com/spf13/cobra"

	objectCreateCmd "github.com/anyproto/anytype-cli/cmd/object/create"
	objectDeleteCmd "github.com/anyproto/anytype-cli/cmd/object/delete"
	objectGetCmd "github.com/anyproto/anytype-cli/cmd/object/get"
	objectSearchCmd "github.com/anyproto/anytype-cli/cmd/object/search"
	objectUpdateCmd "github.com/anyproto/anytype-cli/cmd/object/update"
)

func NewObjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "object <command>",
		Short: "Manage objects",
		Long:  "Create, search, update, and delete Anytype objects",
	}

	cmd.AddCommand(objectCreateCmd.NewCreateCmd())
	cmd.AddCommand(objectSearchCmd.NewSearchCmd())
	cmd.AddCommand(objectGetCmd.NewGetCmd())
	cmd.AddCommand(objectUpdateCmd.NewUpdateCmd())
	cmd.AddCommand(objectDeleteCmd.NewDeleteCmd())

	return cmd
}
