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
	"bytes"
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

var config *Config

type basicAuthDecorator struct {
}

func (bad basicAuthDecorator) Decorate(request *http.Request) {
	config := getConfig()
	request.Header.Set("Authorization", "Basic "+config.Auth)
}

type configSourceDecorator struct {
}

func (csd configSourceDecorator) Decorate(request *http.Request) {
	config := getConfig()
	query := request.URL.Query()
	query.Add("source", config.Key)
	request.URL.RawQuery = query.Encode()
}

type logDecorator struct {
}

func (ld logDecorator) Decorate(request *http.Request) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	request.Write(writer)
	writer.Flush()
	log.Printf(buffer.String())
}

// Execute adds all child commands to the root command GridCmd and sets flags
// appropriately.
func Execute() {
	addCommands()
	setup()
	GridCmd.Execute()
}

func addCommands() {
	GridCmd.AddCommand(addCmd)
	GridCmd.AddCommand(configureCmd)
	GridCmd.AddCommand(lookupCmd)
	GridCmd.AddCommand(lsCmd)
	GridCmd.AddCommand(pullCmd)
	GridCmd.AddCommand(versionCmd)
}

const (
	defaultBaseURL = "https://gridte.rsgis.erdc.dren.mil/te_ba/"
)

func setup() {
	rf := grid.GetRequestFactory()
	rf.BaseURL = defaultBaseURL
	rf.AddDecorator(new(basicAuthDecorator))
	rf.AddDecorator(new(configSourceDecorator))
	rf.AddDecorator(new(logDecorator))
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

func getConfig() *Config {
	if config != nil {
		return config
	}

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
	b, err := ioutil.ReadAll(file)
	json.Unmarshal(b, &config)
	return config
}

// GetTransport is provided as a convenience to setup the GRiD
// BasicAuthTransport with Basic Authorization string and API key.
func GetTransport() grid.BasicAuthTransport {

	config := getConfig()

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
