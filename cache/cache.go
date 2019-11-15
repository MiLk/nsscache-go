// cache contains the struct to manipulate the cache data in memory before writing to the disk
package cache

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

type ACL func(e Entry) bool

type Option func(c *Cache)

func WithACL(a ACL) Option {
	return func(c *Cache) { c.acls = append(c.acls, a) }
}

func NewCache(opts ...Option) *Cache {
	c := Cache{}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// Cache is an in-memory struct representing the cache to be used by libnss-cache
type Cache struct {
	entries []Entry // Entries contained in the cache
	acls    []ACL
}

// Add adds new entries to the cache
func (c *Cache) Add(es ...Entry) {
	for _, e := range es {
		c.addOne(e)
	}
}

// addOne adds a new entry to the cache
func (c *Cache) addOne(e Entry) {
	for _, acl := range c.acls {
		if !acl(e) {
			return
		}
	}

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

// Index generates an index for the given cache on a particular
// column.  This is required for caches beyond a libnss-cache defined
// size in order for them to be read correctly.
func (c *Cache) Index(col int) bytes.Buffer {
	ordered := make([]string, len(c.entries))
	mapped := make(map[string]Entry, len(c.entries))
	for i := range c.entries {
		key := c.entries[i].Column(col)
		ordered[i] = key
		mapped[key] = c.entries[i]
	}

	// libnss-cache depends on the indexes being ordered in order
	// to accelerate the system with a binary search.
	sort.Strings(ordered)

	var b bytes.Buffer
	var offset int64
	for _, key := range ordered {
		b.WriteString(key)
		b.WriteByte(0)
		fmt.Fprintf(&b, "%08d", offset)
		for i := 0; i < 32-len(key)-1; i++ {
			b.WriteByte(0)
		}
		b.WriteString("\n")
		offset += int64(len(mapped[key].String())) + 1
	}
	return b
}
