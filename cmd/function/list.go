package function

import (
	"morty/cliconfig"
	"morty/client/gateway"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List functions",
	Long:    "List functions",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := cmd.Context().Value(cliconfig.CurrentCtxKey{}).(*cliconfig.Context)

		gw := gateway.NewClient(ctx.Gateway)

		opts := &gateway.ListFnRequest{}

		functions, err := gw.ListFn(cmd.Context(), opts)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "id"})
		for k, v := range functions {
			table.Append([]string{k, v})
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
