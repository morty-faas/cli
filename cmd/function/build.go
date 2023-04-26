package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/morty-faas/cli/cliconfig"
	"github.com/morty-faas/cli/function"
	"github.com/morty-faas/cli/pkg/archive"

	"github.com/briandowns/spinner"
	morty "github.com/morty-faas/morty/pkg/client/controller"
	registry "github.com/morty-faas/morty/pkg/client/registry"
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
		client := cmdContext.Value(cliconfig.ControllerClientContextKey{}).(*morty.APIClient)
		registry := cmdContext.Value(cliconfig.RegistryClientContextKey{}).(*registry.APIClient)

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
		archive, err := os.Open(zipFile)
		if err != nil {
			return err
		}

		// Ask the registry to build the function

		s.Start()

		fnUri, _, err := registry.FunctionsApi.V1FunctionsBuildPost(cmdContext).Name(f.Name).Runtime(f.Runtime).Archive(archive).Execute()
		if err != nil {
			return err
		}

		// Temporary awaiting the fix on the registry
		fnUri = strings.Replace(fnUri, "\"", "", -1)

		// Create the function to be able to invoke it
		image := ctx.Registry + fnUri
		createFnRequest := morty.CreateFunctionRequest{
			Name:  &f.Name,
			Image: &image,
		}

		if _, res, err := client.FunctionApi.CreateFunction(cmdContext).CreateFunctionRequest(createFnRequest).Execute(); err != nil {
			// If an error is returned by the API, parse it
			if res != nil && res.Body != nil {
				apiError := &morty.Error{}
				if err := json.NewDecoder(res.Body).Decode(apiError); err != nil {
					return err
				}
				return errors.New(apiError.GetMessage())
			}
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
