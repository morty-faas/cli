package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"log"
    "net/http"
    "io"
)

const RUNTIME_TEMPLATES_ENDPOINTS = "https://raw.githubusercontent.com/polyxia-org/runtimes/main/template/"

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


type Node19Function struct {
   function
}

func (f *Node19Function) init() {
    responseHandler, err := http.Get(RUNTIME_TEMPLATES_ENDPOINTS + f.function.runtime + "/function/handler.js")
    if err != nil {
        log.Fatal("cannot fetch necessary content: ", err)
    }
    defer responseHandler.Body.Close()
    functionTemplateHandler, err := io.ReadAll(responseHandler.Body)
	writeFile(f.function.getWorkingDir() + "/handler.js", string(functionTemplateHandler))
    
    responsePackageJson, err := http.Get(RUNTIME_TEMPLATES_ENDPOINTS + f.function.runtime + "/function/package.json")
    if err != nil {
        log.Fatal("cannot fetch necessary content: ", err)
    }
    defer responsePackageJson.Body.Close()
    functionTemplatePackageJson, err := io.ReadAll(responsePackageJson.Body)
    writeFile(f.function.getWorkingDir() + "/package.json", string(functionTemplatePackageJson))

}

func newNode19() iFunction {
    return &Node19Function{
        function: function{
            name:  "default",
            runtime: "node-19",
        },
    }
}

type PythonFunction struct {
    function
}

func (f *PythonFunction) init() {
	fmt.Println("Python function not implemented yet (", f.function.getWorkingDir(), "/main.py", ")")
}

func newPython() iFunction {
    return &PythonFunction{
        function: function{
            name:  "default",
            runtime: "node-19",
        },
    }
}

func getFunction(runtime string) (iFunction, error) {
    if runtime == "node-19" {
        return newNode19(), nil
    }
    if runtime == "python" {
        return newPython(), nil
    }
    return nil, fmt.Errorf("Wrong gun type passed")
}

func New(name string, runtime string) {
	function, err := getFunction(runtime)
	if err != nil {
			log.Fatal("factory error: ", err)
		return
	}

	function.setName(name)
	workingDir := function.getWorkingDir()
	os.MkdirAll(workingDir, 0755)
	function.init()

	fmt.Println("Function ",function.getName()," initialised! (\"", workingDir, "\")")
}

func writeFile (filename string, content string) {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Fatal("cannot create files: ", err)
	}
}