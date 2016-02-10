package grid

import (
	"fmt"
	"net/url"
)

// AOIService handles communication with the AOI related
// methods of the GRiD API.
//
// GRiD API docs:
// https://github.com/CRREL/GRiD-API/blob/v0.0/composed_api.rst#get-a-users-aoi-list
// https://github.com/CRREL/GRiD-API/blob/v0.0/composed_api.rst#get-aoi-details
type GeonamesService struct {
	client *Client
}

type GeonamesResponse struct {
	Name             string `json:"name,omitempty"`
	ProvidedGeometry string `json:"provided_geometry,omitempty"`
}

func (s *GeonamesService) Lookup(geom string) (*GeonamesResponse, *Response, error) {
	v := url.Values{}
	v.Set("geom", geom)
	v.Add("source", "toasted_filament")
	vals := v.Encode()
	qurl := fmt.Sprintf("api/v1/geoname/?%v", vals)

	req, err := s.client.NewRequest("GET", qurl, nil)

	geonamesResponse := new(GeonamesResponse)
	resp, err := s.client.Do(req, geonamesResponse)
	return geonamesResponse, resp, err
}
