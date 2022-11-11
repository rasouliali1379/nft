package contract

import (
	"context"
	"io"
	file "nft/internal/file/model"
	"os"
)

type IFileService interface {
	UploadImage(c context.Context, imageFile file.Image) (string, error)
	GetImageUrl(c context.Context, imageFile file.Image) (string, error)
}

type IFileRepository interface {
	AddTemp(c context.Context, imageFile file.Image) (string, error)
	Get(c context.Context, filePath string) (*os.File, error)
	Upload(c context.Context, bucket string, file io.Reader, name string) (file.Image, error)
	GetUrl(c context.Context, bucket string, name string) (string, error)
}
