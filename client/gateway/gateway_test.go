package gateway

import (
	"context"
	"encoding/json"
	"morty/pkg/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateFn_CreateFunctionOk(t *testing.T) {
	expected := &CreateFnResponse{
		Name:   "hello-world",
		Rootfs: "http://url-to-rootfs.go",
	}

	s := testutils.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})
	defer s.Close()

	client := NewClient(s.URL)

	res, err := client.CreateFn(context.TODO(), &CreateFnRequest{
		Name:   "hello-world",
		Rootfs: "http://url-to-rootfs.go",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func Test_CreateFn_ErrorCorrectlyParsed(t *testing.T) {
	expected := &APIError{
		Message: "an error occurred",
	}

	s := testutils.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(expected)
	})
	defer s.Close()

	client := NewClient(s.URL)

	_, err := client.CreateFn(context.TODO(), &CreateFnRequest{})

	assert.Error(t, err)
	assert.Equal(t, expected, err)
}

func Test_InvokeFn_InvokeFnOk(t *testing.T) {
	opts := &InvokeFnRequest{
		FnName:  "get-weather",
		Method:  "GET",
		Body:    "",
		Headers: nil,
	}

	expected := `{"city": "Montpellier", "temp_celsius": 30}`

	s := testutils.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	})
	defer s.Close()

	client := NewClient(s.URL)

	res, err := client.InvokeFn(context.TODO(), opts)
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}
