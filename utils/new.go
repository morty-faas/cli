package utils

import (
	"fmt"
	"os"
	"log"
    "net/http"
    "io"
)

type Node19Function struct {
   function
}

func (f *Node19Function) init() {
    for _, file := range []string{"handler.js", "package.json"} {
        response, err := http.Get(RUNTIME_TEMPLATES_ENDPOINTS + f.function.runtime + "/function/" + file)
        if err != nil {
            log.Fatal("cannot fetch necessary content: ", err)
        }
        defer response.Body.Close()
        template, err := io.ReadAll(response.Body)
        writeFile(f.function.getWorkingDir() + "/"+file, string(template))
    }
}

func newNode19() iFunction {
    return &Node19Function{
        function: function{
            name:  "default-node-function",
            runtime: string(Node19),
        },
    }
}

type PythonFunction struct {
    function
}

func (f *PythonFunction) init() {
    for _, file := range []string{"handler.py", "requirements.txt"} {
        response, err := http.Get(RUNTIME_TEMPLATES_ENDPOINTS + f.function.runtime + "/function/" + file)
        if err != nil {
            log.Fatal("cannot fetch necessary content: ", err)
        }
        defer response.Body.Close()
        template, err := io.ReadAll(response.Body)
        writeFile(f.function.getWorkingDir() + "/"+file, string(template))
    }
}

func newPython() iFunction {
    return &PythonFunction{
        function: function{
            name:  "default-python-function",
            runtime: string(Python3),
        },
    }
}

func getFunction(runtime string) (iFunction, error) {
    Runtime(runtime).CheckValidityOrExit()

    if runtime == string(Node19) {
        return newNode19(), nil
    }
    if runtime == string(Python3)  {
        return newPython(), nil
    }
    return nil, fmt.Errorf("ERROR: This runtime isn't supported yet.")
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

	fmt.Println("Function ",function.getName()," initialised!")
    fmt.Println("You can now start developing your function by editing the files in the working directory. (",workingDir,")")
}

func writeFile (filename string, content string) {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Fatal("ERROR: Cannot create file(s): ", err)
	}
}