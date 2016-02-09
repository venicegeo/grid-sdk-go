package grid

import "fmt"

// ExportService handles communication with the Export related
// methods of the GRiD API.
//
// GRiD API docs: https://github.com/CRREL/GRiD-API/blob/v0.0/composed_api.rst#get-export-details
type ExportService struct {
	client *Client
}

type File struct {
	URL  string `json:"url"`
	Pk   int    `json:"pk"`
	Name string `json:"name"`
}

type ExportDetail struct {
	ExportFiles []File `json:"exportfiles"`
}

func (s *ExportService) ListByPk(pk int) ([]File, *Response, error) {
	url := fmt.Sprintf("api/v0/export/%v/", pk)

	req, err := s.client.NewRequest("GET", url, nil)

	authstr := GetAuth()
	if authstr == "" {
		un, pw := Logon()
		req.SetBasicAuth(un, pw)
	} else {
		req.Header.Add("authorization", "Basic "+authstr)
	}
	exportDetail := new(ExportDetail)
	resp, err := s.client.Do(req, exportDetail)
	return exportDetail.ExportFiles, resp, err
}
