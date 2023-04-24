package config

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	ErrContextNameRequired = errors.New("you must provide a valid context name")
	RootCmd                = &cobra.Command{
		Use:   "config",
		Short: "Manage Morty CLI configuration",
		Long:  `Manage Morty CLI configuration, update your contexts etc.`,
	}
)

func init() {
	RootCmd.AddCommand(addContextCmd)
	RootCmd.AddCommand(currentContextCmd)
	RootCmd.AddCommand(useContextCmd)
	RootCmd.AddCommand(listContextCmd)
	RootCmd.AddCommand(removeContextCmd)
}

func validateContextName(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return ErrContextNameRequired
	}
	if args[0] == "" {
		return ErrContextNameRequired
	}
	return nil
}
