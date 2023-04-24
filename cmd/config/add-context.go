package config

import (
	"fmt"

	"github.com/morty-faas/cli/cliconfig"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addContextCmd = &cobra.Command{
	Use:   "add-context NAME",
	Short: "Add a new context",
	Long:  `Add a new context to your configuration to allow easy interaction with a Morty instance.`,
	Args:  validateContextName,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Safe call, validation is performed by validateArgs automatically by cobra
		name := args[0]

		registry, _ := cmd.Flags().GetString("registry")
		gateway, _ := cmd.Flags().GetString("gateway")

		cfg := cmd.Context().Value(cliconfig.CtxKey{}).(*cliconfig.Config)

		ctx := &cliconfig.Context{
			Name:     name,
			Gateway:  gateway,
			Registry: registry,
		}

		log.Debugf("Adding context '%s' (gateway: %s, registry: %s)", name, gateway, registry)

		// We add the context to our configuration and we set it to the current context
		// before saving the configuration on disk.
		if err := cfg.AddContext(ctx); err != nil {
			return err
		}

		if err := cfg.UseContext(ctx.Name); err != nil {
			return err
		}

		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Printf("Success ! Your context '%s' has been saved and it is now the active context.\n", ctx.Name)

		return nil
	},
}

func init() {
	addContextCmd.PersistentFlags().String("gateway", "http://localhost:8080", "The URL of the Morty instance gateway.")
	addContextCmd.PersistentFlags().String("registry", "http://localhost:8081", "The URL of the Morty instance registry.")
}
