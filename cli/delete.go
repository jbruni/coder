package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/coder/coder/cli/cliui"
	"github.com/coder/coder/codersdk"
)

// nolint
func deleteWorkspace() *cobra.Command {
	var orphan bool
	cmd := &cobra.Command{
		Annotations: workspaceCommand,
		Use:         "delete <workspace>",
		Short:       "Delete a workspace",
		Aliases:     []string{"rm"},
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cliui.Prompt(cmd, cliui.PromptOptions{
				Text:      "Confirm delete workspace?",
				IsConfirm: true,
				Default:   cliui.ConfirmNo,
			})
			if err != nil {
				return err
			}

			client, err := CreateClient(cmd)
			if err != nil {
				return err
			}
			workspace, err := namedWorkspace(cmd, client, args[0])
			if err != nil {
				return err
			}

			var state []byte

			if orphan {
				cliui.Warn(cmd.ErrOrStderr(), "Orphaning workspace",
					"Template edit permission is required to orphan workspaces.",
				)

				state, err = client.WorkspaceBuildState(cmd.Context(), workspace.LatestBuild.ID)
				if err != nil {
					return err
				}
				// If there's alreay no state, orphanage makes no sense.
				if len(state) > 0 {
					state, err = codersdk.OrphanTerraformState(state)
					if err != nil {
						return err
					}
				}
			}

			before := time.Now()
			build, err := client.CreateWorkspaceBuild(cmd.Context(), workspace.ID, codersdk.CreateWorkspaceBuildRequest{
				Transition:       codersdk.WorkspaceTransitionDelete,
				ProvisionerState: state,
			})
			if err != nil {
				return err
			}

			err = cliui.WorkspaceBuild(cmd.Context(), cmd.OutOrStdout(), client, build.ID, before)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nThe %s workspace has been deleted at %s!\n", cliui.Styles.Keyword.Render(workspace.Name), cliui.Styles.DateTimeStamp.Render(time.Now().Format(time.Stamp)))
			return nil
		},
	}
	cmd.Flags().BoolVar(&orphan, "orphan", false,
		`Delete workspace without deleting its resources. This can delete a
workspace in a broken state, but may also lead to unaccounted cloud resources.`,
	)
	cliui.AllowSkipPrompt(cmd)
	return cmd
}
