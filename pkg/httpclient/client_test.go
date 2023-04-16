package httpclient

import (
	"context"
	"morty/pkg/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	t.Parallel()

	s := testutils.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer s.Close()

	client := NewClient(s.URL)
	res, err := client.Get(context.TODO(), "/", http.Header{})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.MethodGet, res.Request.Method)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_Post(t *testing.T) {
	t.Parallel()

	s := testutils.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer s.Close()

	client := NewClient(s.URL)
	res, err := client.Post(context.TODO(), "/", nil, http.Header{})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.MethodPost, res.Request.Method)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_Put(t *testing.T) {
	t.Parallel()

	s := testutils.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer s.Close()

	client := NewClient(s.URL)
	res, err := client.Put(context.TODO(), "/", nil, http.Header{})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.MethodPut, res.Request.Method)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_Patch(t *testing.T) {
	t.Parallel()

	s := testutils.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer s.Close()

	client := NewClient(s.URL)
	res, err := client.Patch(context.TODO(), "/", nil, http.Header{})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.MethodPatch, res.Request.Method)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_Delete(t *testing.T) {
	t.Parallel()

	s := testutils.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer s.Close()

	client := NewClient(s.URL)
	res, err := client.Delete(context.TODO(), "/", http.Header{})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.MethodDelete, res.Request.Method)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_makeRequest_createTheRequestWithNoError(t *testing.T) {
	t.Parallel()

	c := NewClient("http://example.com")
	_, err := c.makeRequest(context.TODO(), http.MethodGet, "/hello", nil, http.Header{})
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}
