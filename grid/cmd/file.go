package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "List Files",
	Long:  "",
	Run:   GetExportFiles,
}

func init() {
	fileCmd.Flags().IntVarP(&pk, "pk", "", 0, "Primary key")
}

// GetExportFiles returns the user's exported data for a given AOI.
func GetExportFiles(cmd *cobra.Command, args []string) {
	tp := GetTransport()

	// github client configured to use test server
	client := grid.NewClient(tp.Client())
	a, resp, err := client.Export.ListByPk(pk)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err.Error())
	}
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 8, 3, '\t', 0)
	fmt.Fprintln(w, "NAME\tPRIMARY KEY")

	for _, v := range a {
		fmt.Fprintf(w, "%s\t%d\n", v.Name, v.Pk)
	}

	w.Flush()
}
