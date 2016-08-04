package s3

import (
	"errors"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
	"time"
)

type S3UploadImageHelper struct {}

var (
	S3Bucket            *s3.Bucket
	errorDetectingExtension = errors.New(`unable to determine file extensions`)
	errorInvalidConfig = errors.New(`Invalid config`)
)

func New() (helper *S3UploadImageHelper){
	return &S3UploadImageHelper{}

}
func (helper *S3UploadImageHelper) SetupS3Connection(config []string) (e error){
	if len(config) < 3 {
		return errorInvalidConfig
	}
	accessKey, secret, bucket := config[0], config[1], config[2]
	auth, err := aws.GetAuth(accessKey, secret, ``, time.Time{})
	if err != nil {
		return err
	}
	conn := s3.New(auth, aws.USEast) // our S3 is only accessible from USEast (east-1) region
	S3Bucket = conn.Bucket(bucket)
	return nil
}

func (helper *S3UploadImageHelper) UploadImage(file []byte, uploadPath string) (error) {
	return uploadImage(S3Bucket, file, uploadPath)
}

// take the base64-encoded string of profile image, upload to S3 and get the URL
var uploadImage = func(bucket *s3.Bucket, file []byte, uploadPath string) (e error) {
	if err := bucket.Put(uploadPath, file, `text/plain`, s3.PublicRead, s3.Options{}); err != nil {
		return err
	}
	return nil
}
