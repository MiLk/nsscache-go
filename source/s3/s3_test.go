package s3

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/MiLk/nsscache-go/cache"
	"github.com/stretchr/testify/assert"
)

func TestS3Source_FillPasswdCache(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := `{
  "name": "foo",
  "passwd": "x",
  "uid": 1000,
  "gid": 1000,
  "gecos": "Mr Foo",
  "dir": "/home/foo",
  "shell": "/bin/bash"
}`
	svc := &mockS3Client{
		getObjectResp: mockGetObjectResponse(r),
	}
	prefix := fmt.Sprintf("secret/%s", "nsscache-test")
	src := CreateS3Source(svc, prefix, "testing-bucket")
	c := cache.NewCache()

	assert.Nil(t, src.FillPasswdCache(c))

	var b bytes.Buffer
	n, err := c.WriteTo(&b)

	assert.Nil(t, err)
	assert.EqualValues(t, 43, n)
	assert.Equal(t, "foo:x:1000:1000:Mr Foo:/home/foo:/bin/bash\n", b.String())
}

func TestS3Source_FillShadowCache(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := `{
  "name": "foo",
  "passwd": "123!!",
  "lstchg": "1000",
  "max": "23",
  "expire": "44"
}`
	svc := &mockS3Client{
		getObjectResp: mockGetObjectResponse(r),
	}
	prefix := fmt.Sprintf("secret/%s", "nsscache-test")
	src := CreateS3Source(svc, prefix, "testing-bucket")
	c := cache.NewCache()

	assert.Nil(t, src.FillShadowCache(c))

	var b bytes.Buffer
	n, err := c.WriteTo(&b)

	assert.Nil(t, err)
	assert.EqualValues(t, 25, n)
	assert.Equal(t, "foo:123!!:1000::23:::44:\n", b.String())
}

func TestS3Source_FillGroupCache(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := `{
  "name": "group",
  "passwd": "123!!",
  "gid": 1000,
  "mem": ["foo", "var", "baz"]
}`

	svc := &mockS3Client{
		getObjectResp: mockGetObjectResponse(r),
	}
	prefix := fmt.Sprintf("secret/%s", "nsscache-test")
	src := CreateS3Source(svc, prefix, "testing-bucket")
	c := cache.NewCache()

	assert.Nil(t, src.FillGroupCache(c))

	var b bytes.Buffer
	n, err := c.WriteTo(&b)

	assert.Nil(t, err)
	assert.EqualValues(t, 29, n)
	assert.Equal(t, "group:123!!:1000:foo,var,baz\n", b.String())
}
