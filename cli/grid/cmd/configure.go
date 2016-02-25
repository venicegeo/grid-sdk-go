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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cobra"
)

func logon() {
	// prompt user for username and password and base64 encode it
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GRiD Username: ")
	username, _ := r.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("GRiD Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	fmt.Print("\nGRiD API Key: ")
	key, _ := r.ReadString('\n')
	key = strings.TrimSpace(key)

	// get the appropriate path for the config.json, depends on platform
	var path string
	if runtime.GOOS == "windows" {
		path = os.Getenv("HOMEPATH")
	} else {
		path = os.Getenv("HOME")
	}
	path = path + string(filepath.Separator) + ".grid"

	// TODO(chambbj): I think this does throw an error on Windows. Need to
	// better understand platform-specific behavior.
	err := os.Mkdir(path, 0777)
	// if err != nil {
	// log.Fatal(err)
	// }

	fileandpath := path + string(filepath.Separator) + "config.json"
	file, err := os.Create(fileandpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	// encode the configuration details as JSON
	config := Config{Auth: auth, Key: key}
	json.NewEncoder(file).Encode(config)
}

// Config represents the config JSON structure.
type Config struct {
	Auth string `json:"auth"`
	Key  string `json:"key"`
}

// getConfig extracts the authoriztion string and API key from the config file.
func getConfig() Config {
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
		logon()
		// fmt.Println("No authentication. Please run 'grid configure' first.")
	}
	var config Config
	b, err := ioutil.ReadAll(file)
	json.Unmarshal(b, &config)
	return config
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the CLI",
	Long: `
Configure the GRiD CLI with the user's GRiD credentials.

This function will prompt the user for their GRiD username and password, which
is encoded in the user's config.json file`,
	Run: func(cmd *cobra.Command, args []string) {
		logon()
	},
}
