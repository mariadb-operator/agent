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

type requestOptions struct {
	headers map[string]string
	query   map[string]string
	body    interface{}
}

type requestOption func(*requestOptions)

func withHeaders(headers map[string]string) requestOption {
	return func(ro *requestOptions) {
		ro.headers = headers
	}
}

func withQuery(query map[string]string) requestOption {
	return func(ro *requestOptions) {
		ro.query = query
	}
}

func withBody(body interface{}) requestOption {
	return func(ro *requestOptions) {
		ro.body = body
	}
}

func (c *Client) newRequestWithContext(ctx context.Context, method, path string, reqOpts ...requestOption) (*http.Request, error) {
	opts := requestOptions{}
	for _, setOpt := range reqOpts {
		setOpt(&opts)
	}
	baseUrl, err := buildURL(*c.baseUrl, path, opts.query)
	if err != nil {
		return nil, fmt.Errorf("error building URL: %v", err)
	}

	setHeaders := func(r *http.Request) {
		r.Header.Set("Content-Type", jsonMediaType)
		r.Header.Set("Accept", jsonMediaType)
		for k, v := range c.headers {
			r.Header.Set(k, v)
		}
		for k, v := range opts.headers {
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
	if opts.body != nil {
		bodyBytes, err := json.Marshal(opts.body)
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

func (c *Client) newRequest(method, path string, reqOpts ...requestOption) (*http.Request, error) {
	return c.newRequestWithContext(context.Background(), method, path, reqOpts...)
}

func buildURL(baseUrl url.URL, path string, query map[string]string) (*url.URL, error) {
	if !strings.HasSuffix(baseUrl.Path, "/") {
		baseUrl.Path += "/"
	}
	baseUrl.Path += path

	if query != nil {
		q := baseUrl.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		baseUrl.RawQuery = q.Encode()
	}

	newUrl, err := url.Parse(baseUrl.String())
	if err != nil {
		return nil, fmt.Errorf("error building URL: %v", err)
	}
	return newUrl, nil
}
