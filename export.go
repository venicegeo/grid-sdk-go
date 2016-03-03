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

	"github.com/venicegeo/pzsvc-sdk-go"
)

// ExportFile represents the export file object that is returned by the export
// endpoint.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#exportfiles-object
type ExportFile struct {
	URL  string `json:"url,omitempty"`
	Pk   int    `json:"pk,omitempty"`
	Name string `json:"name,omitempty"`
}

// TDASet represents the TDA set object that is returned by the export endpoint.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#tda-set-object
type TDASet struct {
	Status    string `json:"status,omitempty"`
	TDAType   string `json:"tda_type,omitempty"`
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Pk        int    `json:"pk,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

// ExportDetail represents the export detail object that is returned by the
// export endpoint.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#export-detail-object
type ExportDetail struct {
	ExportFiles []ExportFile `json:"exportfiles,omitempty"`
	TDASets     []TDASet     `json:"tda_set,omitempty"`
}

// GetExport returns export details for the export specified by the user-provided
// primary key.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#get-export-details
func GetExport(pk int) (*ExportDetail, error) {
	url := fmt.Sprintf("api/v1/export/%v/", pk)

	exportDetail := new(ExportDetail)
	request := sdk.GetRequestFactory().NewRequest("GET", url)
	err := sdk.DoRequest(request, &doRequestCallback{unmarshal: exportDetail})
	return exportDetail, err
}

// DownloadByPk downloads the file specified by the user-provided primary key.
func DownloadByPk(pk int) (*sdk.DownloadCallback, error) {
	url := fmt.Sprintf("export/download/file/%v/", pk)

	request := sdk.GetRequestFactory().NewRequest("GET", url)
	dc := sdk.DownloadCallback{}
	err := sdk.DoRequest(request, &dc)
	return &dc, err
}
