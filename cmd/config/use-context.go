package config

import (
	"fmt"
	"morty/cliconfig"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var useContextCmd = &cobra.Command{
	Use:   "use-context",
	Short: "Update your current context",
	Long:  `Update your current context with one of the contexts available in your configuration`,
	Args:  validateContextName,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := cmd.Context().Value(cliconfig.CtxKey{}).(*cliconfig.Config)

		name := args[0]

		log.Debugf("Updating active context to '%s'", name)

		if err := cfg.UseContext(name); err != nil {
			return err
		}

		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Printf("Your current context has been set to '%s'.\n", name)

		return nil
	},
}
