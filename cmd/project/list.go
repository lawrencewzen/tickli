package project

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/sho0pi/tickli/internal/api"
	"github.com/sho0pi/tickli/internal/types"
	"github.com/spf13/cobra"
)

type listOptions struct {
	filter string
}

func filterProjectByName(projects []types.Project, name string) ([]types.Project, error) {
	var matched []types.Project
	nameLower := strings.ToLower(name)
	for i := range projects {
		if strings.Contains(strings.ToLower(projects[i].Name), nameLower) {
			matched = append(matched, projects[i])
		}
	}
	if len(matched) == 0 {
		return nil, fmt.Errorf("no project found with name '%s'", name)
	}
	return matched, nil
}

func newListCommand(client *api.Client) *cobra.Command {
	opts := &listOptions{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available projects",
		Long: `Print all available projects, one per line, optionally filtered by name.

The output is plain text (ID and name columns) so it can be piped into other
commands or scripts.`,
		Example: `  # List all projects
  tickli project list

  # Filter projects by name
  tickli project list -f "work"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := client.ListProjects()
			if err != nil {
				return errors.Wrap(err, "failed to fetch projects")
			}

			if opts.filter != "" {
				projects, err = filterProjectByName(projects, opts.filter)
				if err != nil {
					return err
				}
			}

			if len(projects) == 0 {
				fmt.Println("No projects found.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME")
			for i := range projects {
				fmt.Fprintf(w, "%s\t%s\n", projects[i].ID, projects[i].Name)
			}
			return w.Flush()
		},
	}

	cmd.Flags().StringVarP(&opts.filter, "filter", "f", "", "Only show projects with names containing the provided text")

	return cmd
}
