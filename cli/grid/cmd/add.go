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
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var addCmd = &cobra.Command{
	Use:   "add [WKT geometry]...",
	Short: "Add an AOI",
	Long: `
Attempt to create new Areas of Interest (AOIs) within GRiD by passing one or
more WKT geometries.

This function queries GRiD's Geonames endpoint with the provided geometries and
automatically uses the returned values as the AOI names.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a WKT geometry")
			cmd.Usage()
			return
		}

		for _, geom := range args {
			// get suggested name for the current geometry
			a, _, err := g.Lookup(geom)
			if err != nil {
				log.Fatal(err)
			}

			// create a new AOI for the current geometry with suggested name
			b, _, err := g.AddAOI(a.Name, geom, true)
			if err != nil {
				log.Fatal(err)
			}

			// all this to simply parse the result and print the pk of the new AOI
			success := (*b)["success"].(bool)
			if success {
				delete((*b), "success")

				d, err := json.Marshal(b)
				if err != nil {
					log.Fatal(err)
				}

				c := new(grid.AOIResponse)
				err = json.Unmarshal(d, &c)
				if err != nil {
					log.Fatal(err)
				}

				for _, v := range *c {
					for _, v := range v.AOIs {
						fmt.Printf("Successfully created AOI \"%v\" with primary key \"%v\" at %v\n", v.Fields.Name, v.Pk, v.Fields.CreatedAt)
					}
				}
			}
		}
	},
}
