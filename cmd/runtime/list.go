package runtime

import (
	"fmt"

	"github.com/morty-faas/cli/runtime"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available function runtimes.",
	Long:    `List available function runtimes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		runtimes, err := runtime.List()
		if err != nil {
			return err
		}

		fmt.Println("Available runtimes:")
		for _, runtime := range runtimes {
			fmt.Printf("- %s\n", runtime)
		}
		return nil
	},
}
