package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func newNode19() iFunction {
	return &function{
		Name:          "default-node-function",
		Runtime:       string(Node19),
		requiredFiles: []string{"handler.js", "package.json"},
	}
}

func newPython() iFunction {
	return &function{
		Name:          "default-python-function",
		Runtime:       string(Python3),
		requiredFiles: []string{"handler.py", "requirements.txt"},
	}
}

func newGo() iFunction {
	return &function{
		Name:          "default-go-function",
		Runtime:       string(Go119),
		workingDir:    "default-go-function",
		requiredFiles: []string{"handler.go", "go.mod"},
	}
}

func newRust() iFunction {
	return &function{
		Name:  "default-rust-function",
		Runtime: string(Rust167),
		requiredFiles: []string{"handler.rs", "Cargo.toml"},
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
	case string(Rust167):
		return newRust(), nil
	default:
		return nil, fmt.Errorf("ERROR: This runtime isn't supported yet.")
	}
}

func New(name string, path string, runtime string) {
	function, err := getFunction(runtime)
	function.setWorkingDir(path)
	if err != nil {
		log.Fatal("ERROR: internal factory error.", err)
		return
	}

	function.setName(name)
	if _, err := os.Stat(function.getWorkingDir()); !os.IsNotExist(err) {
		log.Fatal("ERROR: Function already exists. Please consider using a different name.")
	}
	err = os.MkdirAll(function.getWorkingDir(), 0755)
	if err != nil {
		log.Fatal("ERROR: Cannot create function directory: ", err)
	}
	function.init()

	fmt.Println("Function " + function.getName() + " initialized!")
	fmt.Println("You can now start to develop your function by editing the files in the working directory (" + function.getWorkingDir() + ").")
}

func writeFile(filename string, content string) {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Fatal("ERROR: Cannot create file(s): ", err)
	}
}

func fetchTemplateFiles(f function) error {
	for _, file := range f.requiredFiles {
		response, err := http.Get(RUNTIME_TEMPLATES_ENDPOINTS + f.Runtime + "/function/" + file)
		if err != nil || response.StatusCode != 200 {
			return errors.New("ERROR: Cannot fetch template files. Please check your internet connection.\n\nReverting changes...")
		}
		defer response.Body.Close()
		template, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("ERROR: An error occurred while reading the template file: %q", err)
		}
		writeFile(f.getWorkingDir()+"/"+file, string(template))
	}
	return nil
}

func create_config_file(f function) error {
	config, err := json.Marshal(f)
	if err != nil {
		return fmt.Errorf("ERROR: An error occurred while creating the content of the config file: %q", err)
	}
	err = os.MkdirAll(f.getWorkingDir()+"/.morty", 0755)
	if err != nil {
		return fmt.Errorf("ERROR: Cannot create config file: %q", err)
	}
	writeFile(f.getWorkingDir()+"/.morty/config.json", string(config))
	return nil
}
