package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL   = "https://api.form3.tech/"
	apiVersion       = 1
	userAgent        = "go-form3"
	jsonApiMediaType = "application/vnd.api+json"
)

// A Client manages communication with the Form3 API.
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent used when communicating with the Form3 API.
	UserAgent string

	common service

	Accounts *AccountsService
}

type service struct {
	client *Client
}

// NewClient returns a new Form3 API client. If a nil httpClient is
// provided, a new http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c
	c.Accounts = (*AccountsService)(&c.common)
	return c
}

// ListOptions specifies the optional parameters to various List methods that
// support offset pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	PageNumber int `url:"page[number],omitempty"`

	// For paginated result sets, the number of results to include per page.
	PageSize int `url:"page[size],omitempty"`
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", jsonApiMediaType)
	}
	req.Header.Set("Accept", jsonApiMediaType)
	req.Header.Set("Date", time.Now().String())
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// Response is a Form3 API response. This wraps the standard http.Response
// returned from Form3
type Response struct {
	*http.Response
}

// Response content when status code is 400 Bad Request
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"error_message"` // error message
	Code     string         `json:"error_code"`    // more detail on individual errors
}

type Links struct {
	// Link to this endpoint or resource.
	Self *string `json:"self"`
	// Link to the first page of a paginated response. Only available in list call responses.
	First *string `json:"first"`
	// Link to the last page of a paginated response. Only available in list call responses.
	Last *string `json:"last"`
	// Link to the next page of a paginated response. Only available in list call responses.
	Next *string `json:"next"`
	// Link to the previous page of a paginated response. Only available in list call responses.
	Prev *string `json:"prev"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Code)
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range
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

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }
