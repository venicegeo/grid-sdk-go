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

// Config represents the config JSON structure.
type Config struct {
	Auth string `json:"auth"`
	Key  string `json:"key"`
	URL  string `json:"url"`
}

/*
logon is called whenever all fields of the config file need to be updated, or
or upon config file creation.
*/
func logon() {
	// prompt user for username and password and base64 encode it
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GRiD Username: ")
	username, _ := r.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("GRiD Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println()

	fmt.Print("GRiD API Key: ")
	key, _ := r.ReadString('\n')
	key = strings.TrimSpace(key)

	fmt.Print("GRiD Base URL: ")
	baseURL, _ := r.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)

	file, err := createConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	// encode the configuration details as JSON
	config := Config{Auth: auth, Key: key, URL: baseURL}
	json.NewEncoder(file).Encode(config)
}

// getConfig extracts config file contents.
func getConfig() Config {
	path := getConfigFilePath()
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

// getConfigFilePath returns the full path to the config file.
func getConfigFilePath() string {
	// get the appropriate path for the config.json, depends on platform
	var path string
	if runtime.GOOS == "windows" {
		path = os.Getenv("HOMEPATH")
	} else {
		path = os.Getenv("HOME")
	}
	path = path + string(filepath.Separator) + ".grid"
	return path
}

// createConfigFile creates the config file for writing, overwriting existing.
func createConfigFile() (*os.File, error) {
	path := getConfigFilePath()

	// TODO(chambbj): I think this does throw an error on Windows. Need to
	// better understand platform-specific behavior.
	err := os.Mkdir(path, 0777)
	// if err != nil {
	// log.Fatal(err)
	// }

	fileandpath := path + string(filepath.Separator) + "config.json"
	file, err := os.Create(fileandpath)
	return file, err
}

// updateBaseURL rewrites the config file, updating only the base URL.
func updateBaseURL(baseURL string) {
	currentConfig := getConfig()
	file, err := createConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// encode the configuration details as JSON
	config := Config{Auth: currentConfig.Auth, Key: currentConfig.Key, URL: baseURL}
	json.NewEncoder(file).Encode(config)
}

// updateBaseURL rewrites the config file, updating only the base URL.
func updateAPIKey(key string) {
	currentConfig := getConfig()
	file, err := createConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// encode the configuration details as JSON
	config := Config{Auth: currentConfig.Auth, Key: key, URL: currentConfig.URL}
	json.NewEncoder(file).Encode(config)
}

var baseURL, key string

func init() {
	configureCmd.Flags().StringVarP(&baseURL, "base_url", "b", "", "GRiD Base URL")
	configureCmd.Flags().StringVarP(&key, "key", "k", "", "GRiD API Key")
}

var configureCmd = &cobra.Command{
	Use:   "configure [-b base_url][-k key]",
	Short: "Configure the CLI",
	Long: `
Configure the GRiD CLI with the user's GRiD credentials.

This function will prompt the user for their GRiD username and password, which
is encoded in the user's config.json file`,
	Run: func(cmd *cobra.Command, args []string) {
		if baseURL != "" {
			updateBaseURL(baseURL)
		} else if key != "" {
			updateAPIKey(key)
		} else {
			logon()
		}
	},
}
