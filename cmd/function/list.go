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

		functions, err := gw.ListFn(cmd.Context(), &gateway.ListFnRequest{})
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "NAME", "IMAGE"})
		for _, fn := range *functions {
			table.Append([]string{fn.Id, fn.Name, fn.ImageURL})
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
