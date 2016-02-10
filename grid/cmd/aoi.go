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
	fmt.Fprintln(w, "PRIMARY KEY\tNAME\tCREATED AT")

	for _, v := range *a {
		for _, v := range v.AOIs {
			fmt.Fprintf(w, "%v\t%v\t%v\n", v.Pk, v.Fields.Name, v.Fields.CreatedAt)
		}
	}

	w.Flush()
}
