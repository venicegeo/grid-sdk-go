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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/venicegeo/grid-sdk-go"
)

// Config represents the config JSON structure.
type config struct {
	Auth string `json:"auth"`
	Key  string `json:"key"`
}

func getConfig() config {
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
	var c config
	b, err := ioutil.ReadAll(file)
	json.Unmarshal(b, &c)
	return c
}

func getTransport() grid.BasicAuthTransport {
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi")
}

func checkError(err error, status int, writer http.ResponseWriter) bool {
	if err != nil || status < 200 || status > 299 {
		var errorText string
		if err == nil {
			errorText = "S'broke"
		} else {
			errorText = err.Error()
		}
		http.Error(writer, errorText, status)
		return true
	}
	return false
}

func lookupHandler(writer http.ResponseWriter, request *http.Request) {
	geom := request.FormValue("geom")

	switch request.FormValue("version") {
	case "0":
		http.Error(writer, "This method is not supported by this version of the API", http.StatusMethodNotAllowed)
	case "1":
		fallthrough
	default:
		tp := getTransport()
		client := grid.NewClient(tp.Client())

		geomObject, resp, err := client.Geonames.Lookup(geom, tp.Key)
		if err != nil {
			log.Fatal(err.Error())
		}

		if !checkError(err, resp.StatusCode, writer) {
			fmt.Fprintf(writer, geomObject.Name)
		}
	}
}

func getAoiHandler(writer http.ResponseWriter, request *http.Request) {
	pk := request.FormValue("pk")

	switch request.FormValue("version") {
	case "0":
		http.Error(writer, "This method is not supported by this version of the API", http.StatusMethodNotAllowed)
	case "1":
		fallthrough
	default:
		tp := getTransport()
		client := grid.NewClient(tp.Client())

		pki, err := strconv.Atoi(pk)
		aoiObject, resp, err := client.AOI.Get(pki, tp.Key)
		if err != nil {
			log.Fatal(err.Error())
		}

		if !checkError(err, resp.StatusCode, writer) {
			bytes, _ := json.Marshal(aoiObject)
			fmt.Fprintf(writer, string(bytes))
		}
	}
}

func addAoiHandler(writer http.ResponseWriter, request *http.Request) {
	geom := request.FormValue("geom")
	name := request.FormValue("name")
	subscribe, error := strconv.ParseBool(request.FormValue("subscribe"))
	if error != nil {
		subscribe = false
	}

	switch request.FormValue("version") {
	case "0":
		http.Error(writer, "This method is not supported by this version of the API", http.StatusMethodNotAllowed)
	case "1":
		fallthrough
	default:
		tp := getTransport()
		client := grid.NewClient(tp.Client())

		uploadObject, resp, err := client.AOI.Add(name, geom, tp.Key, subscribe)
		if err != nil {
			log.Fatal(err.Error())
		}

		if !checkError(err, resp.StatusCode, writer) {
			bytes, _ := json.Marshal(uploadObject)
			fmt.Fprintf(writer, string(bytes))
		}
	}
}

func main() {

	log.Printf("Hello world.")
	http.HandleFunc("/", handler)
	http.HandleFunc("/lookup", lookupHandler)
	http.HandleFunc("/getaoi", getAoiHandler)
	http.HandleFunc("/addaoi", addAoiHandler)
	http.ListenAndServe(":8080", nil)
}
