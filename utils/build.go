package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	runCommand(fmt.Sprintf("docker export $(docker create %s) -o %s.tar", name, name))
}

func createExt4(name string) error {
	file, err := os.Stat(name + ".tar")
	if err != nil {
		return err
	}
	size := int64(float64(file.Size()) / 1000 / 1000 * 1.25) // 25% size increase for ext4
	runCommand(fmt.Sprintf("dd if=/dev/zero of=./%s.ext4 bs=1M count=%d > /dev/null", name, size))
	runCommand(fmt.Sprintf("mkfs.ext4 ./%s.ext4", name))
	err = os.Mkdir("rootfsdir", 0755)
	if err != nil {
		log.Fatal("ERROR: cannot create rootfs dir", err)
	}
	runCommand(fmt.Sprintf("sudo mount ./%s.ext4 rootfsdir", name))
	runCommand(fmt.Sprintf("sudo tar -C rootfsdir -xf ./%s.tar", name))
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

func Build(name string, runtime string, folder string, buildArgs []string) {
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

	//TODO: upload rootfs
	rootfsPath := filepath.Join(rootPath, name+".ext4")
	runCommand(fmt.Sprintf("cp %s.ext4 %s", name, rootfsPath))
	compressLZ4(rootfsPath)
	err = os.RemoveAll(buildPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Build complete! (\"", rootfsPath, "\")")

}
