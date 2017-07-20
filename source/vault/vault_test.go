package vault

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"

	"github.com/MiLk/nsscache-go/cache"
)

func setupVault(t *testing.T) (net.Listener, *api.Client, error) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)

	client, err := api.NewClient(nil)
	if err != nil {
		return ln, nil, err
	}

	client.SetAddress(addr)
	client.SetToken(token)

	return ln, client, nil

}

func addEntry(c *api.Client, prefix, name string, e cache.Entry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = c.Logical().Write(fmt.Sprintf("%s/%s", prefix, name), map[string]interface{}{
		"value": b,
	})
	return err
}

func TestNewSource(t *testing.T) {
	s, err := NewSource()
	assert.Nil(t, err)
	assert.Equal(t, "https://127.0.0.1:8200", s.client.Address())

	os.Setenv("VAULT_SKIP_VERIFY", "azerty")
	_, err = NewSource()
	assert.NotNil(t, err)
	os.Unsetenv("VAULT_SKIP_VERIFY")
}

func TestVaultSource_FillPasswdCache(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	ln, client, err := setupVault(t)
	defer ln.Close()
	assert.Nil(t, err)

	prefix := fmt.Sprintf("secret/%s/%s", "nsscache-test", "passwd")
	entry := cache.PasswdEntry{
		Name:   "foo",
		Passwd: "x",
		UID:    1000,
		GID:    1000,
		GECOS:  "Mr Foo",
		Dir:    "/home/foo",
		Shell:  "/bin/bash",
	}
	assert.Nil(t, addEntry(client, prefix, entry.Name, &entry))
	entry.Name = "bar"
	entry.UID = 1001
	entry.GECOS = "Mrs Bar"
	entry.Dir = "/home/bar"
	assert.Nil(t, addEntry(client, prefix, entry.Name, &entry))

	s, err := NewSource(Client(client), Prefix("nsscache-test"))
	assert.Nil(t, err)

	c := cache.NewCache()
	err = s.FillPasswdCache(c)
	assert.Nil(t, err)

	var b bytes.Buffer
	n, err := c.WriteTo(&b)
	assert.Nil(t, err)
	assert.EqualValues(t, 87, n)
	expected := `bar:x:1001:1000:Mrs Bar:/home/bar:/bin/bash
foo:x:1000:1000:Mr Foo:/home/foo:/bin/bash
`
	assert.Equal(t, expected, b.String())

	// JSON error
	_, err = client.Logical().Write(fmt.Sprintf("%s/%s", prefix, "invalid"), map[string]interface{}{
		"value": base64.StdEncoding.EncodeToString([]byte("{[}}]'[")),
	})
	assert.Nil(t, err)

	c = cache.NewCache()
	err = s.FillPasswdCache(c)
	assert.NotNil(t, err)

	// Base64 error
	_, err = client.Logical().Write(fmt.Sprintf("%s/%s", prefix, "invalid"), map[string]interface{}{
		"value": "foobar{[}}]'[",
	})
	assert.Nil(t, err)

	c = cache.NewCache()
	err = s.FillPasswdCache(c)
	assert.NotNil(t, err)

	// Empty path
	s, err = NewSource(Client(client), Prefix("nsscache-empty"))
	assert.Nil(t, err)

	c = cache.NewCache()
	err = s.FillPasswdCache(c)
	assert.Nil(t, err)

	b.Reset()
	n, err = c.WriteTo(&b)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, n)
	assert.Equal(t, "", b.String())
}

func TestVaultSource_FillShadowCache(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	ln, client, err := setupVault(t)
	defer ln.Close()
	assert.Nil(t, err)

	prefix := fmt.Sprintf("secret/%s/%s", "nsscache-test", "shadow")
	_, err = client.Logical().Delete(prefix)
	assert.Nil(t, err)

	entry := cache.ShadowEntry{
		Name:   "foo",
		Passwd: "!!",
		Lstchg: cache.Int32(17321),
	}
	assert.Nil(t, addEntry(client, prefix, entry.Name, &entry))
	entry.Name = "bar"
	assert.Nil(t, addEntry(client, prefix, entry.Name, &entry))

	s, err := NewSource(Client(client), Prefix("nsscache-test"))
	assert.Nil(t, err)

	c := cache.NewCache()
	err = s.FillShadowCache(c)
	assert.Nil(t, err)

	var b bytes.Buffer
	n, err := c.WriteTo(&b)
	assert.Nil(t, err)
	assert.EqualValues(t, 38, n)
	expected := `bar:!!:17321::::::
foo:!!:17321::::::
`
	assert.Equal(t, expected, b.String())
}

func TestVaultSource_FillGroupCache(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	ln, client, err := setupVault(t)
	defer ln.Close()
	assert.Nil(t, err)

	prefix := fmt.Sprintf("secret/%s/%s", "nsscache-test", "group")
	_, err = client.Logical().Delete(prefix)
	assert.Nil(t, err)

	entry := cache.GroupEntry{
		Name:   "foo",
		Passwd: "*",
		GID:    1000,
	}
	assert.Nil(t, addEntry(client, prefix, entry.Name, &entry))

	s, err := NewSource(Client(client), Prefix("nsscache-test"))
	assert.Nil(t, err)

	c := cache.NewCache()
	err = s.FillGroupCache(c)
	assert.Nil(t, err)

	var b bytes.Buffer
	n, err := c.WriteTo(&b)
	assert.Nil(t, err)
	assert.EqualValues(t, 12, n)
	expected := `foo:*:1000:
`
	assert.Equal(t, expected, b.String())
}
