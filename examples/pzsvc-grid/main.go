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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/venicegeo/grid-sdk-go"
)

var g *grid.Grid

func init() {
	var err error
	g, err = grid.New()
	if err != nil {
		log.Printf("Problem creating GRiD client. Do you have a valid credentials file?")
		log.Fatal(err)
	}
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
