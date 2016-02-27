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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GeonamesService handles communication with the Geoname related methods of the
// GRiD API.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#lookup-geoname
type GeonamesService struct {
	client *Client
}

// Geoname represents the geoname object that is returned by the geoname
// endpoint.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#geoname-object
type Geoname struct {
	Name     string `json:"name,omitempty"`
	Geometry string `json:"provided_geometry,omitempty"`
}

// Lookup the suggested name for the given geometry.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#lookup-geoname
func (s *GeonamesService) Lookup(geom string) (*Geoname, *Response, error) {
	if geom == "" {
		return nil, nil, errors.New("Please provide a WKT geometry string")
	}

	v := url.Values{}
	v.Set("geom", geom)
	// v.Add("source", key)
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v1/geoname/?%v", vals)

	req, err := s.client.NewRequest("GET", qurl, nil)

	name := new(Geoname)
	resp, err := s.client.Do(req, name)
	return name, resp, err
}

// Lookup looks up the suggested name for the given geometry.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#lookup-geoname
func Lookup(geom string) (*Geoname, *HttpError) {
	geoname := new(Geoname)
	if geom == "" {
		return geoname, &HttpError{Error: "Please provide a WKT geometry string", Status: http.StatusBadRequest}
	}

	v := url.Values{}
	v.Set("geom", geom)
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v1/geoname/?%v", vals)

	request := GetRequestFactory().NewRequest("GET", qurl)

	response, err := GetClient().Do(request)
	if err != nil {
		return nil, &HttpError{Error: err.Error(), Status: http.StatusInternalServerError}
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	eo := ErrorCheck(&body)
	if eo == nil {
		err = json.Unmarshal(body, geoname)
		if err != nil {
			eo = &HttpError{Error: err.Error(), Status: http.StatusBadRequest}
		}
	}
	return geoname, eo
}
