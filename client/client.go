package client

import (
	"context"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	baseURL string
	inner   *http.Client
}

// NewClient returns a new custom HTTP client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		inner:   &http.Client{},
	}
}

// Get perform a GET request using the given options
func (c *Client) Get(context context.Context, path string, headers http.Header) (*http.Response, error) {
	return c.Generic(context, http.MethodGet, path, nil, headers)
}

// Post perform a POST request using the given options
func (c *Client) Post(context context.Context, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	return c.Generic(context, http.MethodPost, path, body, headers)
}

// Patch perform a PATCH request using the given options
func (c *Client) Patch(context context.Context, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	return c.Generic(context, http.MethodPatch, path, body, headers)
}

// Put perform a PUT request using the given options
func (c *Client) Put(context context.Context, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	return c.Generic(context, http.MethodPut, path, body, headers)
}

// Delete perform a DELETE request using the given options
func (c *Client) Delete(context context.Context, path string, headers http.Header) (*http.Response, error) {
	return c.Generic(context, http.MethodDelete, path, nil, headers)
}

// Generic perform a request using the given options
func (c *Client) Generic(context context.Context, method, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	req, err := c.makeRequest(context, method, path, body, headers)
	if err != nil {
		return nil, err
	}
	log.Debugf("Sending %s request on '%s'", req.Method, req.URL)
	return c.inner.Do(req)
}

// makeRequest creates a new http.Request usign the given parameters
func (c *Client) makeRequest(context context.Context, method, uri string, body io.Reader, headers http.Header) (*http.Request, error) {
	url := c.baseURL + "/" + uri

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	request.Header = headers
	return request.WithContext(context), nil
}
