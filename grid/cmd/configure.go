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
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the CLI",
	Long:  "",
	Run:   Configure,
}

// Configure configures the CLI.
func Configure(cmd *cobra.Command, args []string) {
	un, pw := logon()
	data := []byte(un + ":" + pw)
	str := base64.StdEncoding.EncodeToString(data)
	var path string
	if runtime.GOOS == "windows" {
		path = os.Getenv("HOMEPATH")
	} else {
		path = os.Getenv("HOME")
	}
	path = path + string(filepath.Separator) + ".grid"
	err := os.Mkdir(path, 0777)
	// if err != nil {
	// log.Fatal(err)
	// }
	fileandpath := path + string(filepath.Separator) + "credentials"
	file, err := os.Create(fileandpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintln(file, str)
}
