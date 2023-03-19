package function

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	templateEndpoint   = "github.com/polyxia-org/morty-runtimes.git//template"
	mortyWorkspaceFile = "morty.yaml"
)

type (
	Options struct {
		Name      string
		Runtime   string
		Directory string
	}
	function struct {
		Name    string `yaml:"name"`
		Runtime string `yaml:"runtime"`
	}
)

var (
	ErrFunctionDirectoryRequired = errors.New("directory is required")
	ErrFunctionDirectoryExists   = errors.New("directory already exists")
	ErrFunctionRuntimeInvalid    = errors.New("runtime not found")
)

// New initialize a new function on the disk with the given options
func New(opts *Options) error {
	if opts.Directory == "" {
		return ErrFunctionDirectoryRequired
	}

	log.Infof("Creating new function workspace into '%s'", opts.Directory)

	if err := downloadTemplate(opts); err != nil {
		return err
	}

	return injectWorkspaceFile(opts)
}

// downloadTemplate will download the template for this workspace based on the given runtime.
func downloadTemplate(opts *Options) error {
	templateUrl := templateEndpoint + "/" + opts.Runtime + "/function"
	log.Debugf("Downloading remote folder: '%s' into '%s'", templateUrl, opts.Directory)

	if err := getter.Get(opts.Directory, templateUrl); err != nil {
		// If an error occurs, the directory will be created
		// Try to remove it, but don't go in error for this
		if err := os.Remove(opts.Directory); err != nil {
			log.Errorf("failed to remove directory %s: %v", opts.Directory, err)
		}

		if strings.Contains(err.Error(), "not found") {
			return ErrFunctionRuntimeInvalid
		}
		return err
	}

	return nil
}

// injectWorkspaceFile inject the morty.yaml file into a workspace directory
func injectWorkspaceFile(opts *Options) error {
	log.Debugf("Injecting '%s' into workspace directory '%s'", mortyWorkspaceFile, opts.Directory)

	f := &function{
		Name:    opts.Name,
		Runtime: opts.Runtime,
	}

	by, err := yaml.Marshal(f)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(opts.Directory, mortyWorkspaceFile), by, 0644)
}
