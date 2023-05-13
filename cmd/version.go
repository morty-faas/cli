package cmd

import (
	"fmt"

	"github.com/morty-faas/cli/build"
	"github.com/morty-faas/cli/cliconfig"
	morty "github.com/morty-faas/morty/pkg/client/controller"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the current version of the CLI and of the Morty server of the current context.",
	Long:  "Get the current version of the CLI and of the Morty server of the current context.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmd.Context().Value(cliconfig.ControllerClientContextKey{}).(*morty.APIClient)

		metadata, _, err := client.ConfigurationApi.GetServerMetadata(cmd.Context()).Execute()
		if err != nil {
			return err
		}

		fmt.Printf("CLI    : '%s' (commit: %s)\n", build.Version, build.GitCommit)
		fmt.Printf("Server : '%s' (commit: %s)\n", metadata.GetVersion(), metadata.GetGitCommit())
		return nil
	},
}
