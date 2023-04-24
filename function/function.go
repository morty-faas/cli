package function

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/morty-faas/cli/pkg/debug"

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
		Path    string `yaml:"-"`
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
func New(opts *Options) (*function, error) {
	if opts.Directory == "" {
		return nil, ErrFunctionDirectoryRequired
	}

	log.Infof("Creating new function workspace into '%s'", opts.Directory)

	if err := downloadTemplate(opts); err != nil {
		return nil, err
	}

	if err := injectWorkspaceFile(opts); err != nil {
		return nil, err
	}

	f := &function{
		Name:    opts.Name,
		Runtime: opts.Runtime,
		Path:    opts.Directory,
	}

	log.Debugf("Function created: %v", debug.JSON(f))

	return f, nil
}

// NewFromFile will read the Morty workspace file and return a function
func NewFromFile(folder string) (*function, error) {
	path := filepath.Join(folder, mortyWorkspaceFile)
	log.Debugf("Reading function config file: %s\n", path)

	byteValue, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("Unable to read %s", path)
		return nil, err
	}
	f := &function{}
	yaml.Unmarshal([]byte(byteValue), f)
	f.Path = folder

	log.Debugf("Function loaded from file %s: %v", path, debug.JSON(f))

	return f, nil
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
