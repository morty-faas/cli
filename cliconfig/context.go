package cliconfig

import (
	"errors"
	"io/fs"
	"os"
	"os/user"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type CtxKey struct{}
type CurrentCtxKey struct{}
type ControllerClientContextKey struct{}

const (
	mortyConfigDefaultLocation = "$HOME/.morty/config.yaml"
	mortyConfigEnvVarKey       = "MORTYCONFIG"
)

var (
	ErrContextNotFound              = errors.New("context not found")
	ErrContextAlreadyExistsWithName = errors.New("context already exists with the same name")
)

type (
	Config struct {
		// The location of the Configuration on disk. Useful to save back the Configuration.
		// Will not be marshaled
		location string
		// The current active context
		Current string `yaml:"current"`
		// The available contexts
		Contexts []Context `yaml:"contexts"`
	}
	// Context holds information about a single Morty instance.
	Context struct {
		// The name of the current context
		Name string `yaml:"name"`
		// The address of the controller
		Controller string `yaml:"controller"`
		// The address of the registry
		Registry string `yaml:"registry"`
	}
)

// GetCurrentContext returns the current context
func (c *Config) GetCurrentContext() (*Context, error) {
	return c.GetContext(c.Current)
}

// GetContext returns the context associated with the given name.
func (c *Config) GetContext(name string) (*Context, error) {
	for _, ctx := range c.Contexts {
		if ctx.Name == name {
			return &ctx, nil
		}
	}
	return nil, ErrContextNotFound
}

// UseContext set the current context to the given context
func (c *Config) UseContext(name string) error {
	if !c.hasContext(name) {
		return ErrContextNotFound
	}

	c.Current = name
	return nil
}

// AddContext add the given context to the list of known contexts
func (c *Config) AddContext(context *Context) error {
	if c.hasContext(context.Name) {
		return ErrContextAlreadyExistsWithName
	}

	context.Controller = sanitizeUrl(context.Controller)
	context.Registry = sanitizeUrl(context.Registry)

	c.Contexts = append(c.Contexts, *context)
	return nil
}

// hasContext return true if the given context exists, false otherwise
func (c *Config) hasContext(name string) bool {
	for _, ctx := range c.Contexts {
		if ctx.Name == name {
			return true
		}
	}
	return false
}

// defaultConfig returns the default Configuration
func defaultConfig() *Config {
	return &Config{
		location: sanitizeDefaultConfigLocation(),
		Current:  "localhost",
		Contexts: []Context{
			{
				Name:       "localhost",
				Controller: "http://localhost:8080",
				Registry:   "http://localhost:8081",
			},
		},
	}
}

func (c *Config) Save() error {
	log.Debugf("Saving configuration to %s", c.location)
	by, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(c.location, by, fs.FileMode(os.O_WRONLY))
}

// Load the Configuration from the host environment.
func Load() (*Config, error) {
	// By default, we want to load the Configuration from the default location.
	// But if the user has specified a custom location using the MORTYConfig environment variable,
	// we will try to load the Configuration from it.
	path, pathFromEnv := sanitizeDefaultConfigLocation(), os.Getenv(mortyConfigEnvVarKey)
	if pathFromEnv != "" {
		path = pathFromEnv
	}

	log.Infof("Loading configuration from path: %s", path)

	v := viper.New()
	v.SetConfigFile(path)

	if err := v.MergeInConfig(); err != nil {
		// If the Configuration file is not found, and if we use the default location, then we can create
		// the Configuration file on the disk for future operations
		if path == sanitizeDefaultConfigLocation() {
			log.Debug("Configuration file doesn't exist on disk, initializing a new one")
			return initConfigFile()
		}
		return nil, err
	}

	// If we don't have any keys (for example, when the Configuration file has been created),
	// we return to the caller the default Configuration
	if len(v.AllKeys()) == 0 {
		return defaultConfig(), nil
	}

	// Unmarshal the viper Configuration into our final Config
	Config := &Config{}
	if err := v.Unmarshal(Config); err != nil {
		return nil, err
	}

	Config.location = path

	return Config, nil
}

// initConfig file will create the default Configuration file on disk
func initConfigFile() (*Config, error) {
	fileLocation := sanitizeDefaultConfigLocation()

	// Create the parent directories if they don't exists yet
	if err := os.MkdirAll(path.Dir(fileLocation), 0755); err != nil {
		return nil, err
	}

	if _, err := os.Create(fileLocation); err != nil {
		return nil, err
	}

	return Load()
}

func sanitizeDefaultConfigLocation() string {
	user, _ := user.Current()
	return strings.Replace(mortyConfigDefaultLocation, "$HOME", user.HomeDir, -1)
}

// sanitizeUrl will remove the trailing slash from the url if present.
func sanitizeUrl(url string) string {
	if strings.HasSuffix(url, "/") {
		return strings.TrimSuffix(url, "/")
	}
	return url
}
