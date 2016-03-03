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
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Download File",
	Long: `
Download the file(s) specified by the given primary key(s).`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO(chambbj): allow for multiple pks/downloads
		if len(args) == 0 {
			fmt.Println("Please provide a primary key for the file to download.")
			cmd.Usage()
			return
		}

		pk, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err.Error())
		}
		dc, err := grid.DownloadByPk(pk)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("Downloaded %v", dc.FileName)
	},
}
