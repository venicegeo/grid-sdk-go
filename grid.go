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

// HTTPError represents any HTTP error
type HTTPError struct {
	Status  int
	Message string
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("%d: %v", err.Status, err.Message)
}
