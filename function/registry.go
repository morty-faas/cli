package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Registry struct {
	url string
}

type BuildResponseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Uploads a file to the registry and returns the url of the image
func (r *Registry) UploadFile(name, runtime, filePath string) (string, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	part, err := writer.CreateFormFile("archive", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	writer.WriteField("name", name)
	writer.WriteField("runtime", runtime)

	writer.Close()

	log.Infof("Uploading %s to registry\n", name)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/functions/build", r.url), requestBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode >= 400 {
		error := BuildResponseError{}
		json.NewDecoder(res.Body).Decode(&error)
		log.Error(error)
		return "", fmt.Errorf("The function could not be pushed\n")
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Unable to read response body")
		return "", err
	}

	var uri string
	if err = json.Unmarshal(responseBody, &uri); err != nil {
		return "", err
	}
	return uri, nil
}
