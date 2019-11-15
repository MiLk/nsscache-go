// Package vault is a source implementation that retrieves cache data
// from a path in a Hashicorp Vault key/value store.
package vault

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"

	"github.com/MiLk/nsscache-go/cache"
)

// Source contains the Vault API client and complete path to the cache
// data within the vault.
type Source struct {
	client    *api.Client
	prefix    string
	mountPath string
}

// Option represents a function which will make some change to the
// source during initialization.
type Option func(*Source)

// Client is an option function which will set the source's client to
// the one that is provided.
func Client(c *api.Client) Option {
	return func(s *Source) { s.client = c }
}

// Prefix is an option function which will set the source's prefix to
// the value provided.
func Prefix(p string) Option {
	return func(s *Source) { s.prefix = p }
}

// MountPath is an option function which will set the source's
// mountpath for the key/value store to the one provided.
func MountPath(m string) Option {
	return func(s *Source) { s.mountPath = m }
}

// NewSource creates a new Vault source using the options provided.
// If no options are provided a client is initialized with the default
// values.
func NewSource(opts ...Option) (*Source, error) {
	s := Source{
		prefix:    "nsscache",
		mountPath: "secret",
	}

	for _, opt := range opts {
		opt(&s)
	}

	if s.client == nil {
		cl, err := api.NewClient(nil)
		if err != nil {
			return nil, err
		}
		s.client = cl
	}

	return &s, nil
}

// Client is a convenience function to retrieve the Vault API client
// from the source.
func (s *Source) Client() *api.Client {
	return s.client
}

func (s *Source) list(name string, c *cache.Cache, createEntry func() cache.Entry) error {
	prefix := fmt.Sprintf("%s/%s", s.prefix, name)
	sec, err := s.client.Logical().List(fmt.Sprintf("%s/metadata/%s", s.mountPath, prefix))
	if err != nil {
		return errors.Wrap(err, "list from vault")
	}

	// No secret at that path
	if sec == nil {
		return nil
	}

	keys := sec.Data["keys"].([]interface{})
	for _, k := range keys {
		sec, err := s.client.Logical().Read(fmt.Sprintf("%s/data/%s/%s", s.mountPath, prefix, k))
		if err != nil {
			return errors.Wrap(err, "read from vault")
		}
		value := sec.Data["data"].(map[string]interface{})["value"].(string)
		b := bytes.NewBufferString(value)
		b64 := base64.NewDecoder(base64.StdEncoding, b)
		e := createEntry()
		err = json.NewDecoder(b64).Decode(e)
		if err != nil {
			return errors.Wrap(err, "json decoding")
		}
		c.Add(e)
	}
	return nil
}

// FillPasswdCache reads entries from the Vault and uses them to fill
// the passwd cache.
func (s *Source) FillPasswdCache(c *cache.Cache) error {
	return s.list("passwd", c, func() cache.Entry {
		return &cache.PasswdEntry{}
	})
}

// FillShadowCache reads entries from the Vault and uses them to fill
// the shadow cache.
func (s *Source) FillShadowCache(c *cache.Cache) error {
	return s.list("shadow", c, func() cache.Entry {
		return &cache.ShadowEntry{}
	})
}

// FillGroupCache reads entries from the Vault and uses them to fill
// the group cache.
func (s *Source) FillGroupCache(c *cache.Cache) error {
	return s.list("group", c, func() cache.Entry {
		return &cache.GroupEntry{}
	})
}
