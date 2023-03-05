package cmd

import (
	"fmt"
	"log"
	"morty/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
    Use:   "build <NAME>",
    Short: "Build a rootfs to be run in morty FaaS",
    Long:  `This command allow you to package a function into a rootfs that can be run in morty FaaS.`,
    Args:  func(cmd *cobra.Command, args []string) error {
    if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
            return fmt.Errorf("you must specify a function name")
    }
    if args[0] != "" && !strings.Contains(args[0], " ") {
        return nil
    }
    return fmt.Errorf("invalid function name specified")
    },
    Run: func(cmd *cobra.Command, args []string) {
        flags := cmd.Flags()
        name := args[0]
        buildArgs, _ := flags.GetStringArray("build-arg")

        path, _ := flags.GetString("path")
        if path != "" { // if path is provided, use manual build
            // check if runtime is provided and valid
            runtime, err := flags.GetString("runtime")
            if err != nil {
                log.Fatal("Error getting runtime flag: ", err)
            }
            utils.Runtime(runtime).CheckValidityOrExit()

            // check if path exists
            if _, err := os.Stat(path); os.IsNotExist(err) {
                log.Fatal("ERROR: path provided for code folder does not exists")
            }
            utils.Build(name, runtime, path, buildArgs)
        } else {  // if path is not provided, use intuitive build
            utils.BuildIntuitive(name, buildArgs)
        }
    },
}

func init() {
    rootCmd.AddCommand(buildCmd)

    buildCmd.Flags().StringP("path", "p", "", "path of the function to build")
	buildCmd.Flags().StringP("runtime", "r", "", fmt.Sprintf("Runtime of the function : %s", utils.GetAvailableRuntimesAsString()))
    buildCmd.Flags().StringArrayP("build-arg", "b", []string{}, "Add a build-arg for Docker (KEY=VALUE)")

}