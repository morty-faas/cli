package config

import (
	"fmt"
	"github.com/morty-faas/cli/cliconfig"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var removeContextCmd = &cobra.Command{
	Use:   "remove-context NAME",
	Short: "Remove a context",
	Long:  `Remove a context from your configuration.`,
	Args:  validateContextName,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Safe call, validation is performed by validateArgs automatically by cobra
		name := args[0]

		cfg := cmd.Context().Value(cliconfig.CtxKey{}).(*cliconfig.Config)

		log.Debugf("Remove context '%s'", name)

		currentContext, _ := cfg.GetCurrentContext()

		if err := cfg.RemoveContext(name); err != nil {
			return err
		}

		// if we delete the current context, we set the first context as the current context
		if currentContext.Name == name {
			if err := cfg.UseContext(cfg.Contexts[0].Name); err != nil {
				return err
			}
		}

		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Printf("Success ! Your context '%s' has been deleted.\n", name)

		if currentContext.Name == name {
			fmt.Printf("Your current context has been set to '%s'.\n", cfg.Contexts[0].Name)
		}

		return nil
	},
}
