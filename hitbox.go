// Package hitbox provides a client for using the hitbox.tv API.
package hitbox

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

const (
	librayVersion  = "0.1"
	defaultBaseUrl = "http://api.hitbox.tv/"
	userAgent      = "hitbox/" + librayVersion
)

// GamesResponse represents the response from a request on the Games endpoint.
type GamesResponse struct {
	Request    Request    `json:"request,omitempty"`
	Categories []Category `json:"categories,omitempty"`
}

// Request represents the original request as it was received by the hitbox.tv API.
type Request struct {
	This string `json:"this,omitempty"`
}

// Category represents a game category.
type Category struct {
	CategoryID        int64     `json:"category_id,omitempty,string"`
	CategoryName      string    `json:"category_name,omitempty"`
	CategoryNameShort string    `json:"category_name_short,omitempty"`
	CategorySeoKey    string    `json:"category_seo_key,omitempty"`
	CategoryViewers   int64     `json:"category_viewers,omitempty,string"`
	CategoryLogoSmall string    `json:"category_logo_small,omitempty"`
	CategoryLogolarge string    `json:"category_logo_large,omitempty"`
	CategoryUpdated   Timestamp `json:"category_updated,omitempty"`
}

// Timestamp represents a hitbox.tv time that can be unmarshalled from a JSON string.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Time is expected in the format yyyy-mm-dd hh:mm:ss.
func (t *Timestamp) UnmarshalJSON(data []byte) (err error) {
	str := string(data[1 : len(data)-1])
	(*t).Time, err = time.Parse("2006-01-02 15:04:05", str)
	return
}

// A Client manages communication with the hitbox.tv API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent used when communicating with the hitbox.tv API.
	UserAgent string
}

// NewClient returns a new hitbox.tv API client. If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseUrl, _ := url.Parse(defaultBaseUrl)
	c := &Client{client: httpClient, BaseURL: baseUrl, UserAgent: userAgent}
	return c
}

// Games returns a list of games objects sorted by number of current viewers on hitbox.
func (c *Client) Games() (*GamesResponse, *http.Response, error) {
	urlStr := "games"
	req, err := c.newRequest("GET", urlStr)
	if err != nil {
		return nil, nil, err
	}

	gResp := new(GamesResponse)
	resp, err := c.do(req, gResp)
	if err != nil {
		return nil, resp, err
	}

	return gResp, resp, err

}

// newRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) newRequest(method, urlStr string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp, errors.New("hitbox.tv api responded with http " + string(resp.StatusCode))
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}
