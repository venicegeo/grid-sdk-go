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
	"net/http"
	"net/url"
)

// RequestDecorator decorates an http.Request
type RequestDecorator interface {
	Decorate(*http.Request)
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
