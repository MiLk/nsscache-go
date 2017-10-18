package s3

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
)

type mockS3Client struct {
	s3iface.S3API
	getObjectResp *s3.GetObjectOutput
}

func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return m.getObjectResp, nil
}

func mockGetObjectResponse(s string) *s3.GetObjectOutput {
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(s))),
	}
}

func TestCreateS3Source(t *testing.T) {
	svc := &mockS3Client{
		getObjectResp: &s3.GetObjectOutput{},
	}

	src := CreateS3Source(svc, "prefix", "bucket")
	assert.NotNil(t, src)
}

func TestDownloadS3Data(t *testing.T) {
	svc := &mockS3Client{
		getObjectResp: mockGetObjectResponse("This is a test"),
	}

	res, err := DownloadS3Data(svc, "bucket", "key")
	assert.Equal(t, []byte("This is a test"), res)
	assert.Nil(t, err)
}
