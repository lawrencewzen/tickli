package project

import (
	"github.com/sho0pi/tickli/internal/api"
	"github.com/sho0pi/tickli/internal/utils"
	"github.com/spf13/cobra"
)

// NewProjectCommand returns a cobra command for `project` subcommands
func NewProjectCommand() *cobra.Command {
	var client api.Client
	cmd := &cobra.Command{
		Use:     "project",
		Short:   "Work with TickTick projects.",
		Aliases: []string{"list"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			client = utils.LoadClient()
			return nil
		},
	}

	cmd.AddCommand(
		newListCommand(&client),
		newUseProjectCmd(&client),
	)

	return cmd
}
