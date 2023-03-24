package cliconfig

import (
	"io/fs"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_initConfigFile_CreateFile(t *testing.T) {
	_, err := initConfigFile()
	assert.NoError(t, err)
}

func Test_Load_ReturnDefaultConfigurationWhenEmptyFile(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, defaultConfig(), c)
}

func Test_Load_ReturnUserProvidedConfigFromDefaultLocation(t *testing.T) {
	location := sanitizeDefaultConfigLocation()
	expectedCfg := &Config{
		location: location,
		Contexts: []Context{
			{
				Name:     "test",
				Gateway:  "http://gateway.morty.test",
				Registry: "http://registry.morty.test",
			},
		},
	}
	b, _ := yaml.Marshal(expectedCfg)
	os.WriteFile(location, b, fs.FileMode(os.O_RDWR))

	c, err := Load()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, expectedCfg, c)
}

func Test_hasContext(t *testing.T) {
	config := defaultConfig()
	assert.True(t, config.hasContext("localhost"))
	assert.False(t, config.hasContext("production"))
}

func Test_UseContext(t *testing.T) {
	config := defaultConfig()
	config.Contexts = append(config.Contexts, Context{
		Name:     "test",
		Gateway:  "test",
		Registry: "test",
	})

	assert.True(t, config.hasContext("localhost"))
	assert.True(t, config.hasContext("test"))

	assert.NoError(t, config.UseContext("test"))
	assert.Equal(t, config.Current, "test")
}

func Test_GetContext(t *testing.T) {
	expected := Context{
		Name:     "test",
		Gateway:  "test",
		Registry: "test",
	}

	config := defaultConfig()
	config.Contexts = append(config.Contexts, expected)

	ctx, err := config.GetContext("test")

	assert.Equal(t, &expected, ctx)
	assert.NoError(t, err)
}

func Test_GetCurrentContext(t *testing.T) {
	expected := Context{
		Name:     "test",
		Gateway:  "test",
		Registry: "test",
	}

	config := defaultConfig()
	config.Contexts = append(config.Contexts, expected)
	config.UseContext("test")

	ctx, err := config.GetCurrentContext()

	assert.Equal(t, &expected, ctx)
	assert.NoError(t, err)
}

func Test_AddContext(t *testing.T) {
	expected := Context{
		Name:     "test",
		Gateway:  "test",
		Registry: "test",
	}

	config := defaultConfig()

	err := config.AddContext(&expected)
	assert.NoError(t, err)
	assert.Len(t, config.Contexts, 2)

	err = config.AddContext(&expected)
	assert.ErrorIs(t, err, ErrContextAlreadyExistsWithName)
}

func Test_SanitizeUrl(t *testing.T) {
	url := "http://localhost:8080/"
	expected := "http://localhost:8080"

	assert.Equal(t, expected, sanitizeUrl(url))

	url = "http://localhost:8080"
	assert.Equal(t, expected, sanitizeUrl(url))
}
