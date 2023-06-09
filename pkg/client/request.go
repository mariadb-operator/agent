package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	jsonMediaType = "application/json"
)

func (c *Client) newRequestWithContext(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	baseUrl, err := buildURL(*c.baseUrl, path)
	if err != nil {
		return nil, fmt.Errorf("error building URL: %v", err)
	}

	setHeaders := func(r *http.Request) {
		r.Header.Set("Content-Type", jsonMediaType)
		r.Header.Set("Accept", jsonMediaType)
		for k, v := range c.headers {
			r.Header.Set(k, v)
		}
	}

	if method == http.MethodGet {
		req, err := http.NewRequestWithContext(ctx, method, baseUrl.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("error creating GET request: %v", err)
		}
		setHeaders(req)
		return req, nil
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling body: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseUrl.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	setHeaders(req)
	return req, nil
}

func buildURL(baseUrl url.URL, path string) (*url.URL, error) {
	baseUrl.Path = strings.TrimSuffix(baseUrl.Path, "/")
	baseUrl.Path += path

	newUrl, err := url.Parse(baseUrl.String())
	if err != nil {
		return nil, fmt.Errorf("error building URL: %v", err)
	}
	return newUrl, nil
}
