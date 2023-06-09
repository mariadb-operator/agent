package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mariadb-operator/agent/pkg/errors"
)

const (
	jsonMediaType = "application/json"
)

type Option func(*Client) error

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		if httpClient == nil {
			httpClient = http.DefaultClient
		}

		c.httpClient = httpClient
		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		if timeout == 0 {
			timeout = 1 * time.Minute
		}

		c.httpClient.Timeout = timeout
		return nil
	}
}

type Client struct {
	Bootstrap   *Bootstrap
	GaleraState *GaleraState
	Recovery    *Recovery

	baseUrl    *url.URL
	httpClient *http.Client
	headers    map[string]string
}

func NewClient(baseUrl string, opts ...Option) (*Client, error) {
	url, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}
	client := &Client{
		baseUrl:    url,
		httpClient: http.DefaultClient,
		headers:    make(map[string]string, 0),
	}
	for _, setOpt := range opts {
		if err := setOpt(client); err != nil {
			return nil, fmt.Errorf("error setting option: %v", err)
		}
	}

	client.Bootstrap = &Bootstrap{
		Client: client,
	}
	client.GaleraState = &GaleraState{
		Client: client,
	}
	client.Recovery = &Recovery{
		Client: client,
	}
	return client, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error doing request: %v", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	if res.StatusCode >= 400 {
		var apiErr errors.APIError
		if err := decoder.Decode(&apiErr); err != nil {
			return fmt.Errorf("error decoding body into error: %v", err)
		}
		return errors.NewError(res.StatusCode, apiErr.Error())
	}

	if v == nil {
		return nil
	}
	if err := decoder.Decode(&v); err != nil {
		return fmt.Errorf("error decoding body: %v", err)
	}
	return nil
}
