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

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an AOI",
	Long:  "",
	Run:   RunAdd,
}

func RunAdd(cmd *cobra.Command, args []string) {
	geom := "POLYGON((30 10,40 40,20 40,10 20,30 10))"
	//geom := args[0]
	tp := GetTransport()

	// github client configured to use test server
	client := grid.NewClient(tp.Client())
	a, _, err := client.AOI.Add("foo", geom, true)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(a.Success)
}