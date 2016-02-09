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

type AOI struct {
	Name       string `json:"name,omitempty"`
	Geometry   string `json:"geometry,omitempty"`
	Notes      string `json:"notes,omitempty"`
	IsActive   bool   `json:"is_active,omitempty"`
	Source     string `json:"source,omitempty"`
	NumExports string `json:"num_exports,omitempty"`
	Pk         int    `json:"pk,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

type AOIList struct {
	AOIs []AOI `json:"self_aoi_list"`
}

type Export struct {
	Status    string `json:"status,omitempty"`
	StartedAt string `json:"stated_at,omitempty"`
	Name      string `json:"name,omitempty"`
	Datatype  string `json:"datatype,omitempty"`
	HSRS      int    `json:"hsrs,omitempty"`
	URL       string `json:"url,omitempty"`
	Pk        int    `json:"pk,omitempty"`
}

type AOIDetail struct {
	ExportSet []Export `json:"export_set"`
}

func (s *AOIService) List(geom string) ([]AOI, *Response, error) {
	url := "api/v0/aoi/"

	req, err := s.client.NewRequest("GET", url, nil)

	authstr := GetAuth()
	if authstr == "" {
		un, pw := Logon()
		req.SetBasicAuth(un, pw)
	} else {
		req.Header.Add("authorization", "Basic "+authstr)
	}
	aoiList := new(AOIList)
	resp, err := s.client.Do(req, aoiList)
	return aoiList.AOIs, resp, err
}

func (s *AOIService) ListByPk(pk int) ([]Export, *Response, error) {
	url := fmt.Sprintf("api/v0/aoi/%v/", pk)

	req, err := s.client.NewRequest("GET", url, nil)

	authstr := GetAuth()
	if authstr == "" {
		un, pw := Logon()
		req.SetBasicAuth(un, pw)
	} else {
		req.Header.Add("authorization", "Basic "+authstr)
	}
	aoiDetail := new(AOIDetail)
	resp, err := s.client.Do(req, aoiDetail)
	return aoiDetail.ExportSet, resp, err
}
