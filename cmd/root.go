/*
Copyright © 2023 polyxia-org
*/
package cmd

import (
	"context"
	"io"
	"log"
	"morty/cliconfig"
	"morty/cmd/config"
	"morty/cmd/function"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	logEnvVarKey = "MORTY_LOG"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "morty",
	Short: "Morty CLI is used to interact with the Morty serverless platform.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// By default, if there is no value for the environment variable
		// MORTY_LOG, all logs are disabled.
		//
		// The program will produce only
		// outputs that comes from fmt.Print instructions.
		// But in some cases, it would be great to have the opportunity
		// to see what is going on in the execution of the program.
		//
		// If a value is provided, a logger will be configured. If the value
		// can't be parsed, the default log level will be applied, INFO.
		loglevel := os.Getenv(logEnvVarKey)
		if loglevel != "" {
			lvl, err := logrus.ParseLevel(loglevel)
			if err != nil {
				lvl = logrus.InfoLevel
			}
			logrus.SetLevel(lvl)
		} else {
			// Remove the entire output of logrus.* calls
			logrus.SetOutput(io.Discard)
		}

		// Add the configuration into the root context, so all commands can get it
		// and the configuration is loaded only once in the code.
		cfg, err := cliconfig.Load()
		if err != nil {
			log.Fatal(err)
		}

		logrus.Infof("Active context : %s", cfg.Current)

		ctx := context.WithValue(cmd.Context(), cliconfig.CtxKey, cfg)
		cmd.SetContext(ctx)
	},
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
	rootCmd.AddCommand(function.RootCmd)
}
