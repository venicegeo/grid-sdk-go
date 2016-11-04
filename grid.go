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

/*
Package grid provides a client for using the GRiD API.

Portions of the grid package borrow heavily from
https://github.com/google/go-github, a Go library for accessing the GitHub API,
which is released under a BSD-style license
(https://github.com/google/go-github/blob/master/LICENSE), with additional
inspiration drawn from https://github.com/Medium/medium-sdk-go, a similar
library for accessing the Medium API, and released under Apache v2.0.
*/
package grid

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
)

const (
	defaultBaseURL = "https://rsgis.erdc.dren.mil/te_ba/"
	apiKey         = "CM69OHTGZJ2F08ET"
)

// All the types

// AOIArray represents the AOI object that is returned by the AOI list endpoint.
type AOIArray struct {
	AOIList []struct {
		Name      string `json:"name,omitempty"`
		CreatedAt string `json:"created_at,omitempty"`
		IsActive  bool   `json:"is_active,omitempty"`
		Source    string `json:"source,omitempty"`
		User      int    `json:"user,omitempty"`
		Geometry  string `json:"geometry,omitempty"`
		Notes     string `json:"notes,omitempty"`
		Pk        int    `json:"pk,omitempty"`
	} `json:"aoi_list,omitempty"`
}

// AOIDetail represents the AOI object that is returned by the AOI detail
// endpoint.
type AOIDetail struct {
	Name                 string                    `json:"name,omitempty"`
	CreatedAt            string                    `json:"created_at,omitempty"`
	IsActive             bool                      `json:"is_active,omitempty"`
	Source               string                    `json:"source,omitempty"`
	User                 int                       `json:"user,omitempty"`
	Geometry             string                    `json:"geometry,omitempty"`
	Notes                string                    `json:"notes,omitempty"`
	Pk                   int                       `json:"pk,omitempty"`
	ExportSet            []Export                  `json:"export_set,omitempty"`
	PointcloudIntersects []PointcloudDatasetSimple `json:"pointcloud_intersects,omitempty"`
	RasterIntersects     []RasterDatasetSimple     `json:"raster_intersects,omitempty"`
}

// Config represents the config JSON structure.
type Config struct {
	Auth string `json:"auth"`
	URL  string `json:"url"`
}

/*
An Error reports more details on an individual error in an ErrorResponse.
These are the possible validation error codes:
    missing:
        resource does not exist
    missing_field:
        a required field on a resource has not been set
    invalid:
        the formatting of a field is invalid
    already_exists:
        another resource has the same valid as this field
GitHub API docs: http://developer.github.com/v3/#client-errors
*/
type Error struct {
	Resource string `json:"resource"` // resource on which the error occurred
	Field    string `json:"field"`    // field on which the error occurred
	Code     string `json:"code"`     // validation error code
}

/*
An ErrorResponse reports one or more errors caused by an API request.
GitHub API docs: http://developer.github.com/v3/#client-errors
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
	Errors   []Error        `json:"errors"`  // more detail on individual errors
}

// Export represents the export object that is returned as part of an AOIDetail.
type Export struct {
	Status    string `json:"status,omitempty"`
	Name      string `json:"name,omitempty"`
	Datatype  string `json:"datatype,omitempty"`
	HSRS      string `json:"hsrs,omitempty"`
	URL       string `json:"url,omitempty"`
	Pk        int    `json:"pk,omitempty"`
	StartedAt string `json:"started_at,omitempty"`
	User      int    `json:"user,omitempty"`
}

// ExportDetail represents the export object that is returned as part of an
// AOIDetail.
type ExportDetail struct {
	Status            string       `json:"status,omitempty"`
	Name              string       `json:"name,omitempty"`
	Datatype          string       `json:"datatype,omitempty"`
	HSRS              string       `json:"hsrs,omitempty"`
	URL               string       `json:"url,omitempty"`
	Pk                int          `json:"pk,omitempty"`
	StartedAt         string       `json:"started_at,omitempty"`
	User              int          `json:"user,omitempty"`
	RGB               bool         `json:"rgb,omitempty"`
	Intensity         bool         `json:"intensity,omitempty"`
	DimClassification bool         `json:"dim_classification,omitempty"`
	FileExportOptions string       `json:"file_export_options,omitempty"`
	GenerateDEM       bool         `json:"generate_dem,omitempty"`
	CellSpacing       float32      `json:"cell_spacing,omitempty"`
	Notes             string       `json:"notes,omitempty"`
	Classification    string       `json:"classification,omitempty"`
	PCLTerrain        string       `json:"pcl_terrain,omitempty"`
	SRIHRES           float32      `json:"sri_hres,omitempty"`
	ExportFiles       []Exportfile `json:"exportfiles,omitempty"`
	TDASet            []TDA        `json:"tda_set,omitempty"`
	TaskID            string       `json:"task_id,omitempty"`
}

// Exportfile represents the export file object that is returned by the export
// endpoint.
type Exportfile struct {
	Datatype string `json:"datatype,omitempty"`
	Name     string `json:"name,omitempty"`
	Pk       int    `json:"pk,omitempty"`
	URL      string `json:"url,omitempty"`
}

// GenerateExportObject represents the output from a Generate Export operation
type GenerateExportObject struct {
	TaskID   string `json:"task_id,omitempty"`
	ExportID int    `json:"export_id,omitempty"`
}

// GeneratePointCloudExportOptions represents the options for a Generate Point
// Cloud Export Operation
type GeneratePointCloudExportOptions struct {
	Intensity         bool
	DimClassification bool
	Hsrs              string //EPSG code
	FileExportOptions string //individual or collect
	FileExportFormat  string
	Compressed        bool
	SendEmail         bool
	GenerateDem       bool
	CellSpacing       float32
	PclTerrain        string  // urban, mountainous, suburban, or foliated
	SriHResolution    float32 // Horizontal resolution
	DecimationRadius  float32
	RetileSize        float32
	RetileArea        float32
}

// Geoname represents the geoname object that is returned by the geoname
// endpoint.
type Geoname struct {
	Name string `json:"name,omitempty"`
	Geom string `json:"geom,omitempty"`
}

// Grid defines the GRiD client.
type Grid struct {
	Auth string
	// Base URL for API requests.  Defaults to GRiD TE, but can be
	// set to a domain endpoint to use with other instances.  BaseURL should
	// always be specified with a trailing slash.
	BaseURL   *url.URL
	Transport http.RoundTripper
}

// PointcloudCollect represents the pointcloud collect object that is returned
// as part of an AOIDetail.
type PointcloudCollect struct {
	Datatype string `json:"datatype,omitempty"`
	Pk       int    `json:"pk,omitempty"`
	Name     string `json:"name,omitempty"`
}

// PointcloudDatasetSimple ...
type PointcloudDatasetSimple struct {
	Datatype        string  `json:"datatype,omitempty"`
	Name            string  `json:"name,omitempty"`
	Pk              int     `json:"pk,omitempty"`
	Sensor          string  `json:"sensor,omitempty"`
	CollectedAt     string  `json:"collected_at,omitempty"`
	Classification  string  `json:"classification,omitempty"`
	Area            float32 `json:"area,omitempty"`
	Filesize        int     `json:"filesize,omitempty"`
	PointClount     int     `json:"point_count,omitempty"`
	Density         float32 `json:"density,omitempty"`
	PercentCoverage float32 `json:"percent_coverage,omitempty"`
}

// RasterCollect represents the raster collect object that is returned as part
// of an AOIDetail.
type RasterCollect struct {
	Datatype string `json:"datatype,omitempty"`
	Pk       int    `json:"pk,omitempty"`
	Name     string `json:"name,omitempty"`
}

// RasterDatasetSimple ...
type RasterDatasetSimple struct {
	Datatype        string  `json:"datatype,omitempty"`
	Name            string  `json:"name,omitempty"`
	Pk              int     `json:"pk,omitempty"`
	Sensor          string  `json:"sensor,omitempty"`
	CollectedAt     string  `json:"collected_at,omitempty"`
	Classification  string  `json:"classification,omitempty"`
	Area            float32 `json:"area,omitempty"`
	Filesize        int     `json:"filesize,omitempty"`
	PercentCoverage float32 `json:"percent_coverage,omitempty"`
}

/*
Response is a GitHub API response.  This wraps the standard http.Response
returned from GitHub and provides convenient access to things like pagination
links.
*/
type Response struct {
	*http.Response
}

// TaskObject represents the state of a GRiD task
type TaskObject struct {
	Traceback string `json:"task_traceback,omitempty"`
	State     string `json:"task_state,omitempty"`
	Timestamp string `json:"task_tstamp,omitempty"`
	Name      string `json:"task_name,omitempty"`
	TaskID    string `json:"task_id,omitempty"`
}

// TDA ...
type TDA struct {
	CreatedAt string `json:"created_at,omitempty"`
	Name      string `json:"name,omitempty"`
	Notes     string `json:"notes,omitempty"`
	Pk        int    `json:"pk,omitempty"`
	Status    string `json:"status,omitempty"`
	TDAType   string `json:"tda_type,omitempty"`
	URL       string `json:"url,omitempty"`
}

// TDASet represents the TDA set object that is returned by the export endpoint.
type TDASet struct {
	Status    string `json:"status,omitempty"`
	TDAType   string `json:"tda_type,omitempty"`
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Pk        int    `json:"pk,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

/*
CheckResponse checks the API response for errors, and returns them if present.
A response is considered an error if it has a status code outside the 200 range.
API error responses are expected to have either no response body, or a JSON
response body that maps to ErrorResponse. Any other response body will be
silently ignored.
*/
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message, r.Errors)
}

/*
sanitizeURL was originally used (in the GitHub code) to redact the client_id
and client_secret tokens from the URL which may be exposed to the user,
specifically in the ErrorResponse error message.

We may not need it, but then again, maybe we will have our own sanitization.
*/
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	return uri
}

// New returns a new GRiD API client.
func New() (*Grid, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	if config.URL == "" {
		config.URL = defaultBaseURL
	}
	parsedBaseURL, _ := url.Parse(config.URL)
	return &Grid{
		Auth:    config.Auth,
		BaseURL: parsedBaseURL,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}, nil
}

/*
Lookup the suggested name for the given geometry.

GRiD API docs:
https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#lookup-geoname
*/
func (g *Grid) Lookup(geom string) (*Geoname, *Response, error) {
	if geom == "" {
		return nil, nil, errors.New("Please provide a WKT geometry string")
	}

	v := url.Values{}
	v.Set("geom", geom)
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v2/geoname?%v", vals)

	req, err := g.NewRequest("GET", qurl, nil)

	name := new(Geoname)
	resp, err := g.Do(req, name)
	return name, resp, err
}

/*
NewRequest creates an API request. A relative URL can be provided in urlStr, in
which case it is resolved relative to the BaseURL of the Client. Relative URLs
should always be specified without a preceding slash. If  specified, the value
pointed to by body is JSON encoded and included as the request body.
*/
func (g *Grid) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := g.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+g.Auth)

	a := req.URL.Query()
	a.Add("source", apiKey)
	req.URL.RawQuery = a.Encode()

	return req, nil
}

/*
Do sends an API request and returns the API response.  The API response is JSON
decoded and stored in the value pointed to by v, or returned as an error if an
API error has occurred.  If v implements the io.Writer interface, the raw
response body will be written to v, without attempting to first decode it.
*/
func (g *Grid) Do(req *http.Request, v interface{}) (*Response, error) {
	client := &http.Client{
		Transport: g.Transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	return response, err
}

/*
ListAOIs retrieves all AOIs intersecting the optional geometry.

GRiD API docs:
https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#get-a-users-aoi-list
*/
func (g *Grid) ListAOIs(geom string) (*AOIArray, *Response, error) {
	v := url.Values{}
	if geom != "" {
		v.Set("geom", geom)
	}
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v2/aoi?%v", vals)

	req, err := g.NewRequest("GET", qurl, nil)

	aoiList := new(AOIArray)
	resp, err := g.Do(req, aoiList)

	return aoiList, resp, err
}

/*
GetAOI returns AOI details for the AOI specified by the user-provided primary
key.

GRiD API docs:
https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#get-aoi-details
*/
func (g *Grid) GetAOI(pk int) (*AOIDetail, *Response, error) {
	url := fmt.Sprintf("api/v2/aoi/%v", pk)

	req, err := g.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	aoiDetail := new(AOIDetail)
	resp, _ := g.Do(req, aoiDetail)

	return aoiDetail, resp, nil
}

/*
AddAOI uploads the given geometry to create a new AOI.

GRiD API docs:
https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#add-aoi
*/
func (g *Grid) AddAOI(name, geom string, subscribe bool) (*AOIDetail, *Response, error) {
	if name == "" {
		return nil, nil, errors.New("Please provide an AOI name and WKT geometry string")
	}

	if geom == "" {
		return nil, nil, errors.New("Please provide a WKT geometry string")
	}

	v := url.Values{}
	v.Set("geom", geom)
	v.Add("name", name)
	if subscribe {
		v.Add("subscribe", "True")
	}
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v2/aoi/add?%v", vals)

	req, err := g.NewRequest("GET", qurl, nil)
	addAOIResponse := new(AOIDetail)
	resp, err := g.Do(req, addAOIResponse)
	return addAOIResponse, resp, err
}

/*
GetExport returns export details for the export specified by the user-provided
primary key.

GRiD API docs:
https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#get-export-details
*/
func (g *Grid) GetExport(pk int) (*ExportDetail, *Response, error) {
	qurl := fmt.Sprintf("api/v2/export/%v", pk)

	req, err := g.NewRequest("GET", qurl, nil)

	exportDetail := new(ExportDetail)
	resp, err := g.Do(req, exportDetail)
	return exportDetail, resp, err
}

// DownloadByPk downloads the file specified by the user-provided primary key.
func (g *Grid) DownloadByPk(pk int) (*Response, error) {
	url := fmt.Sprintf("export/download/file/%v/", pk)

	req, err := g.NewRequest("GET", url, nil)

	file, err := os.Create("temp")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	resp, err := g.Do(req, file)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cd := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return nil, err
	}
	err = os.Rename(file.Name(), params["filename"])
	return resp, err
}

/*
NewGeneratePointCloudExportOptions is a factory method for a
GeneratePointCloudExportOptions that provides all defaults
*/
func NewGeneratePointCloudExportOptions() *GeneratePointCloudExportOptions {
	return &GeneratePointCloudExportOptions{
		Intensity:         true,
		DimClassification: true,
		FileExportOptions: "individual",
		FileExportFormat:  "las12",
		Compressed:        true,
		SendEmail:         false,
		GenerateDem:       false,
		CellSpacing:       1.0,
	}
}

/*
GeneratePointCloudExport does just that for the given PK and set of products

GRiD API docs:
https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#generate-point-cloud-export
*/
func (g *Grid) GeneratePointCloudExport(pk int, products []string, options *GeneratePointCloudExportOptions) (*GenerateExportObject, *Response, error) {
	if options == nil {
		options = NewGeneratePointCloudExportOptions()
	}
	v := url.Values{}
	for inx := 0; inx < len(products); inx++ {
		v.Add("products", products[inx])
	}
	if !options.Compressed {
		v.Set("compressed", "False")
	}
	if !options.DimClassification {
		v.Set("dim_classification", "False")
	}
	if options.FileExportOptions != "" {
		v.Set("file_export_options", options.FileExportOptions)
	}
	if options.FileExportFormat != "" {
		v.Set("file_export_format", options.FileExportFormat)
	}
	if options.GenerateDem {
		v.Set("generate_dem", "True")
		if options.CellSpacing != 1.0 {
			cellSpacing := fmt.Sprintf("%f", options.CellSpacing)
			v.Set("cell_spacing", cellSpacing)
		}
	}
	if options.Hsrs != "" {
		v.Set("hsrs", options.Hsrs)
	}
	if !options.Intensity {
		v.Set("intensity", "False")
	}
	if options.PclTerrain != "" {
		v.Set("pcl_terrain", options.PclTerrain)
	}
	if options.SendEmail {
		v.Set("send_email", "True")
	}
	if options.SriHResolution != 0 {
		srihres := fmt.Sprintf("%f", options.SriHResolution)
		v.Set("sri_hres", srihres)
	}
	if options.DecimationRadius != 0 {
		radius := fmt.Sprintf("%f", options.DecimationRadius)
		v.Set("decimation_radius", radius)
	}
	if options.RetileSize != 0 {
		rsize := fmt.Sprintf("%f", options.RetileSize)
		v.Set("retile_size", rsize)
	}
	if options.RetileArea != 0 {
		rarea := fmt.Sprintf("%f", options.RetileArea)
		v.Set("retile_area", rarea)
	}
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v2/aoi/%v/generate/pointcloud?%v", pk, vals)

	req, err := g.NewRequest("GET", qurl, nil)
	// fmt.Printf("%+v\n", req)
	geo := new(GenerateExportObject)
	resp, err := g.Do(req, geo)
	return geo, resp, err
}

/*
TaskDetails returns the details for a GRiD task

GRiD API docs:
https://github.com/CRREL/GRiD-API/blob/master/composed_api.rst#generate-point-cloud-export
*/
func (g *Grid) TaskDetails(pk string) (*TaskObject, *Response, error) {
	taskObject := new(TaskObject)
	url := fmt.Sprintf("api/v2/task/%v/", pk)
	req, err := g.NewRequest("GET", url, nil)
	resp, err := g.Do(req, taskObject)
	return taskObject, resp, err
}

// GetConfig extracts config file contents.
func GetConfig() (Config, error) {
	var config Config
	path := getConfigFilePath()
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	b, _ := ioutil.ReadAll(file)
	json.Unmarshal(b, &config)
	return config, nil
}

// getConfigFilePath returns the full path to the config file.
// https://github.com/starkandwayne/cf-cli/blob/master/cf/configuration/config_helpers.go#L9-L20
func getConfigFilePath() string {
	configDir := filepath.Join(userHomeDir(), ".grid")

	err := os.MkdirAll(configDir, 0777)
	if err != nil {
		panic(err)
	}

	return filepath.Join(configDir, "config.json")
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

// CreateConfigFile creates the config file for writing, overwriting existing.
func CreateConfigFile() (*os.File, error) {
	path := getConfigFilePath()
	file, err := os.Create(path)
	return file, err
}
