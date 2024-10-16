package uploader

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"mime/multipart"
	"os"
	"sync"
)

// S3Uploader is a struct that wraps around AWS SDK and provides
// methods for uploading to S3 bucket.
type S3Uploader struct {
	client   *s3.Client
	bucket   string
	uploader *manager.Uploader
}

// NewS3Uploader creates a new S3Uploader with the provided bucket name.
func NewS3Uploader(bucketName string) (*S3Uploader, error) {
	awsConfig, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to load aws SDK config: %v", err)
	}

	s3Client := s3.NewFromConfig(awsConfig)
	uploader := manager.NewUploader(s3Client)

	return &S3Uploader{
		client:   s3Client,
		bucket:   bucketName,
		uploader: uploader,
	}, nil
}

// UploadImages uploads the provided image files to S3 bucket.
func (u *S3Uploader) UploadImages(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	urls := make([]string, 0, len(files))
	var errs []error
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, fileHeader := range files {
		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()

			file, err := fh.Open()
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			defer file.Close()
			filename, err := generateUniqueFilename(fh.Filename, fh.Header.Get("Content-Type"))
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			_, err = u.uploader.Upload(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(u.bucket),
				Key:         aws.String(filename),
				Body:        file,
				ContentType: aws.String(fh.Header.Get("Content-Type")),
				ACL:         types.ObjectCannedACLPublicRead,
			})

			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			url := fmt.Sprintf("https://s3.amazonaws.com/%s/%s", u.bucket, filename)

			mu.Lock()
			urls = append(urls, url)
			mu.Unlock()
		}(fileHeader)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, fmt.Errorf("multiple errors occured: %v", errs)
	}
	return urls, nil
}
