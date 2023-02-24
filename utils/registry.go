package utils 

import (
	"os"
	"net/http"
	"log"
	"fmt"
	"bytes"
	"encoding/json"
)

type Registry struct {
	url string
}

func NewRegistry() *Registry {
	url := os.Getenv("MORTY_REGISTRY_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	return &Registry{
		url,
	}
}

func (r *Registry) GetPresignedUrl(name string) string {
	endpoint := fmt.Sprintf("%s/v1/functions/%s/upload-link", r.url, name)

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal("ERROR: Cannot connect to registry.", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Unable to get presigned url: %v", err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

  presignedUrl := PresignedUrl{}
	json.NewDecoder(res.Body).Decode(&presignedUrl)
	return presignedUrl.UploadLink
}

func (r *Registry) UploadFile(name, filePath string) {
	presignedUrl := r.GetPresignedUrl(name)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Unable to open %s: %v", filePath, err)
	}
	fileBuffer := bytes.NewBuffer(file)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", presignedUrl, fileBuffer)
	if err != nil {
		log.Fatalf("Unable to create PUT request: %v", err)
	}
	_, err = client.Do(req)
	if err != nil {
		log.Fatalf("Unable to PUT: %v", err)
	}
}

type PresignedUrl struct {
	HttpMethod string `json:"httpMethod"`
	UploadLink string `json:"uploadLink"`
}