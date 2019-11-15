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

type Source struct {
	client    *api.Client
	prefix    string
	mountPath string
}

type Option func(*Source)

func Client(c *api.Client) Option {
	return func(s *Source) { s.client = c }
}

func Prefix(p string) Option {
	return func(s *Source) { s.prefix = p }
}

func MountPath(m string) Option {
	return func(s *Source) { s.mountPath = m }
}

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

func (s *Source) FillPasswdCache(c *cache.Cache) error {
	return s.list("passwd", c, func() cache.Entry {
		return &cache.PasswdEntry{}
	})
}

func (s *Source) FillShadowCache(c *cache.Cache) error {
	return s.list("shadow", c, func() cache.Entry {
		return &cache.ShadowEntry{}
	})
}

func (s *Source) FillGroupCache(c *cache.Cache) error {
	return s.list("group", c, func() cache.Entry {
		return &cache.GroupEntry{}
	})
}
