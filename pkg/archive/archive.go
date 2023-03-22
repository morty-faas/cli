package archive

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Create a `dest` archive from `src` directory
func Zip(src, dest string) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk the directory tree
	return filepath.Walk(src, func(filePath string, fileInfo os.FileInfo, err error) error {
		log.Tracef("Adding '%s' to '%s'\n", filePath, dest)
		if err != nil {
			return err
		}

		// Ignore directories
		if fileInfo.IsDir() {
			return nil
		}

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Get the file info
		info, err := file.Stat()
		if err != nil {
			return err
		}

		// Create a zip file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = fileInfo.Name()

		// Add the file to the archive
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})
}
