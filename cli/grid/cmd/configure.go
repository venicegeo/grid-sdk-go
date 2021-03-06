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
	"log"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/venicegeo/grid-sdk-go"
)

func readLine(prompt string) (input string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(input)), nil
}

func readPassword(prompt string) (passwd string, err error) {
	fmt.Print(prompt)

	password, err := gopass.GetPasswd()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(password)), nil
}

/*
logon is called whenever all fields of the config file need to be updated, or
or upon config file creation.
*/
func logon() error {
	username, err := readLine("GRiD Username: ")
	if err != nil {
		return err
	}

	password, err := readPassword("GRiD Password: ")
	if err != nil {
		return err
	}

	baseURL, err := readLine("GRiD Base URL: ")
	if err != nil {
		return err
	}

	file, err := grid.CreateConfigFile()
	if err != nil {
		return err
	}
	defer file.Close()

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	// encode the configuration details as JSON
	config := grid.Config{Auth: auth, URL: baseURL}
	json.NewEncoder(file).Encode(config)

	return nil
}

// updateBaseURL rewrites the config file, updating only the base URL.
func updateBaseURL(baseURL string) {
	cfg, err := grid.GetConfig()
	if err != nil {
		err := logon()
		if err != nil {
			panic(err)
		}
	} else {
		file, err := grid.CreateConfigFile()
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// encode the configuration details as JSON
		config := grid.Config{Auth: cfg.Auth, URL: baseURL}
		json.NewEncoder(file).Encode(config)
	}
}

var baseURL string

func init() {
	configureCmd.Flags().StringVarP(&baseURL, "base_url", "b", "", "GRiD Base URL")
}

var configureCmd = &cobra.Command{
	Use:   "configure [-b base_url]",
	Short: "Configure the CLI",
	Long: `
Configure the GRiD CLI with the user's GRiD credentials.

This function will prompt the user for their GRiD username and password, which
is encoded in the user's config.json file`,
	Run: func(cmd *cobra.Command, args []string) {
		if baseURL != "" {
			updateBaseURL(baseURL)
		} else {
			err := logon()
			if err != nil {
				panic(err)
			}
		}
	},
}
