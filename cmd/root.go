package cmd

import (
	"context"
	"io"
	"os"

	"github.com/morty-faas/cli/cliconfig"
	"github.com/morty-faas/cli/cmd/config"
	"github.com/morty-faas/cli/cmd/function"
	"github.com/morty-faas/cli/cmd/runtime"

	morty "github.com/morty-faas/controller/pkg/client"
	registry "github.com/morty-faas/registry/pkg/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	logEnvVarKey = "MORTY_LOG"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
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
		var level log.Level

		envFlag := os.Getenv(logEnvVarKey)
		cliFlag, _ := cmd.Flags().GetCount("verbose")

		switch cliFlag {
		case 1:
			level = log.InfoLevel
		case 2:
			level = log.DebugLevel
		case 3:
			level = log.TraceLevel
		default:
			level, _ = log.ParseLevel(envFlag)
		}

		if level == 0 {
			log.SetOutput(io.Discard)
		} else {
			log.SetLevel(level)
		}

		// Add the configuration into the root context, so all commands can get it
		// and the configuration is loaded only once in the code.
		cfg, err := cliconfig.Load()
		if err != nil {
			log.Fatal(err)
		}

		currentCtx, err := cfg.GetCurrentContext()
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Active context : %s", cfg.Current)

		ctx := context.WithValue(cmd.Context(), cliconfig.CtxKey{}, cfg)
		ctx = context.WithValue(ctx, cliconfig.CurrentCtxKey{}, currentCtx)
		ctx = context.WithValue(ctx, cliconfig.ControllerClientContextKey{}, makeMortyClient(currentCtx.Controller))
		ctx = context.WithValue(ctx, cliconfig.RegistryClientContextKey{}, makeRegistryClient(currentCtx.Registry))

		cmd.SetContext(ctx)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(config.RootCmd)
	RootCmd.AddCommand(function.RootCmd)
	RootCmd.AddCommand(runtime.RootCmd)

	RootCmd.PersistentFlags().CountP("verbose", "v", "Level of verbosity: -v for INFO, -vv for DEBUG, -vvv for TRACE.")
}

func makeMortyClient(baseURL string) *morty.APIClient {
	return morty.NewAPIClient(&morty.Configuration{
		Servers: morty.ServerConfigurations{
			morty.ServerConfiguration{URL: baseURL},
		},
	})
}

func makeRegistryClient(baseURL string) *registry.APIClient {
	return registry.NewAPIClient(&registry.Configuration{
		Servers: registry.ServerConfigurations{
			registry.ServerConfiguration{URL: baseURL},
		},
	})
}
