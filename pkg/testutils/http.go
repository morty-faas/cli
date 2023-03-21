package testutils

import (
	"net/http"
	"net/http/httptest"
)

// NewTestServer creates a new HTTP test server to be run in unit tests
func NewTestServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}
