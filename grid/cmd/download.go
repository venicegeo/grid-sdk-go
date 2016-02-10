package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Download File",
	Long:  "",
	Run:   DownloadFile,
}

func init() {
	pullCmd.Flags().IntVarP(&pk, "pk", "", 0, "Primary key")
}

// DownloadFile downloads a file by pk.
func DownloadFile(cmd *cobra.Command, args []string) {
	tp := GetTransport()

	// github client configured to use test server
	client := grid.NewClient(tp.Client())
	resp, err := client.Export.DownloadByPk(pk)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err.Error())
	}
}
