package utils

import (
	"os"
	"path/filepath"
)

type iFunction interface {
	setName(name string)
	getName() string
	getWorkingDir() string
	init()
}

type function struct {
    Name  string `json:"name"`
    Runtime string	`json:"runtime"`
		requiredFiles []string
}

func (f *function) setName(name string) {
    f.Name = name
}

func (f *function) getName() string {
    return f.Name
}

func (f *function) getWorkingDir() string {
	// Please, use this function to get the path to the function's working directory.
	// This function will be change in the future in order to get the current working directory.
	if _, err := os.Stat(".morty/config.json"); err == nil {
			return  "./"
	}

	if _, err := os.Stat(filepath.Join("./", f.Name, ".morty/config.json")); err == nil {
			return  filepath.Join("./", f.Name)
	}

	return  filepath.Join("./", f.Name)
}

func (f *function) init() {
	fetchTemplateFiles(*f)
	create_config_file(*f)
}