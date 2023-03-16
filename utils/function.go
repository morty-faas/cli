package utils

import (
	"log"
	"os"
)

type iFunction interface {
	setName(name string)
	getName() string
	getWorkingDir() string
	setWorkingDir(path string)
	init()
}

type function struct {
    Name  string `json:"name"`
    Runtime string	`json:"runtime"`
		requiredFiles []string
		workingDir string
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
	return f.workingDir
}

func (f *function) setWorkingDir(path string) {
	f.workingDir = path
}

func (f *function) init() {
	err := fetchTemplateFiles(*f)
	if err != nil {
		os.RemoveAll(f.getWorkingDir())
		log.Fatal(err)
	}
	err = create_config_file(*f)
	if err != nil {
		log.Fatal(err)
	}
}