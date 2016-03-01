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
	"crypto/tls"
	"fmt"
	"net/http"
)

const (
	defaultBaseURL = "https://gridte.rsgis.erdc.dren.mil/te_ba/"
)

var client *http.Client

// GetClient returns our one http.Client, instantiating it if needed
func GetClient() *http.Client {
	if client == nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client = &http.Client{Transport: transport}
	}

	return client
}

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
	request := GetRequestFactory().NewRequest("GET", url)

	err := DoRequest(request, taskObject)
	return taskObject, err
}
