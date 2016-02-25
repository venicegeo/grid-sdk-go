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
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

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

var pk int
var gridCmdV *cobra.Command

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
	AddCommands()
	GridCmd.Execute()
}

// AddCommands adds child commands to the root GridCmd.
func AddCommands() {
	GridCmd.AddCommand(addCmd)
	GridCmd.AddCommand(configureCmd)
	GridCmd.AddCommand(lookupCmd)
	GridCmd.AddCommand(lsCmd)
	GridCmd.AddCommand(pullCmd)
	GridCmd.AddCommand(versionCmd)
}

func init() {
	gridCmdV = GridCmd
}

func logon() (username, password, key string) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GRiD Username: ")
	username, _ = r.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("GRiD Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password = string(bytePassword)

	fmt.Print("\nGRiD API Key: ")
	key, _ = r.ReadString('\n')
	key = strings.TrimSpace(key)
	return
}

// Config represents the config JSON structure.
type Config struct {
	Auth string `json:"auth"`
	Key  string `json:"key"`
}

// GetTransport is provided as a convenience to setup the GRiD
// BasicAuthTransport with Basic Authorization string and API key.
func GetTransport() grid.BasicAuthTransport {
	var path string
	if runtime.GOOS == "windows" {
		path = os.Getenv("HOMEPATH")
	} else {
		path = os.Getenv("HOME")
	}
	path = path + string(filepath.Separator) + ".grid"
	fileandpath := path + string(filepath.Separator) + "config.json"
	file, err := os.Open(fileandpath)
	if err != nil {
		log.Fatal("No authentication. Please run 'grid configure' first.")
	}
	var config Config
	b, err := ioutil.ReadAll(file)
	json.Unmarshal(b, &config)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	tp := grid.BasicAuthTransport{
		Auth:      config.Auth,
		Key:       config.Key,
		Transport: tr,
	}

	return tp
}
