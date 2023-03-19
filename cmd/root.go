/*
Copyright © 2023 polyxia-org
*/
package cmd

import (
	"context"
	"log"
	"morty/cliconfig"
	"morty/cmd/config"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "morty",
	Short: "Morty allows you to manage and invoke functions over Polyxia.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg, err := cliconfig.Load()
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.WithValue(cmd.Context(), cliconfig.CtxKey, cfg)
		cmd.SetContext(ctx)
	},
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(config.RootCmd)
}
