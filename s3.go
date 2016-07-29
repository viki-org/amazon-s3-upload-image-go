package s3

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
	"github.com/viki-org/nd"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	S3Bucket            *s3.Bucket
	originalUploadImage = uploadImage
	errorDetectingExtension = errors.New(`unable to determine file extensions`)
)

func SetupS3Connection(accessKey string, secret string, bucket string) {
	auth, err := aws.GetAuth(accessKey, secret, ``, time.Time{})
	if err != nil {
		log.Fatal(err)
	}
	conn := s3.New(auth, aws.USEast) // our S3 is only accessible from USEast (east-1) region
	S3Bucket = conn.Bucket(bucket)
}

func UploadImage(data, userID string) (string, error) {
	return uploadImage(S3Bucket, data, userID)
}

// take the base64-encoded string of profile image, upload to S3 and get the URL
var uploadImage = func(bucket *s3.Bucket, data, userID string) (string, error) {
	// return url looks like http://0.viki.io/u/somerandomstring1234.png, http://1.viki.io/123asdf.jpg, etc
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatal(err)
		return ``, err
	}
	extension, err := getExtension(decoded)
	if err != nil {
		log.Fatal(err)
		return ``, err
	}
	timestamp := strconv.FormatInt(nd.Now().Unix(), 10)
	filename := hexdigest(append(decoded, timestamp...))[0:10] + `.` + extension
	uploadPath := "uploads/user/profile_image/" + userID + `/` + filename
	subdomain := strconv.Itoa(nd.IntnRand(2))

	filetype := http.DetectContentType([]byte(decoded))
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif":
	case "image/png":
	default:
		return ``, errors.New(`Not valid image`)
	}

	if err := bucket.Put(uploadPath, []byte(decoded), `text/plain`, s3.PublicRead, s3.Options{}); err != nil {
		log.Fatal(err)
		return ``, err
	}
	url := `http://` + subdomain + `.viki.io/u/` + userID + `/` + filename
	return url, nil
}

// find the extension based on (header?) bytes of image files
var getExtension = func(data []byte) (string, error) {
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
		return `gif`, nil
	} // header is GIF8
	if data[0] == 0xFF && data[1] == 0xD8 {
		return `jpg`, nil
	} // JPEG header
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return `png`, nil
	} // PNG header

	return ``, errorDetectingExtension
}

// create a hex checksum of the inpuut data
var hexdigest = func(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func StubUploadImage(url string) {
	uploadImage = func(bucket *s3.Bucket, data, userID string) (string, error) {
		uploadImage = originalUploadImage
		return url, nil
	}
}
