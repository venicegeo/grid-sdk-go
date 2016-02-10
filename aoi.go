package grid

import "fmt"

// AOIService handles communication with the AOI related
// methods of the GRiD API.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/v0.0/composed_api.rst#get-a-users-aoi-list
// https://github.com/CRREL/GRiD-API/blob/v0.0/composed_api.rst#get-aoi-details
type AOIService struct {
	client *Client
}

type Export struct {
	Status    string `json:"status,omitempty"`
	Name      string `json:"name,omitempty"`
	Datatype  string `json:"datatype,omitempty"`
	HSRS      string `json:"hsrs,omitempty"`
	URL       string `json:"url,omitempty"`
	Pk        int    `json:"pk,omitempty"`
	StartedAt string `json:"started_at,omitempty"`
}

type RasterCollect struct {
	Datatype string `json:"datatype,omitempty"`
	Pk       int    `json:"pk,omitempty"`
	Name     string `json:"name,omitempty"`
}

type PointcloudCollect struct {
	Datatype string `json:"datatype,omitempty"`
	Pk       int    `json:"pk,omitempty"`
	Name     string `json:"name,omitempty"`
}

type Fields struct {
	Name         string `json:"name,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	IsActive     bool   `json:"is_active,omitempty"`
	Source       string `json:"source,omitempty"`
	User         int    `json:"user,omitempty"`
	ClipGeometry string `json:"clip_geometry,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type AOI struct {
	Fields Fields `json:"fields,omitempty"`
	Model  string `json:"model,omitempty"`
	Pk     int    `json:"pk,omitempty"`
}

type AOIItem struct {
	ExportSet          []Export            `json:"export_set,omitempty"`
	RasterCollects     []RasterCollect     `json:"raster_collects,omitempty"`
	PointcloudCollects []PointcloudCollect `json:"pointcloud_collects,omitempty"`
	AOIs               []AOI               `json:"aoi,omitempty"`
}

type AOIResponse map[string]AOIItem

func (s *AOIService) List(geom string) (*AOIResponse, *Response, error) {
	url := "api/v1/aoi/?source=toasted_filament"

	req, err := s.client.NewRequest("GET", url, nil)

	aoiList := new(AOIResponse)
	resp, err := s.client.Do(req, aoiList)
	return aoiList, resp, err
}

func (s *AOIService) ListByPk(pk int) (*AOIItem, *Response, error) {
	url := fmt.Sprintf("api/v1/aoi/%v/?source=toasted_filament", pk)

	req, err := s.client.NewRequest("GET", url, nil)

	aoiDetail := new(AOIItem)
	resp, err := s.client.Do(req, aoiDetail)
	return aoiDetail, resp, err
}
