package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type (
	Gateway struct {
		url string
	}

	CreateFunctionRequest struct {
		Name   string `json:"name"`
		Rootfs string `json:"rootfs"`
	}

	CreateFunctionResponse struct {
		Message string `json:"message"`
	}
)

func (g *Gateway) CreateFunction(name, rootfs string) error {
	requestBody, err := json.Marshal(CreateFunctionRequest{Name: name, Rootfs: rootfs})
	if err != nil {
		return err
	}
	res, err := http.Post(g.url+"/functions", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Could not POST to gateway")
		return err
	}
	if res.StatusCode >= 400 {
		body := CreateFunctionResponse{}
		json.NewDecoder(res.Body).Decode(&body)
		log.Error(body.Message)
		return fmt.Errorf("The function could not be created\n")
	}

	return nil
}
