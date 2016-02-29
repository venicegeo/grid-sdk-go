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
	"github.com/venicegeo/grid-sdk-go"
)

var lsCmd = &cobra.Command{
	Use:   "ls [pk...]",
	Short: "List AOI/Export/File details",
	Long: `
List AOI, export, or file details for the provided primary keys.

With no keys specified, the command returns a listing of all of the user's
AOIs.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If there is no primary key provided, we just return a root level listing.
		if len(args) == 0 {
			result := getAOIs()
			printAOI(result.(*grid.AOIResponse))
		}

		// If the user has provided one or more arguments, assume they are primary
		// keys and concurrently query the AOI and export API endpoints for details.
		var results []interface{}
		for _, arg := range args {
			pk, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Error parsing \"%v\". Please provide primary keys as integers. Continuing with remaining keys...\n\n", arg)
				continue
			}

			c := make(chan interface{})
			go func() {
				exports, err := grid.GetAOI(pk)
				if err == nil {
					c <- exports
				} else {
					c <- err
				}
			}()
			go func() {
				exportFiles, err := grid.GetExport(pk)
				if err == nil {
					c <- exportFiles
				} else {
					c <- err
				}
			}()

			for i := 0; i < 2; i++ {
				result := <-c
				results = append(results, result)
			}
		}

		// Depending on the type of the returned objects, print it accordingly.
		for _, v := range results {
			switch u := v.(type) {
			case *grid.AOIItem:
				printExport(u)
			case *grid.ExportDetail:
				printExportFile(u)
			case error:
				// log.Print(u.Error())
			default:
				//NOOP fmt.Println("unknown")
			}
		}
	},
}

func getAOIs() interface{} {
	// get the full list of AOIs
	a, err := grid.ListAOIs("")
	if err != nil {
		log.Fatal(err.Error())
	}

	return a
}

func printAOI(a *grid.AOIResponse) {
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

func printExport(a *grid.AOIItem) {
	for inx := 0; inx < len(a.AOIs); inx++ {
		fmt.Println(a.AOIs[inx].Fields.ClipGeometry)
	}
	if len(a.ExportSet) > 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 3, '\t', 0)
		fmt.Fprintln(w, "PRIMARY KEY\tNAME\tDATATYPE\tSTARTED AT")
		for _, vv := range a.ExportSet {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", vv.Pk, vv.Name, vv.Datatype, vv.StartedAt)
		}
		w.Flush()
	}
}

func printExportFile(a *grid.ExportDetail) {
	if len(a.ExportFiles) > 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 3, '\t', 0)
		fmt.Fprintln(w, "PRIMARY KEY\tNAME")
		for _, vv := range a.ExportFiles {
			fmt.Fprintf(w, "%v\t%v\n", vv.Pk, vv.Name)
		}
		w.Flush()
	}
}
