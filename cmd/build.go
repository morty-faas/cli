package cmd

import (
	"log"
	"morty/utils"
	"os"

	"github.com/spf13/cobra"
)


// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build <PATH>",
	Short: "Build a rootfs to be run in morty FaaS",
	Long:  `This command allow you to package a function into a rootfs that can be run in morty FaaS.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flags().Lookup("name").Value.String()
		runtime := cmd.Flags().Lookup("runtime").Value.String()
		utils.Runtime(runtime).CheckValidityOrExit()

		folder := args[0]
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			log.Fatal("ERROR: path provided for code folder does not exists")
		}
		if stat, err := os.Stat(folder); err == nil && !stat.IsDir() {
			log.Fatal("ERROR: path provided for code folder is not a directory")
		}
		buildArgs, err := cmd.Flags().GetStringArray("build-arg")
		if err != nil {
			log.Fatal("Error getting build-arg flag:", err)
		}
		utils.Build(name, runtime, folder, buildArgs)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("name", "n", "", "Name of the function to build")
	buildCmd.MarkFlagRequired("name")
	buildCmd.Flags().StringP("runtime", "r", "", "Runtime of the function e.g. \"python-3\", \"node-19\"")
	buildCmd.MarkFlagRequired("runtime")
	buildCmd.Flags().StringArrayP("build-arg", "b", []string{}, "Add a build-arg for Docker (KEY=VALUE)")

}
