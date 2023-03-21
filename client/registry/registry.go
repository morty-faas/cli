package registry

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	httpclient "morty/client"
	"morty/pkg/debug"
	"morty/pkg/serdejson"
	"net/http"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type (
	client struct {
		c *httpclient.Client
	}

	BuildFnRequest struct {
		Name    string `json:"name"`
		Runtime string `json:"runtime"`
		Archive string `json:"archive"`
	}
)

// NewClient initiate a new client for the Morty Registry
func NewClient(baseURL string) *client {
	return &client{httpclient.NewClient(baseURL)}
}

// BuildFn send a build request against the registry and returns the function URI in the registry.
func (rc *client) BuildFn(context context.Context, opts *BuildFnRequest) (string, error) {
	log.Debugf("New build request with options: %v", debug.JSON(opts))

	// Build multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("archive", filepath.Base(opts.Archive))
	if err != nil {
		return "", err
	}

	file, err := os.Open(opts.Archive)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(fw, file)
	if err != nil {
		return "", err
	}

	writer.WriteField("name", opts.Name)
	writer.WriteField("runtime", opts.Runtime)

	writer.Close()

	headers := http.Header{}
	headers.Add("Content-Type", writer.FormDataContentType())

	uri := path.Join("v1", "functions", "build")
	res, err := rc.c.Post(context, uri, body, headers)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resourceUri, err := serdejson.Deserialize[string](res.Body)
	return *resourceUri, err
}