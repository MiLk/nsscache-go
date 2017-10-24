package s3

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
)

type MockS3GetObject struct {
	s3iface.S3API
	resp *s3.GetObjectOutput
	err  error
}

func (m *MockS3GetObject) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return m.resp, m.err
}

/*
CreateMockS3GetObjectClient returns an object for using in your tests that will return
always the given resp and err from parameters. For example:

svc := CreateMockS3GetObjectClient("My response", nil)

Now you can use this client anywhere you need a s3iface.S3API
*/
func CreateMockS3GetObjectClient(resp string, err error) *MockS3GetObject {
	return &MockS3GetObject{
		resp: &s3.GetObjectOutput{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(resp))),
		},
		err: err,
	}
}

func TestCreateS3Source(t *testing.T) {
	svc := CreateMockS3GetObjectClient("", nil)
	src := CreateS3Source(svc, "prefix", "bucket")
	assert.NotNil(t, src)
}

func TestDownloadS3Data(t *testing.T) {
	svc := CreateMockS3GetObjectClient("This is a test", nil)
	res, err := DownloadS3Data(svc, "bucket", "key")
	assert.Equal(t, []byte("This is a test"), res)
	assert.Nil(t, err)
}
