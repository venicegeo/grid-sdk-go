package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Get suggested AOI name",
	Long:  "",
	Run:   RunLookup,
}

func RunLookup(cmd *cobra.Command, args []string) {
	// geom := "POLYGON((30 10,40 40,20 40,10 20,30 10))"
	geom := args[0]
	tp := GetTransport()

	// github client configured to use test server
	client := grid.NewClient(tp.Client())
	a, _, err := client.Geonames.Lookup(geom)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(a.Name)
}
