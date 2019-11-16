// Package source defines the interfaces that fill caches.  Any system
// that can provide the information for the caches can be plugged into
// libnss-go by implementing the Source interface.
package source

import "github.com/MiLk/nsscache-go/cache"

// PasswdSource can be satisfied by a type that provide a function for
// filling the passwd cache.  Remember that in most implementations, a
// passwd entry will not be valid without a corresponding shadow
// entry.
type PasswdSource interface {
	FillPasswdCache(c *cache.Cache) error
}

// ShadowSource is satisfied by a type that provides a function for
// filling the shadow cache.
type ShadowSource interface {
	FillShadowCache(*cache.Cache) error
}

// GroupSource is satisfied by a type that provides a function for
// filling the group cache.  Since now sgroup cache is provided, if
// you require passwords in your groups you will need to specify them
// here in an appropriate format.
type GroupSource interface {
	FillGroupCache(*cache.Cache) error
}

// A Source is a type that is capable of completely filling the caches
// for passwd, group, and shadow.  Consumers of libnss-go should
// implement this interface.
type Source interface {
	PasswdSource
	ShadowSource
	GroupSource
}
