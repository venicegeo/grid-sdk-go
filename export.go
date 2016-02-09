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

	exportDetail := new(ExportDetail)
	resp, err := s.client.Do(req, exportDetail)
	return exportDetail.ExportFiles, resp, err
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
	os.Rename(file.Name(), params["filename"])
	return resp, err
}
