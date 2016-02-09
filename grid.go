package grid

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/gopass"
)

const (
	defaultBaseURL = "https://rsgis.erdc.dren.mil/te_ba/"
)

// A Client manages communication with the GRiD API.
type Client struct {
	// HTTP client used to communicate with the API.
	client    *http.Client
	transport *http.Transport

	// Base URL for API requests.  Defaults to GRiD TE, but can be
	// set to a domain endpoint to use with other instances.  BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// Services used for talking to different parts of the GRiD API.
	AOI    *AOIService
	Export *ExportService
	// File     *FileService
	// Download *DownloadService
}

// NewClient returns a new GRiD API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := &Client{client: httpClient, transport: tr, BaseURL: baseURL}
	c.AOI = &AOIService{client: c}
	c.Export = &ExportService{client: c}
	// c.File = &FileService{client: c}
	// c.Download = &DownloadService{client: c}
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// req.Header.Add("Accept", mediaTypeV3)
	// if c.UserAgent != "" {
	// 	req.Header.Add("User-Agent", c.UserAgent)
	// }
	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	// c.Rate = response.Rate

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	return response, err
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// Response is a GitHub API response.  This wraps the standard http.Response
// returned from GitHub and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	// response.populatePageValues()
	// response.populateRate()
	return response
}

/*
An ErrorResponse reports one or more errors caused by an API request.
GitHub API docs: http://developer.github.com/v3/#client-errors
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
	Errors   []Error        `json:"errors"`  // more detail on individual errors
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message, r.Errors)
}

// sanitizeURL redacts the client_id and client_secret tokens from the URL which
// may be exposed to the user, specifically in the ErrorResponse error message.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

/*
An Error reports more details on an individual error in an ErrorResponse.
These are the possible validation error codes:
    missing:
        resource does not exist
    missing_field:
        a required field on a resource has not been set
    invalid:
        the formatting of a field is invalid
    already_exists:
        another resource has the same valid as this field
GitHub API docs: http://developer.github.com/v3/#client-errors
*/
type Error struct {
	Resource string `json:"resource"` // resource on which the error occurred
	Field    string `json:"field"`    // field on which the error occurred
	Code     string `json:"code"`     // validation error code
}

/*
UnauthenticatedRateLimitedTransport allows you to make unauthenticated calls
that need to use a higher rate limit associated with your OAuth application.
	t := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     "your app's client ID",
		ClientSecret: "your app's client secret",
	}
	client := github.NewClient(t.Client())
This will append the querystring params client_id=xxx&client_secret=yyy to all
requests.
See http://developer.github.com/v3/#unauthenticated-rate-limited-requests for
more information.
*/
type TLSClient struct {
	Transport http.Transport
}

// Client returns an *http.Client that makes requests which are subject to the
// rate limit of your OAuth application.
func (t *TLSClient) Client() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

// BasicAuthTransport is an http.RoundTripper that authenticates all requests
// using HTTP Basic Authentication with the provided username and password.  It
// additionally supports users who have two-factor authentication enabled on
// their GitHub account.
type BasicAuthTransport struct {
	Username string // GitHub username
	Password string // GitHub password
	// OTP      string // one-time password for users with two-factor auth enabled

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.
func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// req = cloneRequest(req) // per RoundTrip contract
	req.SetBasicAuth(t.Username, t.Password)
	// if t.OTP != "" {
	// req.Header.Add(headerOTP, t.OTP)
	// }
	return t.transport().RoundTrip(req)
}

// Client returns an *http.Client that makes requests that are authenticated
// using HTTP Basic Authentication.
func (t *BasicAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *BasicAuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
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
