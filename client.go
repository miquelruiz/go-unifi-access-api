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

type Client struct {
	host       string
	port       uint16
	token      string
	httpClient *http.Client
}

func New(host string, port uint16, token string) *Client {
	return NewWithHttpClient(
		host,
		port,
		token,
		&http.Client{},
	)
}

func NewWithHttpClient(host string, port uint16, token string, httpClient *http.Client) *Client {
	return &Client{
		host:       host,
		port:       port,
		token:      token,
		httpClient: httpClient,
	}
}

func (c *Client) CreateUser(u schema.UserRequest) (schema.UserResponse, error) {
	url := "/api/v1/developer/users"
	buffer := bytes.NewBuffer(make([]byte, 0))
	d := json.NewEncoder(buffer)
	if err := d.Encode(&u); err != nil {
		return schema.UserResponse{}, err
	}
	req := c.buildRequest(http.MethodPost, url, "", io.NopCloser(buffer))
	return doRequest[schema.UserResponse](c, req)
}

func (c *Client) GetUser(id string) (schema.UserResponse, error) {
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
	return doRequest[[]schema.UserResponse](c, req)
}

func (c *Client) buildRequest(method string, path string, rawQuery string, body io.ReadCloser) *http.Request {
	return &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme:   "https",
			Host:     c.host + ":" + fmt.Sprintf("%d", c.port),
			Path:     path,
			RawQuery: rawQuery,
		},
		Header: map[string][]string{"Authorization": {"Bearer " + c.token}},
		Body:   body,
	}
}

func doRequest[T any](c *Client, req *http.Request) (T, error) {
	var zero T
	rawresp, err := c.httpClient.Do(req)
	if err != nil {
		return zero, err
	}
	defer rawresp.Body.Close()

	resp, err := parseResponse[T](rawresp.Body)
	if err != nil {
		return zero, err
	}

	if resp.Code != schema.ResponseSuccess {
		return zero, errors.New(resp.Code + ": " + resp.Msg)
	}

	return resp.Data, nil

}

func parseResponse[T any](r io.Reader) (*schema.Response[T], error) {
	var resp schema.Response[T]
	d := json.NewDecoder(r)
	if err := d.Decode(&resp); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	return &resp, nil
}
