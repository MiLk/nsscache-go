// Package nsscache greatly simplifies the task of writing cache
// filling applications for libnss-cache by encapsulating the tasks of
// writing the caches, cache indexes, and doing so atomically into a
// reusable library.
package nsscache

import (
	"fmt"
	"path/filepath"

	"os"

	"github.com/MiLk/nsscache-go/cache"
	"github.com/MiLk/nsscache-go/source"
)

// CacheMap allows you to manage the caches as a group.
type CacheMap map[string]*cache.Cache

// Option is a wrapper type used to specify options on a specific
// cache.
type Option struct {
	CacheName string
	Option    cache.Option
}

// NewCaches creates cache structs for passwd, group and shadow.
func NewCaches(opts ...Option) CacheMap {
	optionMap := map[string][]cache.Option{}
	for _, opt := range opts {
		if optionMap[opt.CacheName] == nil {
			optionMap[opt.CacheName] = []cache.Option{}
		}
		optionMap[opt.CacheName] = append(optionMap[opt.CacheName], opt.Option)
	}

	m := CacheMap{}
	for _, name := range []string{"passwd", "shadow", "group"} {
		if opts, ok := optionMap[name]; ok {
			m[name] = cache.NewCache(opts...)
		} else {
			m[name] = cache.NewCache()
		}
	}
	return m
}

// FillCaches uses the provided source to fill the caches of the
// CacheMap struct.
func (cm *CacheMap) FillCaches(src source.Source) error {
	if c, ok := (*cm)["passwd"]; ok {
		if err := src.FillPasswdCache(c); err != nil {
			return err
		}
	}

	if c, ok := (*cm)["shadow"]; ok {
		if err := src.FillShadowCache(c); err != nil {
			return err
		}
	}

	if c, ok := (*cm)["group"]; ok {
		if err := src.FillGroupCache(c); err != nil {
			return err
		}
	}

	return nil
}

// WriteOptions specifies optional values for writing the caches out.
// The directory will default to '/etc' and the Extension will default
// to 'cache'.
type WriteOptions struct {
	Directory string
	Extension string
}

func defaultWriteOptions() WriteOptions {
	return WriteOptions{
		Directory: "/etc",
		Extension: "cache",
	}
}

// WriteFiles write the content of the cache structs into files that
// libnss-cache can read.
func (cm *CacheMap) WriteFiles(options *WriteOptions) error {
	wo := defaultWriteOptions()
	if options != nil {
		if options.Directory != "" {
			wo.Directory = options.Directory
		}
		if options.Extension != "" {
			wo.Extension = options.Extension
		}
	}

	for _, name := range []string{"passwd", "shadow", "group"} {
		fpath := filepath.Join(wo.Directory, fmt.Sprintf("%s.%s", name, wo.Extension))
		mode := 0644
		if name == "shadow" {
			mode = 0000
		}
		if err := WriteAtomic(fpath, (*cm)[name], os.FileMode(mode)); err != nil {
			return err
		}
	}

	idxCfg := []struct {
		cache  string
		column int
		supext string
	}{
		{"passwd", 0, "ixname"},
		{"passwd", 2, "ixuid"},
		{"group", 0, "ixname"},
		{"group", 2, "ixgid"},
		{"shadow", 0, "ixname"},
	}

	for _, idx := range idxCfg {
		fpath := filepath.Join(wo.Directory, fmt.Sprintf("%s.%s.%s", idx.cache, wo.Extension, idx.supext))
		idx := (*cm)[idx.cache].Index(idx.column)
		if err := WriteAtomic(fpath, &idx, os.FileMode(0644)); err != nil {
			return err
		}
	}

	return nil
}
