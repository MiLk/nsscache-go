// vault is an implementation of a source using Hashicorp Vault to store the data
package vault

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"

	"github.com/milk/nsscache-go/cache"
)

type VaultSource struct {
	client *api.Client
	prefix string
}

type Option func(*VaultSource)

func Client(c *api.Client) Option {
	return func(s *VaultSource) { s.client = c }
}

func Prefix(p string) Option {
	return func(s *VaultSource) { s.prefix = p }
}

func NewSource(opts ...Option) (*VaultSource, error) {
	s := VaultSource{
		prefix: "nsscache",
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

func (s *VaultSource) Client() *api.Client {
	return s.client
}

func (s *VaultSource) list(name string, c *cache.Cache, createEntry func() cache.Entry) error {
	prefix := fmt.Sprintf("secret/%s/%s", s.prefix, name)
	sec, err := s.client.Logical().List(prefix)
	if err != nil {
		return errors.Wrap(err, "list from vault")
	}

	// No secret at that path
	if sec == nil {
		return nil
	}

	keys := sec.Data["keys"].([]interface{})
	for _, k := range keys {
		sec, err := s.client.Logical().Read(fmt.Sprintf("%s/%s", prefix, k))
		if err != nil {
			return errors.Wrap(err, "read from vault")
		}
		value := sec.Data["value"].(string)
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

func (s *VaultSource) FillPasswdCache(c *cache.Cache) error {
	return s.list("passwd", c, func() cache.Entry {
		return &cache.PasswdEntry{}
	})
}

func (s *VaultSource) FillShadowCache(c *cache.Cache) error {
	return s.list("shadow", c, func() cache.Entry {
		return &cache.ShadowEntry{}
	})
}

func (s *VaultSource) FillGroupCache(c *cache.Cache) error {
	return s.list("group", c, func() cache.Entry {
		return &cache.GroupEntry{}
	})
}
