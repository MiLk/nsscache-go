// nsscache allows you to write programs which will populate the cache files used by libnss-cache
package nsscache

import (
	"github.com/milk/nsscache-go/cache"
	"github.com/milk/nsscache-go/source"
)

// CacheMap allows you to manage the caches as a group
type CacheMap map[string]*cache.Cache

// NewCaches creates cache structs for passwd, group and shadow
func NewCaches(opts ...cache.Option) CacheMap {
	m := CacheMap{}
	m["passwd"] = cache.NewCache("passwd", opts...)
	m["group"] = cache.NewCache("group", opts...)
	shadowOpt := append(opts, cache.Mode(0000))
	m["shadow"] = cache.NewCache("shadow", shadowOpt...)
	return m
}

// FillCaches uses the provided source to fill the caches of the CacheMap struct
func (cm *CacheMap) FillCaches(src source.Source) error {
	if err := src.FillPasswdCache((*cm)["passwd"]); err != nil {
		return err
	}

	if err := src.FillShadowCache((*cm)["shadow"]); err != nil {
		return err
	}

	return src.FillGroupCache((*cm)["group"])
}

// WriteFiles write the content of the cache structs into files that libnss-cache can read
func (cm *CacheMap) WriteFiles() error {
	if err := (*cm)["passwd"].WriteFile(); err != nil {
		return err
	}
	if err := (*cm)["shadow"].WriteFile(); err != nil {
		return err
	}
	return (*cm)["group"].WriteFile()
}
