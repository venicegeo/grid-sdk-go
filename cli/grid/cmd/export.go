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
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [AOI] [Collects]...",
	Short: "Initiate a GRiD Export",
	Long: `
Export is used to initiate a GRiD export for the AOI and for each of the provided collects.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := initClient()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var collects []string
		switch len(args) {
		case 0:
			fmt.Println("Please provide an AOI")
			cmd.Usage()
			return
		case 1:
			fmt.Println("Please provide a collect.")
			cmd.Usage()
			return
		default:
			collects = args[1:]
		}
		pk, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error parsing \"%v\". Please provide AOI primary key as integer.\n\n", args[0])
			return
		}

		export, _, err := g.GeneratePointCloudExport(pk, collects, nil)
		if err != nil {
			log.Fatal(err)
		}
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 3, '\t', 0)
		fmt.Fprintln(w, "TASK ID\tEXPORT ID")
		fmt.Fprintf(w, "%v\t%v\n", export.TaskID, export.ExportID)
		w.Flush()
	},
}
