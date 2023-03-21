package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	httpclient "morty/client"
	"morty/pkg/debug"
	"morty/pkg/serdejson"
	"net/http"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

type (
	client struct {
		c *httpclient.Client
	}

	APIError struct {
		Message string `json:"message"`
	}

	CreateFnRequest struct {
		Name   string `json:"name"`
		Rootfs string `json:"rootfs"`
	}

	CreateFnResponse = CreateFnRequest

	InvokeFnRequest struct {
		FnName  string   `json:"functionName"`
		Method  string   `json:"method"`
		Body    string   `json:"body"`
		Headers []string `json:"headers"`
	}
)

// NewClient initiate a new client for the Morty Gateway
func NewClient(baseURL string) *client {
	return &client{httpclient.NewClient(baseURL)}
}

// InvokeFn invoke a function and return the resulting payload.
func (gc *client) InvokeFn(context context.Context, opts *InvokeFnRequest) (string, error) {
	log.Debugf("New invocation request with options: %v", debug.JSON(opts))

	headers := http.Header{}
	// If the caller has passed headers, map them to http.Header
	if opts.Headers != nil {
		for _, header := range opts.Headers {
			splitted := strings.Split(header, ":")
			if len(splitted) != 2 {
				return "", fmt.Errorf("header '%s' is not valid. Please use the correct format: 'Key: Value'", header)
			}
			hKey, hValue := splitted[0], splitted[1]
			headers.Add(hKey, hValue)
		}
	}

	var body io.Reader
	if opts.Body != "" {
		body = bytes.NewBuffer([]byte(opts.Body))
	}

	uri := path.Join("invoke", opts.FnName)
	res, err := gc.c.Generic(context, opts.Method, uri, body, headers)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		// If the API has returned an error, map it to a custom error type APIError
		return "", makeApiError(res.Body)
	}

	fnPayload, err := serdejson.Deserialize[string](res.Body)
	return *fnPayload, err
}

// CreateFn create a function and return an error if there is one
func (gc *client) CreateFn(context context.Context, dto *CreateFnRequest) (*CreateFnResponse, error) {
	log.Debugf("Creating function with payload: %v", debug.JSON(dto))

	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	uri := "functions"
	res, err := gc.c.Post(context, uri, bytes.NewBuffer(body), headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		// If the API has returned an error, map it to a custom error type APIError
		return nil, makeApiError(res.Body)
	}

	return serdejson.Deserialize[CreateFnResponse](res.Body)
}

// Build an APIError from the response body
func makeApiError(body io.Reader) error {
	apiErr := &APIError{}
	if err := json.NewDecoder(body).Decode(apiErr); err != nil {
		return err
	}
	return apiErr
}

// Error is the implementation of the error interface for our custom APIError type
func (e *APIError) Error() string {
	return e.Message
}
