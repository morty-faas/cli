package runtime

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "runtime",
		Aliases: []string{"rt"},
		Short: "Manage runtimes from the CLI",
		Long:  `Manage runtimes from the CLI`,
	}
)

func init() {
	RootCmd.AddCommand(listCmd)
}
