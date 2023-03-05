package cmd

import (
	"fmt"
	"morty/utils"
	"strings"

	"github.com/spf13/cobra"
)

// newCmd represents the build command
var newCmd = &cobra.Command{
	Use:   "new <NAME>",
	Short: "Create a new workspace to develop function",
	Long:  `This command creates a complete workspace with default configuration to write function code.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("you must specify a function name")
		}
		if args[0] != "" && !strings.Contains(args[0], " ") {
			return nil
		}
		return fmt.Errorf("invalid function name specified")
	},
	Run: func(cmd *cobra.Command, args []string) {
		runtime := cmd.Flags().Lookup("runtime").Value.String()
		utils.Runtime(runtime).CheckValidityOrExit()

		name := args[0]

		utils.New(name, runtime)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("runtime", "r", "", fmt.Sprintf("Runtime of the function : %s", utils.GetAvailableRuntimesAsString()))
	newCmd.MarkFlagRequired("runtime")
}
