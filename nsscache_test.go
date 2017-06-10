package nsscache

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"path"

	"github.com/milk/nsscache-go/cache"
	"github.com/pkg/errors"
)

type testSource struct{}

func (s *testSource) FillPasswdCache(c *cache.Cache) error {
	c.Add(&cache.PasswdEntry{
		Name:   "foo",
		Passwd: "x",
		UID:    1000,
		GID:    1000,
		GECOS:  "Mr Foo",
		Dir:    "/home/foo",
		Shell:  "/bin/bash",
	})
	return nil
}

func (s *testSource) FillShadowCache(c *cache.Cache) error {
	lstchg := int32(time.Now().Sub(time.Unix(0, 0)).Hours() / 24)
	c.Add(&cache.ShadowEntry{
		Name:   "foo",
		Passwd: "!!",
		Lstchg: cache.Int32(lstchg),
	})
	return nil
}

func (s *testSource) FillGroupCache(c *cache.Cache) error {
	c.Add(&cache.GroupEntry{
		Name:   "foo",
		Passwd: "*",
		GID:    1000,
	})
	return nil
}

type errorSource map[string]bool

func (s *errorSource) FillPasswdCache(c *cache.Cache) error {
	if (map[string]bool)(*s)["passwd"] {
		return errors.New("error")
	}
	return nil
}

func (s *errorSource) FillShadowCache(c *cache.Cache) error {
	if (map[string]bool)(*s)["shadow"] {
		return errors.New("error")
	}
	return nil
}

func (s *errorSource) FillGroupCache(c *cache.Cache) error {
	if (map[string]bool)(*s)["group"] {
		return errors.New("error")
	}
	return nil
}

func Getent(dir string, args ...string) ([]byte, error) {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return nil, err
	}
	cmdArgs := []string{
		"run", "--rm",
		"-v", fmt.Sprintf("%s:%s:ro", filepath.Join(dir, "passwd.cache"), "/etc/passwd.cache"),
		"-v", fmt.Sprintf("%s:%s:ro", filepath.Join(dir, "shadow.cache"), "/etc/shadow.cache"),
		"-v", fmt.Sprintf("%s:%s:ro", filepath.Join(dir, "group.cache"), "/etc/group.cache"),
		"nsscache-go",
		"getent",
	}
	cmd := exec.Command(docker, append(cmdArgs, args...)...)
	return cmd.CombinedOutput()
}

func TestCacheMap_WriteFiles(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	cm := NewCaches()
	src := testSource{}
	assert.Nil(t, cm.FillCaches(&src))
	assert.Nil(t, cm.WriteFiles(&WriteOptions{
		Directory: dir,
	}))

	var res []byte
	res, err = Getent(dir, "passwd", "foo")
	assert.Nil(t, err)
	assert.Equal(t, "foo:x:1000:1000:Mr Foo:/home/foo:/bin/bash\n", string(res))

	res, err = Getent(dir, "shadow", "foo")
	lstchg := int32(time.Now().Sub(time.Unix(0, 0)).Hours() / 24)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("foo:!!:%d::::::\n", lstchg), string(res))

	res, err = Getent(dir, "group", "foo")
	assert.Nil(t, err)
	assert.Equal(t, "foo:*:1000:\n", string(res))

	assert.Nil(t, cm.WriteFiles(&WriteOptions{
		Directory: dir,
		Extension: "cachetest",
	}))

	_, err = os.Stat(path.Join(dir, "passwd.cachetest"))
	assert.Nil(t, err)
	_, err = os.Stat(path.Join(dir, "group.cachetest"))
	assert.Nil(t, err)
	_, err = os.Stat(path.Join(dir, "shadow.cachetest"))
	assert.Nil(t, err)

	assert.NotNil(t, cm.WriteFiles(&WriteOptions{
		Directory: "/tmp/does_not_exist",
	}))
}

func TestCacheMap_FillCaches(t *testing.T) {
	cm := NewCaches()
	src := testSource{}
	delete(cm, "shadow")
	assert.Nil(t, cm.FillCaches(&src))
	delete(cm, "shadow")
	assert.Nil(t, cm.FillCaches(&src))
	delete(cm, "group")
	assert.Nil(t, cm.FillCaches(&src))
}

func TestCacheMap_FillCaches2(t *testing.T) {
	cm := NewCaches()
	src := errorSource{}
	assert.Nil(t, cm.FillCaches(&src))
	src["group"] = true
	assert.NotNil(t, cm.FillCaches(&src))
	src["shadow"] = true
	assert.NotNil(t, cm.FillCaches(&src))
	src["passwd"] = true
	assert.NotNil(t, cm.FillCaches(&src))
}
