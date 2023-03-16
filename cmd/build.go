package cmd

import (
	"morty/utils"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [PATH]",
	Short: "Build a rootfs to be run in morty FaaS. If PATH is not provided, the current directory will be used.",
	Long:  `This command allow you to package a function into a rootfs that can be run in morty FaaS.`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		buildArgs, _ := flags.GetStringArray("build-arg")

		var path string
		if len(args) > 0 {
			path = args[0]
		} else {
			path = "."
		}

		utils.Build(path, buildArgs)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringArrayP("build-arg", "b", []string{}, "Add a build-arg for Docker (KEY=VALUE)")
}
