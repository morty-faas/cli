package config

import (
	"fmt"
	"morty/cliconfig"

	"github.com/spf13/cobra"
)

var listContextCmd = &cobra.Command{
	Use:   "contexts",
	Short: "List all contexts available in the configuration",
	Long:  `List all contexts available in the configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmd.Context().Value(cliconfig.CtxKey{}).(*cliconfig.Config)

		for i, c := range cfg.Contexts {
			fmt.Printf("Name         : %s\n", c.Name)
			fmt.Printf("Gateway URL  : %s\n", c.Gateway)
			fmt.Printf("Registry URL : %s\n", c.Registry)

			if i < len(cfg.Contexts)-1 {
				fmt.Printf("\n")
			}
		}
	},
}
