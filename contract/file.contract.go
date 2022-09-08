package contract

import (
	"context"
	"io"
	file "nft/src/file/model"
	"os"
)

type IFileService interface {
	// NewImageFile(c context.Context, f file.Image) (string, error)
	UploadKycImage(c context.Context, imageFile file.Image) (string, error)
	UploadNftImage(c context.Context, imageFile file.Image) (string, error)
	GetKycImageUrl(c context.Context, name string) (string, error)
	GetNftImageUrl(c context.Context, name string) (string, error)
}

type IFileRepository interface {
	// Add(c context.Context, file []byte, name string) (string, error)
	AddTemp(c context.Context, imageFile file.Image) (string, error)
	Get(c context.Context, filePath string) (*os.File, error)
	Upload(c context.Context, bucket string, file io.Reader, name string) (file.Image, error)
	GetUrl(c context.Context, bucket string, name string) (string, error)
}
