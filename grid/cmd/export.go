package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "List Exports",
	Long:  "",
	Run:   GetExports,
}

func init() {
	exportCmd.Flags().IntVarP(&pk, "pk", "", 0, "Primary key")
}

// GetExports returns the user's exports for a given AOI.
func GetExports(cmd *cobra.Command, args []string) {
	tp := GetTransport()

	// github client configured to use test server
	client := grid.NewClient(tp.Client())
	a, resp, err := client.AOI.ListByPk(pk)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err.Error())
	}
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 8, 3, '\t', 0)
	fmt.Fprintln(w, "STARTED AT\tNAME\tDATATYPE\tPRIMARY KEY")

	for _, v := range a {
		fmt.Fprintf(w, "%v\t%s\t%s\t%d\n", v.StartedAt, v.Name, v.Datatype, v.Pk)
	}

	w.Flush()
}
