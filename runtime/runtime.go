package runtime

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	runtimeTemplateEndpoint = "https://api.github.com/repos/polyxia-org/morty-runtimes/git/trees/main?recursive=1"
)

type (
	RuntimesResponse struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
		Tree []struct {
			Path string `json:"path"`
			Mode string `json:"mode"`
			Type string `json:"type"`
			Size int    `json:"size"`
			Sha  string `json:"sha"`
			URL  string `json:"url"`
		} `json:"tree"`
		Truncated bool `json:"truncated"`
	}
)

func List() ([]string, error) {
	log.Debugf("GET request on '%s'", runtimeTemplateEndpoint)
	response, err := http.Get(runtimeTemplateEndpoint)
	if err != nil {
		return nil,err
	}
	defer response.Body.Close()

	var runtimes RuntimesResponse
	err = json.NewDecoder(response.Body).Decode(&runtimes)
	if err != nil {
		return nil,err
	}

	var runtimesList []string

	for _, tree := range runtimes.Tree {
		if tree.Type == "tree" && strings.HasPrefix(tree.Path, "template/")  && strings.HasSuffix(tree.Path, "/function") {
			runtimesList = append(runtimesList, strings.Split(tree.Path, "/")[1])
		}
	}

	return runtimesList, nil
}