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
