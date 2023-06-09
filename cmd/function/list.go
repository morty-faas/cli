package function

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/morty-faas/cli/cliconfig"

	morty "github.com/morty-faas/morty/pkg/client/controller"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List functions",
	Long:    "List functions",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdContext := cmd.Context()
		client := cmdContext.Value(cliconfig.ControllerClientContextKey{}).(*morty.APIClient)

		functions, res, err := client.FunctionApi.GetFunctions(cmdContext).Execute()
		if err != nil {
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

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NAME", "VERSION"})
		for _, fn := range functions {
			for _, v := range fn.GetVersions() {
				table.Append([]string{fn.GetName(), v})
			}
		}

		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.SetNoWhiteSpace(true)
		table.Render()

		return nil
	},
}
