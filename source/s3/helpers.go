package s3

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

func CreateS3Client() (s3iface.S3API, error) {
	s := session.Must(session.NewSession())
	return s3.New(s), nil
}

func CreateS3ClientWithConfig(config *aws.Config) (s3iface.S3API, error) {
	s, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return s3.New(s), nil
}

func DownloadS3Data(c s3iface.S3API, bucket string, key string) ([]byte, error) {
	results, err := c.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}
	defer results.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, results.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
