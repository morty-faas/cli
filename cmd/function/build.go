package function

import (
	"fmt"
	"morty/cliconfig"
	"morty/client/gateway"
	"morty/client/registry"
	"morty/function"
	"morty/pkg/archive"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

const (
	zipFile = "function.zip"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [PATH]",
	Short: "Build a function image to be run in Morty. If PATH is not provided, the current directory will be used.",
	Long:  `This command allows you to build a function image that can be run in Morty.`,
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return os.Remove(zipFile)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdContext := cmd.Context()

		ctx := cmdContext.Value(cliconfig.CurrentCtxKey{}).(*cliconfig.Context)

		// Initialize clients
		gw := gateway.NewClient(ctx.Gateway)
		reg := registry.NewClient(ctx.Registry)

		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		// Load the function metadata from the morty.yaml file
		f, err := function.NewFromFile(path)
		if err != nil {
			return err
		}

		// Zip the function code in the current working directory
		if err := archive.Zip(f.Path, zipFile); err != nil {
			return err
		}

		s := makeSpinner(fmt.Sprintf("Building your function '%s'", f.Name))

		// Ask the registry to build the function
		buildFnRequest := &registry.BuildFnRequest{
			Name:    f.Name,
			Runtime: f.Runtime,
			Archive: zipFile,
		}

		s.Start()
		fnUri, err := reg.BuildFn(cmdContext, buildFnRequest)
		if err != nil {
			return err
		}

		// Create the function to be able to invoke it
		createFnRequest := &gateway.CreateFnRequest{
			Name:   f.Name,
			Rootfs: ctx.Registry + fnUri,
		}

		if _, err := gw.CreateFn(cmdContext, createFnRequest); err != nil {
			return err
		}
		s.Stop()

		fmt.Printf("Success ! Your function '%s' has been created !\n", f.Name)

		return nil
	},
}

func makeSpinner(prefix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = prefix + " "
	s.Suffix = "\n"
	return s
}
