package relation

import (
	"github.com/spf13/cobra"

	relationAddToTypeCmd "github.com/anyproto/anytype-cli/cmd/relation/addtotype"
	relationCreateCmd "github.com/anyproto/anytype-cli/cmd/relation/create"
	relationListCmd "github.com/anyproto/anytype-cli/cmd/relation/list"
)

func NewRelationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <command>",
		Short: "Manage relations",
		Long:  "Create and list relations (object fields/properties)",
	}

	cmd.AddCommand(relationCreateCmd.NewCreateCmd())
	cmd.AddCommand(relationListCmd.NewListCmd())
	cmd.AddCommand(relationAddToTypeCmd.NewAddToTypeCmd())

	return cmd
}
