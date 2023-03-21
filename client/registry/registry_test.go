package registry

import (
	"context"
	"encoding/json"
	"morty/pkg/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BuildFn_ResourceURIIsReturned(t *testing.T) {
	expected := "/v1/functions/hello-world"
	server := testutils.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	client := NewClient(server.URL)

	opts := &BuildFnRequest{
		Name:    "hello-world",
		Runtime: "node-19",
		Archive: "../../testdata/function.zip",
	}

	uri, err := client.BuildFn(context.TODO(), opts)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, uri)
	assert.NoError(t, err)
}

func Test_BuildFn_ErrorIfArchiveNotFound(t *testing.T) {
	client := NewClient("http://test.com")

	opts := &BuildFnRequest{
		Name:    "hello-world",
		Runtime: "node-19",
		Archive: "super-nano-jetson.zip",
	}

	_, err := client.BuildFn(context.TODO(), opts)

	assert.Error(t, err)
}
