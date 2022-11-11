package file

import (
	"context"
	"nft/contract"
	"nft/infra/jtrace"
	file "nft/internal/file/model"
	"strings"

	"go.uber.org/fx"
)

type FileService struct {
	fileRepository contract.IFileRepository
}

type FileServiceParams struct {
	fx.In
	FileRepository contract.IFileRepository
}

func NewFileService(params FileServiceParams) contract.IFileService {
	return FileService{
		fileRepository: params.FileRepository,
	}
}

func (f FileService) UploadImage(c context.Context, imageFile file.Image) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "FileService[UploadNftImage]")
	defer span.Finish()

	fileName, err := f.fileRepository.AddTemp(c, imageFile)
	if err != nil {
		return "", err
	}

	reader, err := f.fileRepository.Get(c, fileName)
	if err != nil {
		return "", err
	}

	splittedPath := strings.Split(fileName, "/")

	uploaded, err := f.fileRepository.Upload(c, imageFile.Bucket, reader, splittedPath[len(splittedPath)-1])
	if err != nil {
		return "", err
	}

	return uploaded.FileName, nil
}

func (f FileService) GetImageUrl(c context.Context, imageFile file.Image) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "FileService[GetNftImageUrl]")
	defer span.Finish()
	return f.fileRepository.GetUrl(c, imageFile.Bucket, imageFile.FileName)
}
