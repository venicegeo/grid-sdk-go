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

package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/gorilla/mux"
	"github.com/venicegeo/grid-sdk-go"
)

// GetConfig extracts the authoriztion string and API key from the config file.
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

var g *grid.Grid

func init() {
	config := getConfig()
	g = grid.NewClient(config.Auth, config.Key, "")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi")
}

// processError checks for an error condition.
// If one is found, it is applied to the ResponseWriter.
func processError(err error, writer http.ResponseWriter) {
	if err == nil {
		return
	}
	var status int
	var message string
	switch err.(type) {
	// case grid.HTTPError:
	// 	message = err.(grid.HTTPError).Message
	// 	status = err.(grid.HTTPError).Status
	default:
		message = err.Error()
		status = http.StatusBadRequest
	}
	http.Error(writer, message, status)
}

func lookupHandler(writer http.ResponseWriter, request *http.Request) {
	geom := request.FormValue("geom")

	// switch request.FormValue("version") {
	// case "0":
	// 	http.Error(writer, "This method is not supported by this version of the API", http.StatusMethodNotAllowed)
	// case "1":
	// 	fallthrough
	// default:
	geomObject, _, err := g.Lookup(geom)
	if err == nil {
		bytes, err := json.Marshal(geomObject)
		if err == nil {
			fmt.Fprintf(writer, string(bytes))
		}
	}
	processError(err, writer)
}

func getAoiHandler(writer http.ResponseWriter, request *http.Request) {
	pkStr := request.FormValue("pk")

	pk, err := strconv.Atoi(pkStr)
	if err == nil {
		aoiObject, _, err := g.GetAOI(pk)
		if err == nil {
			bytes, err := json.Marshal(aoiObject)
			if err == nil {
				fmt.Fprintf(writer, string(bytes))
			}
		}
	}
	//
	// switch request.FormValue("version") {
	// case "0":
	// 	http.Error(writer, "This method is not supported by this version of the API", http.StatusMethodNotAllowed)
	// case "1":
	// 	fallthrough
	// default:
	// 	aoiObject, eo := v1.GetAOI(pk)
	// 	if !checkError(eo, writer) {
	// 		bytes, _ := json.Marshal(aoiObject)
	// 		fmt.Fprintf(writer, string(bytes))
	// 	}
	// }
	processError(err, writer)
}

func addAoiHandler(writer http.ResponseWriter, request *http.Request) {
	geom := request.FormValue("geom")
	name := request.FormValue("name")
	subscribe, error := strconv.ParseBool(request.FormValue("subscribe"))
	if error != nil {
		subscribe = false
	}
	addAoiObject, _, err := g.AddAOI(name, geom, subscribe)
	if err == nil {
		bytes, err := json.Marshal(addAoiObject)
		if err == nil {
			fmt.Fprintf(writer, string(bytes))
		}
	}
	processError(err, writer)
}

func generatePointCloudHandler(writer http.ResponseWriter, request *http.Request) {
	aoiStr := request.FormValue("aoi")
	aoi, err := strconv.Atoi(aoiStr)
	if err == nil {
		collectsStr := request.FormValue("collects")
		collects := strings.Split(collectsStr, "+")
		options := grid.NewGeneratePointCloudExportOptions()
		generateObject, _, err := g.GeneratePointCloudExport(aoi, collects, options)
		if err == nil {
			bytes, err := json.Marshal(generateObject)
			if err == nil {
				fmt.Fprintf(writer, string(bytes))
			}
		}
	}
	processError(err, writer)
}

func exportHandler(writer http.ResponseWriter, request *http.Request) {
	exportStr := request.FormValue("export_id")

	export, err := strconv.Atoi(exportStr)
	if err == nil {
		exportObject, _, err := g.GetExport(export)
		if err == nil {
			bytes, err := json.Marshal(exportObject)
			if err == nil {
				fmt.Fprintf(writer, string(bytes))
			}
		}
	}
	processError(err, writer)
}

func taskHandler(writer http.ResponseWriter, request *http.Request) {
	task := request.FormValue("task_id")

	taskObject, _, err := g.TaskDetails(task)
	if err == nil {
		bytes, err := json.Marshal(taskObject)
		if err == nil {
			fmt.Fprintf(writer, string(bytes))
		}
	}
	processError(err, writer)
}

const (
	defaultBaseURL = "https://gridte.rsgis.erdc.dren.mil/te_ba/"
)

type server struct {
	router *mux.Router
}

func (server *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	server.router.ServeHTTP(rw, req)
}

func main() {

	var args = os.Args[1:]
	var port string
	if len(args) > 0 {
		port = ":" + args[0]
	} else {
		port = ":8080"
	}

	router := mux.NewRouter()
	http.Handle("/", &server{router})
	log.Printf("Hello port %v", port)
	router.HandleFunc("/", handler)
	router.HandleFunc("/lookup", lookupHandler)
	router.HandleFunc("/getaoi", getAoiHandler)
	router.HandleFunc("/addaoi", addAoiHandler)
	router.HandleFunc("/task", taskHandler)
	router.HandleFunc("/export", exportHandler)
	router.HandleFunc("/generatepointcloud", generatePointCloudHandler)
	http.ListenAndServe(port, nil)
}
