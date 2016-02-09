package grid

// AOIService handles communication with the AOI related
// methods of the GRiD API.
//
// GRiD API docs: https://github.com/CRREL/GRiD-API/blob/v0.0/composed_api.rst#get-a-users-aoi-list
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
	// aois, err := ioutil.ReadAll(resp.Body)
	// resp.Body.Close()
	// Check(err)
	// a := &AOIList{}
	// err = json.Unmarshal([]byte(aois), &a)
	// Check(err)
	return aoiList.AOIs, resp, err
}
