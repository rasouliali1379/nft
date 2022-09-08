package file

import (
	"context"
	"nft/client/jtrace"
	"nft/config"
	"nft/contract"
	file "nft/src/file/model"
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

func (f FileService) UploadNftImage(c context.Context, imageFile file.Image) (string, error) {
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

	uploaded, err := f.fileRepository.Upload(c, config.C().Storage.Buckets.NFT, reader, splittedPath[len(splittedPath)-1])
	if err != nil {
		return "", err
	}

	return uploaded.FileName, nil
}

func (f FileService) UploadKycImage(c context.Context, imageFile file.Image) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "FileService[UploadKycImage]")
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

	uploaded, err := f.fileRepository.Upload(c, config.C().Storage.Buckets.KYC, reader, splittedPath[len(splittedPath)-1])
	if err != nil {
		return "", err
	}

	return uploaded.FileName, nil
}

func (f FileService) GetNftImageUrl(c context.Context, name string) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "FileService[GetNftImageUrl]")
	defer span.Finish()
	return f.fileRepository.GetUrl(c, config.C().Storage.Buckets.NFT, name)
}

func (f FileService) GetKycImageUrl(c context.Context, name string) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "FileService[GetKycImageUrl]")
	defer span.Finish()
	return f.fileRepository.GetUrl(c, config.C().Storage.Buckets.KYC, name)
}
