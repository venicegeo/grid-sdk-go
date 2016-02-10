package cmd

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
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
	Use:   "grid",
	Short: "grid short",
	Long: `grid is the main command.

grid provide CLI access to GRiD.`,
}

var pk int
var gridCmdV *cobra.Command

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Grid",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Grid v0.1 -- HEAD")
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
	GridCmd.AddCommand(aoiCmd)
	GridCmd.AddCommand(exportCmd)
	GridCmd.AddCommand(fileCmd)
	GridCmd.AddCommand(pullCmd)
	GridCmd.AddCommand(configureCmd)
	GridCmd.AddCommand(versionCmd)
}

func init() {
	gridCmdV = GridCmd
}

func logon() (username, password string) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GRiD Username: ")
	username, _ = r.ReadString('\n')

	fmt.Print("GRiD Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password = string(bytePassword)
	return
}

func getAuth() string {
	var path string
	if runtime.GOOS == "windows" {
		path = os.Getenv("HOMEPATH")
	} else {
		path = os.Getenv("HOME")
	}
	path = path + string(filepath.Separator) + ".grid"
	fileandpath := path + string(filepath.Separator) + "credentials"
	file, err := os.Open(fileandpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	return line
}

func GetTransport() grid.BasicAuthTransport {
	authstr := getAuth()
	var username, password string
	if authstr == "" {
		username, password = logon()
	} else {
		data, _ := base64.StdEncoding.DecodeString(authstr)
		up := strings.Split(string(data), ":")
		username = up[0]
		password = up[1]
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	tp := grid.BasicAuthTransport{
		Username:  strings.TrimSpace(username),
		Password:  strings.TrimSpace(password),
		Transport: tr,
	}

	return tp
}
