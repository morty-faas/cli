package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Node19Function struct {
	function
}

func (f *Node19Function) init() {
    fetchTemplateFiles(f.function, []string{"handler.js", "package.json"})
    create_config_file(f.function)
}

func newNode19() iFunction {
    return &Node19Function{
        function: function{
            Name:  "default-node-function",
            Runtime: string(Node19),
        },
    }
}

type PythonFunction struct {
	function
}

func (f *PythonFunction) init() {
    fetchTemplateFiles(f.function, []string{"handler.py", "requirements.txt"})
    create_config_file(f.function)
}

func newPython() iFunction {
    return &PythonFunction{
        function: function{
            Name:  "default-python-function",
            Runtime: string(Python3),
        },
    }
}

type GoFunction struct {
	function
}

func (f *GoFunction) init() {
	fetchTemplateFiles(f.function, []string{"handler.go", "go.mod"})
}

func newGo() iFunction {
    return &GoFunction{
        function: function{
            Name:  "default-go-function",
            Runtime: string(Go119),
        },
    }
}

func getFunction(runtime string) (iFunction, error) {
	Runtime(runtime).CheckValidityOrExit()

	switch runtime {
	case string(Node19):
		return newNode19(), nil
	case string(Python3):
		return newPython(), nil
	case string(Go119):
		return newGo(), nil
	default:
		return nil, fmt.Errorf("ERROR: This runtime isn't supported yet.")
	}
}

func New(name string, runtime string) {
	function, err := getFunction(runtime)
	if err != nil {
		log.Fatal("ERROR: internal factory error.", err)
		return
	}

	function.setName(name)
	workingDir := function.getWorkingDir()
	if _, err := os.Stat(workingDir); !os.IsNotExist(err) {
		log.Fatal("ERROR: Function already exists. Please consider using a different name.")
	}
	err = os.MkdirAll(workingDir, 0755)
	if err != nil {
		log.Fatal("ERROR: Cannot create function directory: ", err)
	}
	function.init()

	fmt.Println("Function " + function.getName() + " initialized!")
	fmt.Println("You can now start to develop your function by editing the files in the working directory (" + workingDir + ").")
}

func writeFile(filename string, content string) {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Fatal("ERROR: Cannot create file(s): ", err)
	}
}

func fetchTemplateFiles(f function,files []string) {
    revert := false
    for _, file := range files {
        response, err := http.Get(RUNTIME_TEMPLATES_ENDPOINTS + f.Runtime + "/function/" + file)
        if err != nil || response.StatusCode != 200 {
            revert = true
        }
        defer response.Body.Close()
        template, err := io.ReadAll(response.Body)
        writeFile(f.getWorkingDir() + "/"+file, string(template))
    }

    if revert {
        os.RemoveAll(f.getWorkingDir())
        log.Fatal("ERROR: Cannot fetch template files. Please check your internet connection.\n\nReverting changes...")
    }
}

func create_config_file(f function) {
    config, err := json.Marshal(f)
	if err != nil {
		log.Fatal("ERROR: An error occurred while creating the content of the config file: ", err)
	}
    err = os.MkdirAll(f.getWorkingDir() + "/.morty", 0755)
	if err != nil {
		log.Fatal("ERROR: Cannot create config file: ", err)
	}
    writeFile(f.getWorkingDir() + "/.morty/config.json", string(config))
}
