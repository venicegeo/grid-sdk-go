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
	"fmt"
	"net/url"
)

// AOIService handles communication with the AOI related
// methods of the GRiD API.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/v0.0/composed_api.rst#get-a-users-aoi-list
// https://github.com/CRREL/GRiD-API/blob/v0.0/composed_api.rst#get-aoi-details
type GeonamesService struct {
	client *Client
}

type GeonamesResponse struct {
	Name             string `json:"name,omitempty"`
	ProvidedGeometry string `json:"provided_geometry,omitempty"`
}

func (s *GeonamesService) Lookup(geom string) (*GeonamesResponse, *Response, error) {
	v := url.Values{}
	v.Set("geom", geom)
	v.Add("source", "toasted_filament")
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v1/geoname/?%v", vals)

	req, err := s.client.NewRequest("GET", qurl, nil)

	geonamesResponse := new(GeonamesResponse)
	resp, err := s.client.Do(req, geonamesResponse)
	return geonamesResponse, resp, err
}
