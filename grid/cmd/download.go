// Copyright 2016, RadiantBlue Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
