// cache contains the struct to manipulate the cache data in memory before writing to the disk
package cache

import (
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

// WriteTo writes the content of the cache to an io.Writer
func (c *Cache) WriteTo(w io.Writer) (int64, error) {
	total := int64(0)
	for _, e := range c.entries {
		if n, err := e.WriteTo(w); err != nil {
			return total + n, err
		} else {
			total += n
		}
	}
	return total, nil
}
