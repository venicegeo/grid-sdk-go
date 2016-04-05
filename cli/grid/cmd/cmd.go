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

	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

// GridCmd is Grid's root command. Every other command attached to GridCmd is a
// child command to it.
var GridCmd = &cobra.Command{
	Use: "grid",
	Long: `
grid is a command-line interface to the GRiD database.`,
}

var g *grid.Grid

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the GRiD CLI",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grid v0.1 -- HEAD")
	},
}

// Execute adds all child commands to the root command GridCmd and sets flags
// appropriately.
func Execute() {
	GridCmd.AddCommand(addCmd)
	GridCmd.AddCommand(configureCmd)
	GridCmd.AddCommand(exportCmd)
	GridCmd.AddCommand(lookupCmd)
	GridCmd.AddCommand(lsCmd)
	GridCmd.AddCommand(pullCmd)
	GridCmd.AddCommand(taskCmd)
	GridCmd.AddCommand(versionCmd)

	if err := GridCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// Try once. We'll assume that on first run, we will get an error, which
	// requires that we login.
	var err error
	g, err = grid.New()
	if err != nil {
		fmt.Println("It looks like this is your first time running the GRiD CLI. Please follow the prompts to configure the CLI with your credentials.")
		err := logon()
		if err != nil {
			panic("Error creating credentials file!")
		}

		// After that, things should work fine. If we've got a problem now, then
		// panic!
		g, err = grid.New()
		if err != nil {
			panic(err)
		}
	}
}
