package function

import (
	"fmt"
	"morty/cliconfig"
	"morty/function"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [PATH]",
	Short: "Build a function image to be run in Morty. If PATH is not provided, the current directory will be used.",
	Long:  `This command allows you to build a function image that can be run in Morty.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context().Value(cliconfig.CurrentCtxKey{}).(*cliconfig.Context)

		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = fmt.Sprintf("Building function from %s ", path)
		s.Suffix = "\n"
		s.Start()

		opts := &function.BuildOptions{
			Directory: path,
			Registry:  ctx.Registry,
			Gateway:   ctx.Gateway,
		}
		name, err := function.Build(opts)
		s.Stop()
		if err != nil {
			return err
		}
		fmt.Printf("Function %s has been created !\n", name)

		return nil
	},
}

func init() {}
