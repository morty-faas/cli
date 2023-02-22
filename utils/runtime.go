package utils
import (
	"log"
)

const RUNTIME_TEMPLATES_ENDPOINTS = "https://raw.githubusercontent.com/polyxia-org/morty-runtimes/main/template/"

type Runtime string
const (
	Node19 Runtime = "node-19"
	Python3 Runtime = "python-3"
)

func (runtime Runtime) CheckValidityOrExit() {
    switch runtime {
    case Node19, Python3:
        return
    }
		log.Fatal("Bad runtime provided. Please use one of: \"node-19\", \"python-3\".")
}
