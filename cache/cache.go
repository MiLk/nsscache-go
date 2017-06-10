package cache

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/youtube/vitess/go/ioutil2"
)

func NewCache(name string, opts ...Option) *Cache {
	c := Cache{
		dir:  "/etc",
		name: name,
		ext:  "cache",
		perm: 0644,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

type Option func(*Cache)

// Dir is an option allowing to change the directory where the cache file will be written
func Dir(d string) Option {
	return func(c *Cache) { c.dir = d }
}

// Extension is an option allowing to change the extension of the cache file
func Extension(e string) Option {
	return func(c *Cache) { c.ext = e }
}

// Permissions is an option allowing to change the mode of the cache file
func Mode(mode os.FileMode) Option {
	return func(c *Cache) { c.perm = mode }
}

type Cache struct {
	dir     string      // Directory component of the path to the cache file
	name    string      // File name component of the path to the cache file
	ext     string      // Extension of the cache file
	perm    os.FileMode // File mode for the cache file
	entries []Entry     // Entries contained in the cache
}

func (c *Cache) Add(e Entry) {
	c.entries = append(c.entries, e)
}

func (c *Cache) filename() string {
	return path.Join(c.dir, fmt.Sprintf("%s.%s", c.name, c.ext))
}

func (c *Cache) buffer() (*bytes.Buffer, error) {
	var b bytes.Buffer
	for _, e := range c.entries {
		if _, err := b.WriteString(e.String()); err != nil {
			return nil, err
		}
	}
	return &b, nil
}

func (c *Cache) WriteTo(w io.Writer) (int64, error) {
	b, err := c.buffer()
	if err != nil {
		return 0, err
	}
	return io.Copy(w, b)
}

func (c *Cache) WriteFile() error {
	b, err := c.buffer()
	if err != nil {
		return err
	}
	return ioutil2.WriteFileAtomic(c.filename(), b.Bytes(), c.perm)
}
