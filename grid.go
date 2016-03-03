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

package grid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/venicegeo/pzsvc-sdk-go"
)

// ErrorObject represents any error returnable by the GRiD service v1
type ErrorObject struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// TaskObject represents the state of a GRiD task
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#task-object
type TaskObject struct {
	Traceback string `json:"task_traceback,omitempty"`
	State     string `json:"task_state,omitempty"`
	Timestamp string `json:"task_tstamp,omitempty"`
	Name      string `json:"task_name,omitempty"`
	TaskID    string `json:"task_id,omitempty"`
}

// HTTPError represents any HTTP error
type HTTPError struct {
	Status  int
	Message string
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("%d: %v", err.Status, err.Message)
}

// TaskDetails returns the details for a GRiD task
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#generate-point-cloud-export
func TaskDetails(pk string) (*TaskObject, error) {
	taskObject := new(TaskObject)
	url := fmt.Sprintf("api/v1/task/%v/", pk)
	request := sdk.GetRequestFactory().NewRequest("GET", url)
	err := sdk.DoRequest(request, &doRequestCallback{unmarshal: taskObject})
	return taskObject, err
}

type doRequestCallback struct {
	unmarshal interface{}
}

func (drc doRequestCallback) Callback(response *http.Response, err error) error {
	if err != nil {
		// log.Printf("DoRequest failed to Do the request")
		return &HTTPError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	// Check for HTTP errors
	if response.StatusCode < 200 || response.StatusCode > 299 {
		message := fmt.Sprintf("%v returned %v", response.Request.URL.String(), string(body))
		return &HTTPError{Message: message, Status: response.StatusCode}
	}

	// Unmarshal the result to determine if it is in fact an error
	// masquerading as a StatusOK
	eo := new(ErrorObject)
	err = json.Unmarshal(body, eo)
	if err != nil {
		log.Printf("Unmarshal error in ErrorCheck: %v", string(body))
		return &HTTPError{Message: err.Error(), Status: http.StatusNotAcceptable}
	} else if eo.Error != "" {
		// log.Printf("ErrorCheck discovered %v", eo.Error)
		return &HTTPError{Message: eo.Error, Status: http.StatusBadRequest}
	}

	// log.Printf("Body: %v", string(body))

	// If we've gotten this far, hopefully we can unmarshal properly
	err = json.Unmarshal(body, drc.unmarshal)
	if err != nil {
		log.Printf("Unmarshal error in DoRequest: %v", string(body))
		return &HTTPError{Message: err.Error(), Status: http.StatusNotAcceptable}
	}
	return nil

}

// Config represents the config JSON structure.
type Config struct {
	Auth string `json:"auth"`
	Key  string `json:"key"`
}

var config *Config

// ConfigSourceDecorator adds a source (API Key) to a request
// based on the contents of a Config file
type ConfigSourceDecorator struct {
	Project string
}

// Decorate decorates as per sdk.RequestDecorator
func (csd ConfigSourceDecorator) Decorate(request *http.Request) error {
	var err error
	if config == nil {
		config = new(Config)
		bytes, _ := sdk.GetConfig(csd.Project)
		err = json.Unmarshal(bytes, config)
	}
	if err == nil {
		query := request.URL.Query()
		query.Add("source", config.Key)
		request.URL.RawQuery = query.Encode()
	}
	return err
}
