package config

import (
	"fmt"

	"github.com/morty-faas/cli/cliconfig"

	"github.com/spf13/cobra"
)

var currentContextCmd = &cobra.Command{
	Use:   "current-context",
	Short: "Display information about your current context",
	Long:  `Display information about your current context`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := cmd.Context().Value(cliconfig.CtxKey{}).(*cliconfig.Config)

		ctx, err := cfg.GetCurrentContext()
		if err != nil {
			return err
		}

		fmt.Printf("Name         : %s\n", ctx.Name)
		fmt.Printf("Gateway URL  : %s\n", ctx.Gateway)
		fmt.Printf("Registry URL : %s\n", ctx.Registry)

		return nil
	},
}
