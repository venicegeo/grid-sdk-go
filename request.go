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
	"net/url"
)

// RequestDecorator decorates an http.Request
type RequestDecorator interface {
	Decorate(*http.Request) error
}

// RequestFactory creates http.Request decorated as needed
type RequestFactory struct {
	BaseURL    string
	decorators []RequestDecorator
}

var rf *RequestFactory

// GetRequestFactory returns the one RequestFactory
func GetRequestFactory() *RequestFactory {
	if rf == nil {
		rf = new(RequestFactory)
	}
	return rf
}

// AddDecorator adds a RequestDecorator to the RequestFactory
func (rf *RequestFactory) AddDecorator(rd RequestDecorator) {
	rf.decorators = append(rf.decorators, rd)
}

// NewRequest creates a new http.Request, decorating it as needed
func (rf *RequestFactory) NewRequest(method, relativeURL string) *http.Request {
	baseURL, _ := url.Parse(rf.BaseURL)
	parsedRelativeURL, _ := url.Parse(relativeURL)
	resolvedURL := baseURL.ResolveReference(parsedRelativeURL)
	request, _ := http.NewRequest(method, resolvedURL.String(), nil)
	for inx := 0; inx < len(rf.decorators); inx++ {
		rf.decorators[inx].Decorate(request)
	}
	return request
}

// DoRequest performs the request and attempts to unmarshal the response
// into the object provided
func DoRequest(request *http.Request, unmarshal interface{}) error {
	response, err := GetClient().Do(request)
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
	err = json.Unmarshal(body, unmarshal)
	if err != nil {
		log.Printf("Unmarshal error in DoRequest: %v", string(body))
		return &HTTPError{Message: err.Error(), Status: http.StatusNotAcceptable}
	}
	return nil
}
