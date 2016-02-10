package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var aoiCmd = &cobra.Command{
	Use:   "aoi",
	Short: "List AOIs",
	Long:  "",
	Run:   GetAOIs,
}

// GetAOIs returns the user's AOIs.
func GetAOIs(cmd *cobra.Command, args []string) {
	tp := GetTransport()

	// github client configured to use test server
	client := grid.NewClient(tp.Client())
	a, resp, err := client.AOI.List("")
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err.Error())
	}
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 8, 3, '\t', 0)
	fmt.Fprintln(w, "CREATED AT\tPRIMARY KEY\tNAME\tNUMBER OF EXPORTS")

	for _, v := range a {
		fmt.Fprintf(w, "%v\t%d\t%s\t%s\n", v.CreatedAt, v.Pk, v.Name, v.NumExports)
	}

	w.Flush()
}
