package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"nft/client/jtrace"
	model "nft/client/storage/model"
	"nft/config"
	"time"
)

type Aws struct {
	sess *session.Session
}

func (a *Aws) Init(c context.Context) error {
	span, _ := jtrace.T().SpanFromContext(c, "Aws[Init]")
	defer span.Finish()

	var err error

	a.sess, err = session.NewSession(
		&aws.Config{
			Endpoint:         aws.String(config.C().Storage.Url),
			Region:           aws.String(endpoints.UsWest2RegionID),
			S3ForcePathStyle: aws.Bool(true),
			Credentials: credentials.NewStaticCredentials(
				config.C().Storage.Username,
				config.C().Storage.Password,
				"",
			),
		})

	return err
}

func (a *Aws) Add(c context.Context, file model.File) (string, error) {
	span, _ := jtrace.T().SpanFromContext(c, "Aws[Add]")
	defer span.Finish()

	uploader := s3manager.NewUploader(a.sess)
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(file.Bucket),
		//ACL:    aws.String("public-read"),
		Key:  aws.String(file.Name),
		Body: file.Content,
	})

	if err != nil {
		return "", err
	}

	return up.Location, nil
}

func (a *Aws) GetUrl(c context.Context, file model.File, exp time.Duration) (string, error) {
	span, _ := jtrace.T().SpanFromContext(c, "Aws[GetUrl]")
	defer span.Finish()
	svc := s3.New(a.sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(file.Bucket),
		Key:    aws.String(file.Name),
	})

	urlStr, err := req.Presign(exp)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}
