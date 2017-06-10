// source contains the interfaces to implement to create a new source
package source

import "github.com/milk/nsscache-go/cache"

type PasswdSource interface {
	FillPasswdCache(c *cache.Cache) error
}

type ShadowSource interface {
	FillShadowCache(*cache.Cache) error
}

type GroupSource interface {
	FillGroupCache(*cache.Cache) error
}

type Source interface {
	PasswdSource
	ShadowSource
	GroupSource
}
