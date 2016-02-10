package grid

import (
	"fmt"
	"mime"
	"os"
)

// ExportService handles communication with the Export related
// methods of the GRiD API.
//
// GRiD API docs: https://github.com/CRREL/GRiD-API/blob/v0.0/composed_api.rst#get-export-details
type ExportService struct {
	client *Client
}

type ExportFile struct {
	URL  string `json:"url,omitempty"`
	Pk   int    `json:"pk,omitempty"`
	Name string `json:"name,omitempty"`
}

type TDA struct {
	Status    string `json:"status,omitempty"`
	TDAType   string `json:"tda_type,omitempty"`
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Pk        int    `json:"pk,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

type ExportResponse struct {
	ExportFiles []ExportFile `json:"exportfiles,omitempty"`
	TDASet      []TDA        `json:"tda_set,omitempty"`
}

func (s *ExportService) ListByPk(pk int) (*ExportResponse, *Response, error) {
	url := fmt.Sprintf("api/v1/export/%v/?source=toasted_filament", pk)

	req, err := s.client.NewRequest("GET", url, nil)

	exportDetail := new(ExportResponse)
	resp, err := s.client.Do(req, exportDetail)
	return exportDetail, resp, err
}

func (s *ExportService) DownloadByPk(pk int) (*Response, error) {
	url := fmt.Sprintf("export/download/file/%v/", pk)

	req, err := s.client.NewRequest("GET", url, nil)

	file, err := os.Create("temp")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	resp, err := s.client.Do(req, file)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cd := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		panic(err)
	}
	err = os.Rename(file.Name(), params["filename"])
	return resp, err
}
