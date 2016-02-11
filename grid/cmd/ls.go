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
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List AOI/Export/File details",
	Long:  "",
	Run:   RunLs,
}

// GetAOIs returns the user's AOIs.
func getAOIs() interface{} {
	tp := GetTransport()

	// github client configured to use test server
	client := grid.NewClient(tp.Client())
	a, _, err := client.AOI.List("")
	if err != nil {
		fmt.Println(err.Error())
	}

	return a
}

// GetExports returns the user's exports for a given AOI.
func getExports(pk int) interface{} {
	tp := GetTransport()

	// github client configured to use test server
	client := grid.NewClient(tp.Client())
	a, _, err := client.AOI.ListByPk(pk)
	if err != nil {
		fmt.Println(err.Error())
	}

	return a
}

// GetExportFiles returns the user's exported data for a given AOI.
func getExportFiles(pk int) interface{} {
	tp := GetTransport()

	// github client configured to use test server
	client := grid.NewClient(tp.Client())
	a, _, err := client.Export.ListByPk(pk)
	if err != nil {
		fmt.Println(err.Error())
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

func printExportFile(a *grid.ExportResponse) {
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

func RunLs(cmd *cobra.Command, args []string) {
	var pk int
	var results []interface{}

	// If the user has provided an argument, assume it is the primary key and
	// concurrently query the AOI and export API endpoints for details. If there
	// is no primary key, we just return a root level listing.
	if len(args) > 0 {
		pk, _ = strconv.Atoi(args[0])
		c := make(chan interface{})
		go func() { c <- getExports(pk) }()
		go func() { c <- getExportFiles(pk) }()

		for i := 0; i < 2; i++ {
			result := <-c
			results = append(results, result)
		}
	} else {
		result := getAOIs()
		results = append(results, result)
	}

	// Depending on the type of the returned objects, print it accordingly.
	for _, v := range results {
		switch u := v.(type) {
		case *grid.AOIItem:
			printExport(u)
		case *grid.ExportResponse:
			printExportFile(u)
		case *grid.AOIResponse:
			printAOI(u)
		default:
			fmt.Println("unknown")
		}
	}
}
