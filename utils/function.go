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
    name  string
    runtime string
}

func (f *function) setName(name string) {
    f.name = name
}

func (f *function) getName() string {
    return f.name
}

func (f *function) getWorkingDir() string {
    return  filepath.Join("./workspaces", f.name)
}