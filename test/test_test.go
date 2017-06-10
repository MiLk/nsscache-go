package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/milk/nsscache-go/cache"
)

func createPasswdCache(dir string) *cache.Cache {
	passwd := cache.NewCache("passwd", cache.Dir(dir))
	passwd.Add(&cache.PasswdEntry{
		Name:   "foo",
		Passwd: "x",
		UID:    1000,
		GID:    1000,
		GECOS:  "Mr Foo",
		Dir:    "/home/foo",
		Shell:  "/bin/bash",
	})
	return passwd
}

func createShadowCache(dir string) *cache.Cache {
	shadow := cache.NewCache("shadow", cache.Dir(dir))
	lstchg := int32(time.Now().Sub(time.Unix(0, 0)).Hours() / 24)
	shadow.Add(&cache.ShadowEntry{
		Name:   "foo",
		Passwd: "!!",
		Lstchg: cache.Int32(lstchg),
	})
	return shadow
}

func createGroupCache(dir string) *cache.Cache {
	shadow := cache.NewCache("group", cache.Dir(dir))
	shadow.Add(&cache.GroupEntry{
		Name:   "foo",
		Passwd: "*",
		GID:    1000,
	})
	return shadow
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

func TestLibNsscache(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	passwd := createPasswdCache(dir)
	shadow := createShadowCache(dir)
	group := createGroupCache(dir)

	assert.Nil(t, passwd.WriteFile())
	assert.Nil(t, shadow.WriteFile())
	assert.Nil(t, group.WriteFile())

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
}
