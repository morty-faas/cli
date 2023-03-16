package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	cp "github.com/otiai10/copy"
	"github.com/pierrec/lz4"
)

func runCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func fetchTemplate(runtime string) {
	runCommand("git init -b main > /dev/null")
	runCommand("git remote add origin https://github.com/polyxia-org/runtimes.git")
	runCommand("git config core.sparseCheckout true")
	runCommand(fmt.Sprintf("echo %s > .git/info/sparse-checkout", runtime))
	runCommand("git pull origin main > /dev/null")
	runCommand(fmt.Sprintf("mv %s/* .", "template/"+runtime))
}

func copyUserFunction(src string, dest string) error {
	err := os.RemoveAll(filepath.Join(dest, "function"))
	if err != nil {
		return fmt.Errorf("ERROR: Cannot remove file(s). %s", err)
	}
	err = os.MkdirAll(filepath.Join(dest, "function"), 0755)
	if err != nil {
		return fmt.Errorf("ERROR: Cannot create folder(s)/file(s) %s", err)
	}
	return cp.Copy(src, filepath.Join(dest, "function"))
}

func getImageSize(name string) (int, error) {
	cmd := exec.Command("docker", "inspect", name, "--format", "{{.Size}}")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	sizeString := strings.TrimSpace(string(output))
	size, err := strconv.Atoi(sizeString)

	if err != nil {
		return 0, err
	}

	size = int(float64(size) / 1000 / 1000 * 1.25) // 25% size increase for ext4

	return size, nil
}

func buildImage(name string, buildArgs []string) {
	parsedBuildArgs := ""
	//loop through build args and add them to the docker build command
	for _, args := range buildArgs {
		index := strings.Index(args, "=")
		if index == -1 {
			log.Fatal("each build-arg must take the form key=value")
		}

		values := []string{args[0:index], args[index+1:]}

		k := strings.TrimSpace(values[0])
		v := strings.TrimSpace(values[1])
		parsedBuildArgs += fmt.Sprintf("--build-arg %s='%s' ", k, v)
	}

	runCommand(fmt.Sprintf("docker build %s -t %s . > /dev/null", parsedBuildArgs, name))
}

func createExt4(name string) error {

	size, err := getImageSize(name)
	if err != nil {
		return err
	}

	// Create the ext4 mount point
	runCommand(fmt.Sprintf("dd if=/dev/zero of=./%s.ext4 bs=1M count=%d > /dev/null", name, size))
	runCommand(fmt.Sprintf("mkfs.ext4 ./%s.ext4", name))

	if err = os.Mkdir("rootfsdir", 0755); err != nil {
		log.Fatal("ERROR: cannot create rootfs dir", err)
	}

	runCommand(fmt.Sprintf("sudo mount ./%s.ext4 rootfsdir", name))

	// Create a script to export the rootfs from within the Docker container
	runCommand("cat <<EOF > ./rootfs-export.sh\nfor d in app bin etc lib root sbin usr; do tar c \"/\\${d}\" | tar x -C /my-rootfs; done\nfor dir in dev proc run sys var; do mkdir /my-rootfs/\\${dir}; done\nexit\nEOF")

	// export the rootfs from the Docker container
	runCommand(fmt.Sprintf("sudo docker run --rm -i -v ./rootfsdir:/my-rootfs -v /dev/urandom:/dev/random %s:latest sh <./rootfs-export.sh", name))

	// Inject the custom /sbin/init script
	runCommand("cat <<EOF > $(pwd)/alpha_init\n#!/bin/sh\nsource /app/env.sh\n/usr/bin/alpha\nEOF")

	runCommand("sudo mv rootfsdir/sbin/init rootfsdir/sbin/init.old")
	runCommand("sudo cp $(pwd)/alpha_init rootfsdir/sbin/init")
	runCommand("sudo chmod a+x rootfsdir/sbin/init")

	runCommand("sudo umount rootfsdir")
	return nil
}

func compressLZ4(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var buff bytes.Buffer
	w := lz4.NewWriter(&buff)
	w.Write(data)
	w.Close()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fmt.Sprintf("%s.lz4", absPath), buff.Bytes(), 0644)
}

func Build(folder string, buildArgs []string) {
	name := GetFunctionValueFromConfig(folder, "name")
	runtime := GetFunctionValueFromConfig(folder, "runtime")

	Runtime(runtime).CheckValidityOrExit()

	// Check if user is root
	currentUser, e := user.Current()
	if e != nil {
		log.Fatalf("Unable to get current user: %s", e)
	}

	if currentUser.Uid != "0" {
		log.Fatal("You must be root to build a function")
	}

	rootPath, _ := os.Getwd()
	buildPath := "/tmp/morty/builds/" + uuid.New().String()
	if err := os.MkdirAll(buildPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	os.Chdir(buildPath)
	fetchTemplate(runtime)

	os.Chdir(rootPath)
	err := copyUserFunction(folder, buildPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Chdir(buildPath)
	buildImage(name, buildArgs)
	createExt4(name)

	rootfsPath := filepath.Join(rootPath, name+".ext4")
	runCommand(fmt.Sprintf("cp %s.ext4 %s", name, rootfsPath))
	compressLZ4(rootfsPath)

	//upload image to minio
	objectName := fmt.Sprintf("%s.ext4.lz4", name)
	filePath := fmt.Sprintf("%s/%s", rootPath, objectName)

	registry := NewRegistry()
	registry.UploadFile(objectName, filePath)

	// minioClient := storage.New()
	// minioClient.StoreImage(objectName, filePath)

	err = os.RemoveAll(buildPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Build successful! (\"", rootfsPath, "\")")

}

func GetFunctionValueFromConfig(folder string, parameter string) string {
	file, err := os.Open(filepath.Join(folder, ".morty/config.json"))
	if err != nil {
		log.Fatal("ERROR: cannot open config file. Does the .morty folder exist?")
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result[parameter].(string)
}
