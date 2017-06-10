// cache contains the struct to manipulate the cache data in memory before writing to the disk
package cache

import (
	"bytes"
	"io"
)

func NewCache() *Cache {
	return &Cache{}
}

// Cache is an in-memory struct representing the cache to be used by libnss-cache
type Cache struct {
	entries []Entry // Entries contained in the cache
}

// Add adds a new entry to the cache
func (c *Cache) Add(e Entry) {
	c.entries = append(c.entries, e)
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

// WriteTo writes the content of the cache to an io.Writer
func (c *Cache) WriteTo(w io.Writer) (int64, error) {
	b, err := c.buffer()
	if err != nil {
		return 0, err
	}
	return io.Copy(w, b)
}
