package cache

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
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

func TestWithACL(t *testing.T) {
	c := NewCache(WithACL(func(e Entry) bool {
		pe, ok := e.(*PasswdEntry)
		if !ok {
			return true
		}
		return pe.Name == "admin"
	}))

	c.Add(&PasswdEntry{
		Name:   "foo",
		Passwd: "x",
		UID:    1000,
		GID:    1000,
		GECOS:  "Mr Foo",
		Dir:    "/home/foo",
		Shell:  "/bin/bash",
	}, &PasswdEntry{
		Name:   "bar",
		Passwd: "x",
		UID:    1001,
		GID:    1000,
		GECOS:  "Mrs Bar",
		Dir:    "/home/bar",
		Shell:  "/bin/bash",
	}, &PasswdEntry{
		Name:   "admin",
		Passwd: "x",
		UID:    1002,
		GID:    1000,
		GECOS:  "Admin",
		Dir:    "/home/admin",
		Shell:  "/bin/bash",
	})

	var b bytes.Buffer
	n, err := c.WriteTo(&b)
	assert.Nil(t, err)
	assert.EqualValues(t, 46, n)
	expected := "admin:x:1002:1000:Admin:/home/admin:/bin/bash\n"
	assert.Equal(t, expected, b.String())
}

type errorWriter struct{}

func (w *errorWriter) Write(b []byte) (int, error) {
	return 0, errors.New("error")
}

func TestCache_WriteTo(t *testing.T) {
	c := NewCache()
	c.Add(&PasswdEntry{})
	w := &errorWriter{}
	_, err := c.WriteTo(w)
	assert.NotNil(t, err)
}
