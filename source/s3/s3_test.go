package s3

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/MiLk/nsscache-go/cache"
	"github.com/stretchr/testify/assert"
)

func TestS3Source_FillPasswdCache_OK(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := `[{
	"name": "foo",
	"passwd": "x",
	"uid": 1000,
	"gid": 1000,
	"gecos": "Mr Foo",
	"dir": "/home/foo",
	"shell": "/bin/bash"
},
{
	"name": "var",
	"passwd": "x",
	"uid": 1001,
	"gid": 1000,
	"gecos": "Mr Var",
	"dir": "/home/var",
	"shell": "/bin/bash"
}]`
	svc := CreateMockS3GetObjectClient(r, nil)
	prefix := fmt.Sprintf("secret/%s", "nsscache-test")
	src := CreateS3Source(svc, prefix, "testing-bucket")
	c := cache.NewCache()

	assert.Nil(t, src.FillPasswdCache(c))

	var b bytes.Buffer
	n, err := c.WriteTo(&b)

	assert.Nil(t, err)
	assert.EqualValues(t, 86, n)
	expected := `foo:x:1000:1000:Mr Foo:/home/foo:/bin/bash
var:x:1001:1000:Mr Var:/home/var:/bin/bash
`
	assert.Equal(t, expected, b.String())
}

func TestS3Source_FillPasswdCache_JSONDecodingError(t *testing.T) {
	// json decoding error
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := `[{
	"name": "foo",
	"passwd": "x,
	"uid": 1000,
	"gid": 1000,
	"gecos": "Mr Foo",
	"dir": "/home/foo",
	"shell": "/bin/bash"
}]`
	svc := CreateMockS3GetObjectClient(r, nil)
	prefix := fmt.Sprintf("secret/%s", "nsscache-test")
	src := CreateS3Source(svc, prefix, "testing-bucket")
	c := cache.NewCache()

	err = src.FillPasswdCache(c)
	expectedErr := "json decoding: invalid character '\\n' in string literal"
	assert.Equal(t, expectedErr, err.Error())
}

func TestS3Source_FillPasswdCache_BadEntryFormat(t *testing.T) {
	// json does not match entry format
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := `[{
	"name": "Foo",
	"passwd": "x",
	"uid": "1000",
	"gid": 1000,
	"gecos": "Mr Foo",
	"dir": "/home/foo",
	"shell": "/bin/bash"
}]`
	svc := CreateMockS3GetObjectClient(r, nil)
	prefix := fmt.Sprintf("secret/%s", "nsscache-test")
	src := CreateS3Source(svc, prefix, "testing-bucket")
	c := cache.NewCache()

	err = src.FillPasswdCache(c)
	expectedErr := "json does not match entry format: json: cannot unmarshal string into Go struct field PasswdEntry.uid of type uint32"
	assert.Equal(t, expectedErr, err.Error())
}

func TestS3Source_FillPasswdCache_DownloadError(t *testing.T) {
	// error downloading from s3
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := `[{
	"name": "Foo",
	"passwd": "x",
	"uid": "1000",
	"gid": 1000,
	"gecos": "Mr Foo",
	"dir": "/home/foo",
	"shell": "/bin/bash"
}]`
	svc := CreateMockS3GetObjectClient(r, errors.New("some error"))
	prefix := fmt.Sprintf("secret/%s", "nsscache-test")
	src := CreateS3Source(svc, prefix, "testing-bucket")
	c := cache.NewCache()

	err = src.FillPasswdCache(c)
	expectedErr := "downloading from S3: error getting object secret/nsscache-test/passwd from bucket testing-bucket: some error"
	assert.Equal(t, expectedErr, err.Error())
}

func TestS3Source_FillShadowCache_OK(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := `[{
  "name": "foo",
  "passwd": "123!!",
  "lstchg": "1000",
  "max": "23",
  "expire": "44"
}]`

	svc := CreateMockS3GetObjectClient(r, nil)
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

func TestS3Source_FillGroupCache_OK(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := `[{
  "name": "group",
  "passwd": "123!!",
  "gid": 1000,
  "mem": ["foo", "var", "baz"]
}]`

	svc := CreateMockS3GetObjectClient(r, nil)
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
