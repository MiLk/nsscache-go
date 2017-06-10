package nsscache

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type testWriterTo struct{}

func (wt *testWriterTo) WriteTo(w io.Writer) (int64, error) {
	return 0, errors.New("test error")
}

func TestWriteAtomic(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	assert.NotNil(t, WriteAtomic(path.Join(dir, "test"), &testWriterTo{}, 06400))
}
