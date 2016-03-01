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
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var taskCmd = &cobra.Command{
	Use:   "task [Task ID]...",
	Short: "Get task details",
	Long: `
Lookup is used to retrieve the details of a GRiD task, including the status.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a WKT geometry")
			cmd.Usage()
			return
		}

		for _, taskID := range args {
			// get and print the suggested name for the current geometry
			task, err := grid.TaskDetails(taskID)
			if err != nil {
				log.Fatal(err)
			}
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 3, '\t', 0)
			fmt.Fprintln(w, "ID\tNAME\tSTATE")
			fmt.Fprintf(w, "%v\t%v\t%v\n", task.TaskID, task.Name, task.State)
			w.Flush()
		}
	},
}
