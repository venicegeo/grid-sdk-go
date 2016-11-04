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

var geom string

func init() {
	lsCmd.Flags().StringVarP(&geom, "geom", "", "", "WKT Polygon")
}

var lsCmd = &cobra.Command{
	Use:   "ls [pk...]",
	Short: "List AOI/Export/File details",
	Long: `
List AOI, export, or file details for the provided primary keys.

With no keys specified, the command returns a listing of all of the user's
AOIs.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := initClient()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		listAOIs := false
		if len(args) == 0 || geom != "" {
			listAOIs = true
		}
		// If there is no primary key provided, we just return a root level listing.
		if listAOIs {
			a := new(grid.AOIArray)
			if geom == "" {
				// get the full list of AOIs
				b, _, err := g.ListAOIs("")
				if err != nil {
					log.Fatal(err.Error())
				}
				a = b
			} else {
				// get the list of AOIs intersecting the geometry
				b, _, err := g.ListAOIs(geom)
				if err != nil {
					log.Fatal(err.Error())
				}
				a = b
			}

			// bi, _ := json.MarshalIndent(a.AOIList, "", "   ")
			// os.Stdout.Write(bi)

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 3, '\t', 0)
			fmt.Fprintln(w, "PRIMARY KEY\tNAME\tCREATED AT")
			for _, v := range a.AOIList {
				fmt.Fprintf(w, "%v\t%v\t%v\n", v.Pk, v.Name, v.CreatedAt)
			}
			w.Flush()
		}

		// If the user has provided one or more arguments, assume they are primary
		// keys and concurrently query the AOI and export API endpoints for details.
		// var results []interface{}
		for _, arg := range args {
			pk, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Error parsing \"%v\". Please provide primary keys as integers.\n\n", arg) // Continuing with remaining keys...\n\n", arg)
				continue
			}

			c1 := make(chan *grid.AOIDetail)
			c2 := make(chan *grid.ExportDetail)
			// unsure about how to handle errors anymore, as we expect users to be able to supply any type of pk, but GRiD now returns an error if you provide an export pk to the aoi endpoint and vice versa (plus product pks)
			go func() {
				// get information on the AOI specified by the given primary key
				a, _, _ := g.GetAOI(pk)
				// if err != nil {
				// 	log.Fatal(err.Error())
				// }

				c1 <- a
			}()
			go func() {
				// get information on the export specified by the given primary key
				a, _, _ := g.GetExport(pk)
				// if err != nil {
				// 	log.Fatal(err.Error())
				// }

				c2 <- a
			}()

			for i := 0; i < 2; i++ {
				select {
				case a := <-c1:
					fmt.Println()
					fmt.Println("NAME:", a.Name)
					fmt.Println("CREATED AT:", a.CreatedAt)
					fmt.Println("\nRASTER COLLECTS")
					if len(a.RasterIntersects) > 0 {
						w := new(tabwriter.Writer)
						w.Init(os.Stdout, 0, 8, 3, '\t', 0)
						fmt.Fprintln(w, "PRIMARY KEY\tNAME\tDATATYPE")
						for _, vv := range a.RasterIntersects {
							fmt.Fprintf(w, "%v\t%v\t%v\n", vv.Pk, vv.Name, vv.Datatype)
						}
						w.Flush()
					}
					fmt.Println("\nPOINTCLOUD COLLECTS")
					if len(a.PointcloudIntersects) > 0 {
						w := new(tabwriter.Writer)
						w.Init(os.Stdout, 0, 8, 3, '\t', 0)
						fmt.Fprintln(w, "PRIMARY KEY\tNAME\tDATATYPE")
						for _, vv := range a.PointcloudIntersects {
							fmt.Fprintf(w, "%v\t%v\t%v\n", vv.Pk, vv.Name, vv.Datatype)
						}
						w.Flush()
					}
					fmt.Println("\nEXPORTS")
					if len(a.ExportSet) > 0 {
						w := new(tabwriter.Writer)
						w.Init(os.Stdout, 0, 8, 3, '\t', 0)
						fmt.Fprintln(w, "PRIMARY KEY\tNAME\tDATATYPE\tSTARTED AT")
						for _, vv := range a.ExportSet {
							fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", vv.Pk, vv.Name, vv.Datatype, vv.StartedAt)
						}
						w.Flush()
					}
				case b := <-c2:
					if len(b.ExportFiles) > 0 {
						fmt.Println()
						w := new(tabwriter.Writer)
						w.Init(os.Stdout, 0, 8, 3, '\t', 0)
						fmt.Fprintln(w, "PRIMARY KEY\tNAME")
						for _, vv := range b.ExportFiles {
							fmt.Fprintf(w, "%v\t%v\n", vv.Pk, vv.Name)
						}
						w.Flush()
					}
				}
			}
		}
	},
}
