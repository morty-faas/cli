package function

import (
	"morty/pkg/archive"
	"os"

	log "github.com/sirupsen/logrus"
)

type BuildOptions struct {
	Directory string
	Registry  string
	Gateway   string
}

// Create a zip archive from a folder, upload it to the registry and create the function through the gateway
func Build(opts *BuildOptions) (string, error) {
	config, err := functionFromFile(opts.Directory)
	if err != nil {
		return "", err
	}
	name := config.Name
	runtime := config.Runtime
	log.Infof("Fonction configuration: name=%s, runtime=%s\n", name, runtime)

	zipPath := "/tmp/morty/"
	if err = os.MkdirAll(zipPath, os.ModePerm); err != nil {
		return "", err
	}
	log.Debugf("Created %s folder\n", zipPath)

	zipFilename := zipPath + "function.zip"

	if err = archive.Zip(opts.Directory, zipFilename); err != nil {
		return "", err
	}
	log.Debugf("Created zip archive: %s\n", zipFilename)

	registry := new(Registry)
	registry.url = opts.Registry
	log.Infof("Registry URL set to %s\n", registry.url)
	imageUrl, err := registry.UploadFile(name, runtime, zipFilename)
	if err != nil {
		return "", err
	}
	imageUrl = registry.url + imageUrl
	log.Infof("Function %s has been uploaded to registry (imageUrl: %s)\n", name, imageUrl)

	if err = os.Remove(zipFilename); err != nil {
		return "", err
	}
	log.Debugf("Removed %s\n", zipFilename)

	gateway := new(Gateway)
	gateway.url = opts.Gateway
	log.Infof("Gateway URL set to %s\n", gateway.url)
	if err = gateway.CreateFunction(name, imageUrl); err != nil {
		return "", err
	}
	log.Infof("Function %s has been created through the gateway\n", name)

	return name, nil
}
