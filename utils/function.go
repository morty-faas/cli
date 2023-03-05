package utils

import (
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
    return  filepath.Join("./workspaces", f.Name)
}