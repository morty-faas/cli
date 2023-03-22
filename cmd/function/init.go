package function

import (
	"fmt"
	"morty/function"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "init NAME <opt:DIRECTORY>",
	Short: "Initialize a new function workspace",
	Long:  `Initialize a new function workspace`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// By default, directory is optional so use the name of the function instead
		// If no directory is specified, the function will be created into the current working
		// directory of the user.
		dir := name
		if len(args) == 2 {
			dir = args[1]
		}

		runtime, _ := cmd.Flags().GetString("runtime")

		opts := &function.Options{
			Name:      name,
			Runtime:   runtime,
			Directory: dir,
		}

		if _, err := function.New(opts); err != nil {
			return fmt.Errorf("failed to initialize function: %v", err)
		}

		fmt.Printf("Success ! Your function '%s' has been initialized into '%s'\n", opts.Name, opts.Directory)
		return nil
	},
}

func init() {
	createCmd.PersistentFlags().StringP("runtime", "r", "node-19", "The runtime to use for the function.")
}
