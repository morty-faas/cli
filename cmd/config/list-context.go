package config

import (
	"fmt"

	"github.com/morty-faas/cli/cliconfig"

	"github.com/spf13/cobra"
)

var listContextCmd = &cobra.Command{
	Use:   "contexts",
	Short: "List all contexts available in the configuration",
	Long:  `List all contexts available in the configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmd.Context().Value(cliconfig.CtxKey{}).(*cliconfig.Config)

		for i, c := range cfg.Contexts {
			fmt.Printf("Name            : %s\n", c.Name)
			fmt.Printf("Controller URL  : %s\n", c.Controller)
			fmt.Printf("Registry URL    : %s\n", c.Registry)

			if i < len(cfg.Contexts)-1 {
				fmt.Printf("\n")
			}
		}
	},
}
