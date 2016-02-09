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
	// aois, err := ioutil.ReadAll(resp.Body)
	// resp.Body.Close()
	// Check(err)
	// a := &AOIList{}
	// err = json.Unmarshal([]byte(aois), &a)
	// Check(err)
	return exportDetail.ExportFiles, resp, err
}

// func (s *ExportService) DownloadByPK(pk int) (*Response, error) {
// 	// url := "api/v0/export/"
//   //
// 	// tr := &http.Transport{
// 	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	// }
// 	// client := &http.Client{Transport: tr}
//   //
// 	// req, err := http.NewRequest("GET", url+strconv.Itoa(pk)+"/", nil)
//
//   url := fmt.Sprintf("api/v0/export/%v/", pk)
//
//   req, err := s.client.NewRequest("GET", url, nil)
//
// 	authstr := GetAuth()
// 	if authstr == "" {
// 		un, pw := Logon()
// 		req.SetBasicAuth(un, pw)
// 	} else {
// 		req.Header.Add("authorization", "Basic "+authstr)
// 	}
//   exportDetail := new(ExportDetail)
//   resp, err := s.client.Do(req, exportDetail)
// 	// resp, err := client.Do(req)
// 	// exports, err := ioutil.ReadAll(resp.Body)
// 	// defer resp.Body.Close()
// 	// Check(err)
// 	// a := &exportDetail{}
// 	// err = json.Unmarshal([]byte(exports), &a)
// 	// Check(err)
//
// 	durl := "/export/download/file/" + strconv.Itoa(pk) + "/"
// 	req1, err := http.NewRequest("GET", durl, nil)
// 	if authstr == "" {
// 		un, pw := Logon()
// 		req1.SetBasicAuth(un, pw)
// 	} else {
// 		req1.Header.Add("authorization", "Basic "+authstr)
// 	}
// 	resp2, err := client.Do(req1)
// 	cd := resp2.Header.Get("Content-Disposition")
// 	_, params, err := mime.ParseMediaType(cd)
// 	fname := params["filename"]
// 	fmt.Println(fname)
//
// 	file, err := os.Create(fname)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()
//
// 	numBytes, err := io.Copy(file, resp2.Body)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	log.Println("Downloaded", numBytes, "bytes to", fname)
// }

// func GetExportDetail(pk int) exportDetail {
// 	url := "api/v0/export/"
//
// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}
// 	client := &http.Client{Transport: tr}
//
// 	req, err := http.NewRequest("GET", url+strconv.Itoa(pk)+"/", nil)
//
// 	authstr := GetAuth()
// 	if authstr == "" {
// 		un, pw := Logon()
// 		req.SetBasicAuth(un, pw)
// 	} else {
// 		req.Header.Add("authorization", "Basic "+authstr)
// 	}
// 	resp, err := client.Do(req)
// 	exports, err := ioutil.ReadAll(resp.Body)
// 	resp.Body.Close()
// 	Check(err)
// 	a := &exportDetail{}
// 	err = json.Unmarshal([]byte(exports), &a)
// 	Check(err)
// 	return *a
// }
