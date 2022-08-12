package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"nft/client/jtrace"
	model "nft/client/storage/model"
	"nft/config"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	storage *minio.Client
}

func (m *Minio) Init(c context.Context) error {
	span, _ := jtrace.T().SpanFromContext(c, "Minio[Init]")
	defer span.Finish()

	minioClient, err := minio.New(
		config.C().Storage.Url,
		&minio.Options{
			Creds:  credentials.NewStaticV4(config.C().Storage.Username, config.C().Storage.Password, ""),
			Secure: config.C().Storage.SSL,
		})
	if err != nil {
		return fmt.Errorf("error happened while initializing the connection to minio storage: %w", err)
	}

	m.storage = minioClient

	return nil
}

func (m *Minio) Add(c context.Context, file model.File) (string, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "Minio[Add]")
	defer span.Finish()

	fileSize, err := file.Content.(*os.File).Seek(0, io.SeekEnd)
	if err != nil {
		return "", fmt.Errorf("error occurred while getting file size: %w", err)
	}

	objectInfo, err := m.storage.PutObject(
		ctx,
		file.Bucket,
		file.Name,
		file.Content,
		fileSize,
		minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("error occurred while uploading file: %w", err)
	}

	return objectInfo.Key, nil
}

func (m *Minio) GetUrl(c context.Context, file model.File, exp time.Duration) (string, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "Minio[GetUrl]")
	defer span.Finish()

	presignedURL, err := m.storage.PresignedGetObject(ctx, file.Bucket, file.Name, exp, url.Values{})
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("error occurred while getting file url: %w", err)
	}
	log.Println(presignedURL)
	return presignedURL.String(), nil
}
