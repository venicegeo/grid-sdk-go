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
	"net/http"
	"net/url"

	"github.com/venicegeo/pzsvc-sdk-go"
)

// Geoname represents the geoname object that is returned by the geoname
// endpoint.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#geoname-object
type Geoname struct {
	Name     string `json:"name,omitempty"`
	Geometry string `json:"provided_geometry,omitempty"`
}

// Lookup looks up the suggested name for the given geometry.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#lookup-geoname
func Lookup(geom string) (*Geoname, error) {
	geoname := new(Geoname)
	if geom == "" {
		return geoname, &HTTPError{Message: "Please provide a WKT geometry string", Status: http.StatusBadRequest}
	}

	v := url.Values{}
	v.Set("geom", geom)
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v1/geoname/?%v", vals)

	request := sdk.GetRequestFactory().NewRequest("GET", qurl)

	drc := doRequestCallback{unmarshal: geoname}
	err := sdk.DoRequest(request, drc)
	return geoname, err
}
