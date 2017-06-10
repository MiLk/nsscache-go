package cache

import (
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestPasswdEntry_String(t *testing.T) {
	e := PasswdEntry{
		Name:  "foo",
		UID:   1000,
		GID:   1000,
		GECOS: "Mr Foo",
		Dir:   "/home/foo",
		Shell: "/usr/bin/bash",
	}
	expected := "foo:x:1000:1000:Mr Foo:/home/foo:/usr/bin/bash\n"
	assert.Equal(t, expected, e.String())
}

func TestPasswdEntry_WriteTo(t *testing.T) {
	e := PasswdEntry{
		Name:  "foo",
		UID:   1000,
		GID:   1000,
		GECOS: "Mr Foo",
		Dir:   "/home/foo",
		Shell: "/usr/bin/bash",
	}
	expected := "foo:x:1000:1000:Mr Foo:/home/foo:/usr/bin/bash\n"
	var b bytes.Buffer
	assert.Nil(t, writerToError(e.WriteTo(&b)))
	assert.Equal(t, expected, b.String())
}

func TestShadowEntry_String(t *testing.T) {
	e := ShadowEntry{
		Name: "foo",
		Min:  Int32(90),
	}
	expected := "foo:!!::90:::::\n"
	assert.Equal(t, expected, e.String())
}

func TestShadowEntry_WriteTo(t *testing.T) {
	e := ShadowEntry{
		Name: "foo",
		Min:  Int32(90),
	}
	expected := "foo:!!::90:::::\n"
	var b bytes.Buffer
	assert.Nil(t, writerToError(e.WriteTo(&b)))
	assert.Equal(t, expected, b.String())
}

func TestGroupEntry_String(t *testing.T) {
	e := GroupEntry{
		Name: "foo",
		GID:  1000,
	}
	expected := "foo:x:1000:\n"
	assert.Equal(t, expected, e.String())
}

func TestGroupEntry_WriteTo(t *testing.T) {
	e := GroupEntry{
		Name: "foo",
		GID:  1000,
	}
	expected := "foo:x:1000:\n"
	var b bytes.Buffer
	assert.Nil(t, writerToError(e.WriteTo(&b)))
	assert.Equal(t, expected, b.String())
}

func writerToError(i int64, e error) error {
	return e
}
