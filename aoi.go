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
)

// Export represents the export object that is returned as part of an AOIItem.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#export-object
type Export struct {
	Status    string `json:"status,omitempty"`
	Name      string `json:"name,omitempty"`
	Datatype  string `json:"datatype,omitempty"`
	HSRS      string `json:"hsrs,omitempty"`
	URL       string `json:"url,omitempty"`
	Pk        int    `json:"pk,omitempty"`
	StartedAt string `json:"started_at,omitempty"`
}

// RasterCollect represents the raster collect object that is returned as part
// of an AOIItem.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#collect-object
type RasterCollect struct {
	Datatype string `json:"datatype,omitempty"`
	Pk       int    `json:"pk,omitempty"`
	Name     string `json:"name,omitempty"`
}

// PointcloudCollect represents the pointcloud collect object that is returned
// as part of an AOIItem.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#collect-object
type PointcloudCollect struct {
	Datatype string `json:"datatype,omitempty"`
	Pk       int    `json:"pk,omitempty"`
	Name     string `json:"name,omitempty"`
}

// AOI represents the AOI object that is returned by the AOI detail endpoint.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#aoi-object2
type AOI struct {
	// Fields Fields `json:"fields,omitempty"`
	Fields struct {
		Name         string `json:"name,omitempty"`
		CreatedAt    string `json:"created_at,omitempty"`
		IsActive     bool   `json:"is_active,omitempty"`
		Source       string `json:"source,omitempty"`
		User         int    `json:"user,omitempty"`
		ClipGeometry string `json:"clip_geometry,omitempty"`
		Notes        string `json:"notes,omitempty"`
	} `json:"fields,omitempty"`
	Model string `json:"model,omitempty"`
	Pk    int    `json:"pk,omitempty"`
}

// AOIItem represents the AOI object that is returned by the AOI list endpoint.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#aoi-detail-object
type AOIItem struct {
	ExportSet          []Export            `json:"export_set,omitempty"`
	RasterCollects     []RasterCollect     `json:"raster_collects,omitempty"`
	PointcloudCollects []PointcloudCollect `json:"pointcloud_collects,omitempty"`
	AOIs               []AOI               `json:"aoi,omitempty"`
}

// AOIResponse represents the collection of AOIItems returned by the AOI list
// endpoint.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#aoi-object
type AOIResponse map[string]AOIItem

// AddAOIResponse represents the response returned by the AOI add endpoint.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#aoi-detail-object
type AddAOIResponse map[string]interface{}

// ListAOIs retrieves all AOIs intersecting the optional geometry.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#get-a-users-aoi-list
func ListAOIs(geom string) (*AOIResponse, error) {
	url := "api/v1/aoi"

	if geom != "" {
		return nil, &HTTPError{Status: http.StatusNotImplemented, Message: "This method does not currently accept geometries."}
	}
	aoiList := new(AOIResponse)
	request := GetRequestFactory().NewRequest("GET", url)

	err := DoRequest(request, aoiList)
	return aoiList, err
}

// GetAOI returns AOI details for the AOI specified by the user-provided primary
// key.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#get-aoi-details
func GetAOI(pk int) (*AOIItem, error) {
	url := fmt.Sprintf("api/v1/aoi/%v/", pk)

	aoiDetail := new(AOIItem)
	request := GetRequestFactory().NewRequest("GET", url)

	err := DoRequest(request, aoiDetail)
	return aoiDetail, err
}

// AddAOI uploads the given geometry to create a new AOI.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#add-aoi
func AddAOI(name, geom string, subscribe bool) (*AddAOIResponse, error) {
	addAOIResponse := new(AddAOIResponse)
	if name == "" {
		return addAOIResponse, &HTTPError{Message: "Please provide an AOI name and WKT geometry string", Status: http.StatusBadRequest}
	}

	if geom == "" {
		return addAOIResponse, &HTTPError{Message: "Please provide a WKT geometry string", Status: http.StatusBadRequest}
	}

	v := url.Values{}
	v.Set("geom", geom)
	v.Add("name", name)
	if subscribe {
		v.Add("subscribe", "True")
	}
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v1/aoi/add/?%v", vals)

	request := GetRequestFactory().NewRequest("GET", qurl)

	err := DoRequest(request, addAOIResponse)

	return addAOIResponse, err
}
