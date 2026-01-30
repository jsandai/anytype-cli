package object

import (
	"github.com/spf13/cobra"

	objectCreateCmd "github.com/anyproto/anytype-cli/cmd/object/create"
	objectDeleteCmd "github.com/anyproto/anytype-cli/cmd/object/delete"
	objectGetCmd "github.com/anyproto/anytype-cli/cmd/object/get"
	objectGetContentCmd "github.com/anyproto/anytype-cli/cmd/object/getcontent"
	objectSearchCmd "github.com/anyproto/anytype-cli/cmd/object/search"
	objectSetContentCmd "github.com/anyproto/anytype-cli/cmd/object/setcontent"
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
	cmd.AddCommand(objectGetContentCmd.NewGetContentCmd())
	cmd.AddCommand(objectUpdateCmd.NewUpdateCmd())
	cmd.AddCommand(objectDeleteCmd.NewDeleteCmd())
	cmd.AddCommand(objectSetContentCmd.NewSetContentCmd())

	return cmd
}
