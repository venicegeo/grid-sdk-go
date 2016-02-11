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

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an AOI",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a WKT geometry")
			cmd.Usage()
			return
		}

		tp := GetTransport()
		client := grid.NewClient(tp.Client())
		key := GetKey()

		for _, geom := range args {
			a, _, err := client.Geonames.Lookup(geom, key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(a.Name)
			b, _, err := client.AOI.Add(a.Name, geom, key, true)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", *b)
		}
	},
}
