// nsscache allows you to write programs which will populate the cache files used by libnss-cache
package nsscache

import (
	"fmt"
	"path"

	"os"

	"github.com/milk/nsscache-go/cache"
	"github.com/milk/nsscache-go/source"
)

// CacheMap allows you to manage the caches as a group
type CacheMap map[string]*cache.Cache

// NewCaches creates cache structs for passwd, group and shadow
func NewCaches() CacheMap {
	m := CacheMap{}
	for _, name := range []string{"passwd", "shadow", "group"} {
		m[name] = cache.NewCache()
	}
	return m
}

// FillCaches uses the provided source to fill the caches of the CacheMap struct
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

// WriteFiles write the content of the cache structs into files that libnss-cache can read
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
		filepath := path.Join(wo.Directory, fmt.Sprintf("%s.%s", name, wo.Extension))
		mode := 0644
		if name == "shadow" {
			mode = 0000
		}
		if err := WriteAtomic(filepath, (*cm)[name], os.FileMode(mode)); err != nil {
			return err
		}
	}

	return nil
}
