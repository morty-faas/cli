package utils
import (
	"log"
)

const RUNTIME_TEMPLATES_ENDPOINTS = "https://raw.githubusercontent.com/polyxia-org/morty-runtimes/main/template/"

type Runtime string
const (
	Node19 Runtime = "node-19"
	Python3 Runtime = "python-3"
	Go119	 Runtime = "go-1.19"
)

func (runtime Runtime) CheckValidityOrExit() {
    switch runtime {
    case Node19, Python3, Go119:
        return
    }
		log.Fatal("Bad runtime provided. Please use one of: \"node-19\",  \"go-1.19\",  \"python-3\".")
}
