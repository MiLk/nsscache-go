package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	c := NewCache("foo", Dir("/tmp/test/foo"), Extension("bar"))
	assert.Equal(t, "/tmp/test/foo/foo.bar", c.filename())
	c = NewCache("bar", Dir("/tmp/test/foo/"), Extension("foo"))
	assert.Equal(t, "/tmp/test/foo/bar.foo", c.filename())
}

func TestCache_Add(t *testing.T) {
	c := NewCache("cache")
	c.Add(&PasswdEntry{
		Name:   "foo",
		Passwd: "x",
		UID:    1000,
		GID:    1000,
		GECOS:  "Mr Foo",
		Dir:    "/home/foo",
		Shell:  "/bin/bash",
	})

	b, err := c.buffer()
	assert.Nil(t, err)
	expected := "foo:x:1000:1000:Mr Foo:/home/foo:/bin/bash\n"
	assert.Equal(t, expected, b.String())

	c.Add(&PasswdEntry{
		Name:   "bar",
		Passwd: "x",
		UID:    1001,
		GID:    1000,
		GECOS:  "Mrs Bar",
		Dir:    "/home/bar",
		Shell:  "/bin/bash",
	})

	b, err = c.buffer()
	assert.Nil(t, err)
	expected = "foo:x:1000:1000:Mr Foo:/home/foo:/bin/bash\nbar:x:1001:1000:Mrs Bar:/home/bar:/bin/bash\n"
	assert.Equal(t, expected, b.String())
}
