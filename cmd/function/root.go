package function

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:     "function",
		Aliases: []string{"fn"},
		Short:   "Manage functions from the CLI",
		Long:    `Manage functions from the CLI`,
	}
)

func init() {
	RootCmd.AddCommand(createCmd)
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(invokeCmd)
}
