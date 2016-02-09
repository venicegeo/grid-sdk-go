package grid

import (
	"fmt"
	"io"
	"log"
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

	var foo interface{}
	resp, err := s.client.Do(req, foo)

	cd := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(cd)
	fname := params["filename"]
	file, err := os.Create(fname)
	defer file.Close()

	numBytes, err := io.Copy(file, resp.Body)
	log.Println("Downloaded", numBytes, "bytes to", fname)
	return resp, err
}
