package file

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"nft/config"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	storage "nft/infra/storage/model"
	file "nft/internal/file/model"
	"os"
	"strings"
	"time"

	"go.uber.org/fx"
)

type FileRepository struct {
	storage contract.IStorage
}

type FileRepositoryParams struct {
	fx.In
	Storage contract.IStorage
}

func NewFileRepository(params FileRepositoryParams) contract.IFileRepository {
	return &FileRepository{
		storage: params.Storage,
	}
}

func (f FileRepository) AddTemp(c context.Context, imageFile file.Image) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "FileRepository[AddTemp]")
	defer span.Finish()

	splittedName := strings.Split(imageFile.FileName, ".")
	if len(splittedName) != 2 {
		return "", apperrors.ErrInvalidFileExtension
	}

	tempFile, err := os.CreateTemp(config.C().File.TempDir, uuid.NewString()+"-*."+splittedName[len(splittedName)-1])
	if err != nil {
		return "", fmt.Errorf("error while creating temp image file: %w", err)
	}

	_, err = tempFile.Write(imageFile.Content)
	if err != nil {
		return "", fmt.Errorf("error while writing temp image: %w", err)
	}
	return tempFile.Name(), nil
}

func (f FileRepository) Get(c context.Context, filePath string) (*os.File, error) {
	span, c := jtrace.T().SpanFromContext(c, "FileRepository[Get]")
	defer span.Finish()
	return os.Open(filePath)
}

func (f FileRepository) Upload(c context.Context, bucket string, reader io.Reader, name string) (file.Image, error) {
	span, c := jtrace.T().SpanFromContext(c, "FileRepository[Upload]")
	defer span.Finish()

	fileUrl, err := f.storage.Add(c, storage.File{Bucket: bucket, Name: name, Content: reader})
	if err != nil {
		return file.Image{}, err
	}

	return file.Image{
		FileName: name,
		FileUrl:  fileUrl,
	}, nil
}

func (f FileRepository) GetUrl(c context.Context, bucket string, name string) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "FileRepository[GetUrl]")
	defer span.Finish()

	fileUrl, err := f.storage.GetUrl(
		c,
		storage.File{Bucket: bucket, Name: name},
		time.Minute*time.Duration(config.C().Storage.UrlExpInMin),
	)
	if err != nil {
		return "", err
	}

	return fileUrl, nil
}
