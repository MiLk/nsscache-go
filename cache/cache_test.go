package cache

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache_Add(t *testing.T) {
	c := NewCache()
	c.Add(&PasswdEntry{
		Name:   "foo",
		Passwd: "x",
		UID:    1000,
		GID:    1000,
		GECOS:  "Mr Foo",
		Dir:    "/home/foo",
		Shell:  "/bin/bash",
	})

	var b bytes.Buffer
	n, err := c.WriteTo(&b)
	assert.Nil(t, err)
	assert.EqualValues(t, 43, n)
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

	b.Reset()
	n, err = c.WriteTo(&b)
	assert.Nil(t, err)
	assert.EqualValues(t, 87, n)
	expected = "foo:x:1000:1000:Mr Foo:/home/foo:/bin/bash\nbar:x:1001:1000:Mrs Bar:/home/bar:/bin/bash\n"
	assert.Equal(t, expected, b.String())
}
