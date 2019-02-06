package vault

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"

	"github.com/MiLk/nsscache-go/cache"
)

var vaultClient *api.Client

func setupTest(t *testing.T) func(t *testing.T) {
	vault.AddTestLogicalBackend("kv", kv.VersionedKVFactory)

	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)

	client, err := api.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	client.SetAddress(addr)
	client.SetToken(token)

	if err = client.Sys().TuneMount("secret", api.MountConfigInput{
		Options: map[string]string{
			"version": "2",
		},
	}); err != nil {
		t.Fatal(err)
	}
	os.Setenv("VAULT_ADDR", client.Address())
	vaultClient = client

	return func(t *testing.T) {
		ln.Close()
		os.Unsetenv("VAULT_ADDR")
		vaultClient = nil
	}
}

func addEntry(c *api.Client, mountPath, prefix, name string, e cache.Entry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = c.Logical().Write(fmt.Sprintf("%s/data/%s/%s", mountPath, prefix, name), map[string]interface{}{
		"data": map[string]interface{}{
			"value": b,
		},
	})
	return err
}

func TestNewSource(t *testing.T) {
	s, err := NewSource()
	assert.Nil(t, err)
	assert.Equal(t, "https://127.0.0.1:8200", s.client.Address())

	assert.Equal(t, s.client, s.Client())

	os.Setenv("VAULT_SKIP_VERIFY", "azerty")
	_, err = NewSource()
	assert.NotNil(t, err)
	os.Unsetenv("VAULT_SKIP_VERIFY")
}

func TestVaultSource_FillPasswdCache(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "nsscache-go-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	teardownTest := setupTest(t)
	defer teardownTest(t)

	mountPath := "secret"
	prefix := fmt.Sprintf("%s/%s", "nsscache-test", "passwd")
	entry := cache.PasswdEntry{
		Name:   "foo",
		Passwd: "x",
		UID:    1000,
		GID:    1000,
		GECOS:  "Mr Foo",
		Dir:    "/home/foo",
		Shell:  "/bin/bash",
	}
	assert.Nil(t, addEntry(vaultClient, mountPath, prefix, entry.Name, &entry))
	entry.Name = "bar"
	entry.UID = 1001
	entry.GECOS = "Mrs Bar"
	entry.Dir = "/home/bar"
	assert.Nil(t, addEntry(vaultClient, mountPath, prefix, entry.Name, &entry))

	s, err := NewSource(Client(vaultClient), MountPath(mountPath), Prefix("nsscache-test"))
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
	_, err = vaultClient.Logical().Write(fmt.Sprintf("%s/data/%s/%s", mountPath, prefix, "invalid"), map[string]interface{}{
		"data": map[string]interface{}{
			"value": base64.StdEncoding.EncodeToString([]byte("{[}}]'[")),
		},
	})
	assert.Nil(t, err)

	c = cache.NewCache()
	err = s.FillPasswdCache(c)
	assert.NotNil(t, err)

	// Base64 error
	_, err = vaultClient.Logical().Write(fmt.Sprintf("%s/data/%s/%s", mountPath, prefix, "invalid"), map[string]interface{}{
		"data": map[string]interface{}{
			"value": "foobar{[}}]'[",
		},
	})
	assert.Nil(t, err)

	c = cache.NewCache()
	err = s.FillPasswdCache(c)
	assert.NotNil(t, err)

	// Empty path
	s, err = NewSource(Client(vaultClient), Prefix("nsscache-empty"))
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

	teardownTest := setupTest(t)
	defer teardownTest(t)

	mountPath := "secret"
	prefix := fmt.Sprintf("%s/%s", "nsscache-test", "shadow")
	_, err = vaultClient.Logical().Delete(fmt.Sprintf("%s/metadata/%s", mountPath, prefix))
	assert.Nil(t, err)

	entry := cache.ShadowEntry{
		Name:   "foo",
		Passwd: "!!",
		Lstchg: cache.Int32(17321),
	}
	assert.Nil(t, addEntry(vaultClient, mountPath, prefix, entry.Name, &entry))
	entry.Name = "bar"
	assert.Nil(t, addEntry(vaultClient, mountPath, prefix, entry.Name, &entry))

	s, err := NewSource(Client(vaultClient), MountPath(mountPath), Prefix("nsscache-test"))
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

	teardownTest := setupTest(t)
	defer teardownTest(t)

	mountPath := "secret"
	prefix := fmt.Sprintf("%s/%s", "nsscache-test", "group")
	_, err = vaultClient.Logical().Delete(fmt.Sprintf("%s/metadata/%s", mountPath, prefix))
	assert.Nil(t, err)

	entry := cache.GroupEntry{
		Name:   "foo",
		Passwd: "*",
		GID:    1000,
	}
	assert.Nil(t, addEntry(vaultClient, mountPath, prefix, entry.Name, &entry))

	s, err := NewSource(Client(vaultClient), MountPath(mountPath), Prefix("nsscache-test"))
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

func TestVaultSource_List(t *testing.T) {
	s, err := NewSource()
	assert.Nil(t, err)
	err = s.list("name", nil, nil)
	assert.NotNil(t, err)
}
