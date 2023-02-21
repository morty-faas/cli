package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
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
	contentType := "application/lz4"

	info, err := s.client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalf("Unable to upload %s to %s: %v", filePath, bucketName, err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}
