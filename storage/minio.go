package storage

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
	"os"
	"time"
)

type Storage struct {
	client *minio.Client
}

func New() *Storage {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_USER")
	secretAccessKey := os.Getenv("MINIO_PASSWORD")
	useSSL := false

	// Initialize minio client object.
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &Storage{mc}
}

func (s *Storage) StoreImage(objectName string, filePath string) {
	ctx := context.Background()

	bucketName := "images"

	// get presigned url from minio
	presignedURL, err := s.client.PresignedPutObject(ctx, bucketName, objectName, time.Duration(60)*time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	// add file data in buffer
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Unable to open %s: %v", filePath, err)
	}

	fileBuffer := bytes.NewBuffer(file)

	// create http request for file upload
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, presignedURL.String(), fileBuffer)

	if err != nil {
		log.Fatalf("Unable to create PUT request: %v", err)
	}

	// upload file to minio
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Unable to PUT: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Unable to upload %s to %s: %v", filePath, bucketName, res.Status)
	}

	log.Printf("Successfully uploaded %s\n", objectName)
}
