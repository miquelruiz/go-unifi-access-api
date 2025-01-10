package unifiaccessclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/miquelruiz/go-unifi-access-api/schema"
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	token      string
	baseUrl    url.URL
	httpClient HttpClient
}

// New creates a new API client that will send requests to host using token for
// authentication.
//
// A default [http.Client] will be created and used. Use [NewWithHttpClient] to
// provide a custom [http.Client].
//
// An error will be returned if host can't be parsed by [url.Parse].
func New(host string, token string) (*Client, error) {
	return NewWithHttpClient(
		host,
		token,
		&http.Client{},
	)
}

// NewWithHttpClient is the same than [New] but accepts a custom HttpClient
func NewWithHttpClient(host string, token string, httpClient HttpClient) (*Client, error) {
	baseUrl, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return &Client{
		token:      token,
		baseUrl:    *baseUrl,
		httpClient: httpClient,
	}, nil
}

func (c *Client) CreateUser(u schema.UserRequest) (*schema.UserResponse, error) {
	url := "/api/v1/developer/users"
	buffer := bytes.NewBuffer(make([]byte, 0))
	d := json.NewEncoder(buffer)
	if err := d.Encode(&u); err != nil {
		return nil, err
	}
	req := c.buildRequest(http.MethodPost, url, "", io.NopCloser(buffer))
	return doRequest[schema.UserResponse](c, req)
}

func (c *Client) GetUser(id string) (*schema.UserResponse, error) {
	url := "/api/v1/developer/users/" + id
	req := c.buildRequest(http.MethodGet, url, "", nil)
	return doRequest[schema.UserResponse](c, req)
}

func (c *Client) UpdateUser(id string, u schema.UserRequest) error {
	url := "/api/v1/developer/users/" + id
	buffer := bytes.NewBuffer(make([]byte, 0))
	d := json.NewEncoder(buffer)
	if err := d.Encode(&u); err != nil {
		return err
	}

	req := c.buildRequest(http.MethodPut, url, "", io.NopCloser(buffer))
	_, err := doRequest[schema.UserResponse](c, req)
	return err
}

// TODO deal with pagination
func (c *Client) ListUsers() ([]schema.UserResponse, error) {
	url := "/api/v1/developer/users"
	req := c.buildRequest(http.MethodGet, url, "", nil)
	listPtr, err := doRequest[[]schema.UserResponse](c, req)
	if err != nil {
		return nil, err
	}
	return *listPtr, nil
}

func (c *Client) buildRequest(method string, path string, rawQuery string, body io.ReadCloser) *http.Request {
	baseCopy := c.baseUrl
	baseCopy.Path = path
	baseCopy.RawQuery = rawQuery

	return &http.Request{
		Method: method,
		URL:    &baseCopy,
		Header: map[string][]string{"Authorization": {"Bearer " + c.token}},
		Body:   body,
	}
}

func doRequest[T any](c *Client, req *http.Request) (*T, error) {
	rawresp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rawresp.Body.Close()

	if rawresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("endpoint returned %s", rawresp.Status)
	}

	resp, err := parseResponse[T](rawresp.Body)
	if err != nil {
		return nil, err
	}

	if resp.Code != schema.ResponseSuccess {
		return nil, errors.New(resp.Code + ": " + resp.Msg)
	}

	return &resp.Data, nil
}

func parseResponse[T any](r io.Reader) (*schema.Response[T], error) {
	var resp schema.Response[T]
	d := json.NewDecoder(r)
	if err := d.Decode(&resp); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	return &resp, nil
}
