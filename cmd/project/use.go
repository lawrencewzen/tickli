package project

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sho0pi/tickli/internal/api"
	"github.com/sho0pi/tickli/internal/completion"
	"github.com/sho0pi/tickli/internal/config"
	"github.com/sho0pi/tickli/internal/types"
	"github.com/spf13/cobra"
)

// resolveProject finds a single project matching query, which may be either an
// exact project ID or a case-insensitive substring of the project name.
func resolveProject(projects []types.Project, query string) (types.Project, error) {
	// Exact ID match takes precedence.
	for i := range projects {
		if projects[i].ID == query {
			return projects[i], nil
		}
	}

	// Fall back to case-insensitive name matching.
	var matched []types.Project
	q := strings.ToLower(query)
	for i := range projects {
		if strings.Contains(strings.ToLower(projects[i].Name), q) {
			matched = append(matched, projects[i])
		}
	}

	switch len(matched) {
	case 0:
		return types.Project{}, fmt.Errorf("no project matches '%s'", query)
	case 1:
		return matched[0], nil
	default:
		lines := make([]string, len(matched))
		for i := range matched {
			lines[i] = fmt.Sprintf("  %s (%s)", matched[i].Name, matched[i].ID)
		}
		return types.Project{}, fmt.Errorf("'%s' matches multiple projects, be more specific:\n%s",
			query, strings.Join(lines, "\n"))
	}
}

type useProjectOptions struct {
	projectID string
}

func newUseProjectCmd(client *api.Client) *cobra.Command {
	opts := &useProjectOptions{}
	cmd := &cobra.Command{
		Use:   "use <project>",
		Short: "Switch active project context",
		Long: `Switch the active project context for subsequent commands.

Pass a project name (partial, case-insensitive) or an exact project ID.
The selected project becomes the default context for future commands.`,
		Example: `  # Switch by project name (partial match)
  tickli project use "Work Tasks"

  # Switch by exact project ID
  tickli project use abc123def456`,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: completion.ProjectIDs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				opts.projectID = args[0]
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.projectID == "" {
				return errors.New("a project name or ID is required, e.g. `tickli project use \"Work Tasks\"`")
			}

			projects, err := client.ListProjects()
			if err != nil {
				return errors.Wrap(err, "could not fetch projects")
			}

			selectedProject, err := resolveProject(projects, opts.projectID)
			if err != nil {
				return err
			}

			cfg, err := config.Load()
			if err != nil {
				return errors.Wrap(err, "could not load config")
			}

			cfg.DefaultProjectID = selectedProject.ID
			if err := config.Save(cfg); err != nil {
				return errors.Wrap(err, "failed to save config")
			}
			log.Info().
				Str("project_id", cfg.DefaultProjectID).
				Str("project_name", selectedProject.Name).
				Msg("Switched to project")
			return nil
		},
	}

	return cmd
}
