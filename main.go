package grid

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
)

type selfAoiList struct {
	Name       string `json:"name"`
	Geometry   string `json:"geometry"`
	IsActive   bool   `json:"is_active"`
	Source     string `json:"source"`
	NumExports string `json:"num_exports"`
	Pk         int    `json:"pk"`
	CreatedAt  string `json:"created_at"`
}

type aoi struct {
	SelfAoiList []selfAoiList `json:"self_aoi_list"`
}

const (
	baseURL = "https://rsgis.erdc.dren.mil/te_ba"
	gridURL = baseURL + "/api/v0"
)

func GetAOIList() aoi {
	url := gridURL + "/aoi/"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)

	authstr := GetAuth()
	if authstr == "" {
		un, pw := Logon()
		req.SetBasicAuth(un, pw)
	} else {
		req.Header.Add("authorization", "Basic "+authstr)
	}
	resp, err := client.Do(req)
	aois, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	Check(err)
	a := &aoi{}
	err = json.Unmarshal([]byte(aois), &a)
	Check(err)
	return *a
}

type exportSet struct {
	Status    string `json:"status"`
	StartedAt string `json:"stated_at"`
	Name      string `json:"name"`
	Datatype  string `json:"datatype"`
	HSRS      int    `json:"hsrs"`
	URL       string `json:"url"`
	Pk        int    `json:"pk"`
}

type aoiDetail struct {
	ExportSet []exportSet `json:"export_set"`
}

func GetAOIDetail(pk int) aoiDetail {
	url := gridURL + "/aoi/"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url+strconv.Itoa(pk)+"/", nil)

	authstr := GetAuth()
	if authstr == "" {
		un, pw := Logon()
		req.SetBasicAuth(un, pw)
	} else {
		req.Header.Add("authorization", "Basic "+authstr)
	}
	resp, err := client.Do(req)
	exports, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	Check(err)
	a := &aoiDetail{}
	err = json.Unmarshal([]byte(exports), &a)
	Check(err)
	return *a
}

type exportFiles struct {
	URL  string `json:"url"`
	Pk   int    `json:"pk"`
	Name string `json:"name"`
}

type exportDetail struct {
	ExportFiles []exportFiles `json:"exportfiles"`
}

func GetExportDetail(pk int) exportDetail {
	url := gridURL + "/export/"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url+strconv.Itoa(pk)+"/", nil)

	authstr := GetAuth()
	if authstr == "" {
		un, pw := Logon()
		req.SetBasicAuth(un, pw)
	} else {
		req.Header.Add("authorization", "Basic "+authstr)
	}
	resp, err := client.Do(req)
	exports, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	Check(err)
	a := &exportDetail{}
	err = json.Unmarshal([]byte(exports), &a)
	Check(err)
	return *a
}

func DownloadByPK(pk int) {
	url := gridURL + "/export/"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url+strconv.Itoa(pk)+"/", nil)

	authstr := GetAuth()
	if authstr == "" {
		un, pw := Logon()
		req.SetBasicAuth(un, pw)
	} else {
		req.Header.Add("authorization", "Basic "+authstr)
	}
	resp, err := client.Do(req)
	exports, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	Check(err)
	a := &exportDetail{}
	err = json.Unmarshal([]byte(exports), &a)
	Check(err)

	durl := baseURL + "/export/download/file/" + strconv.Itoa(pk) + "/"
	req1, err := http.NewRequest("GET", durl, nil)
	if authstr == "" {
		un, pw := Logon()
		req1.SetBasicAuth(un, pw)
	} else {
		req1.Header.Add("authorization", "Basic "+authstr)
	}
	resp2, err := client.Do(req1)
	cd := resp2.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(cd)
	fname := params["filename"]
	fmt.Println(fname)

	file, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	numBytes, err := io.Copy(file, resp2.Body)
	if err != nil {
		panic(err)
	}

	log.Println("Downloaded", numBytes, "bytes to", fname)
}

func Logon() (un, pw string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	un, _ = reader.ReadString('\n')
	un = strings.TrimSpace(un)
	fmt.Print("Enter password: ")
	pass := gopass.GetPasswd()
	pw = string(pass)
	return
}

func GetAuth() string {
	path := os.Getenv("HOME") + string(filepath.Separator) + ".grid"
	fileandpath := path + string(filepath.Separator) + "credentials"
	file, err := os.Open(fileandpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	return line
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
