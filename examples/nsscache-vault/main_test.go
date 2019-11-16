package main

import (
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
)

func setupVault(t *testing.T) (net.Listener, error) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)

	os.Setenv("VAULT_ADDR", addr)
	os.Setenv("VAULT_TOKEN", token)

	return ln, nil
}

func TestMainE(t *testing.T) {
	ln, err := setupVault(t)
	assert.Nil(t, err)
	defer ln.Close()

	assert.Nil(t, mainE())

	cwd, err := os.Getwd()
	assert.Nil(t, err)

	stat, err := os.Stat(filepath.Join(cwd, "passwd.cache"))
	assert.Nil(t, err)
	assert.EqualValues(t, 0644, stat.Mode())

	stat, err = os.Stat(filepath.Join(cwd, "shadow.cache"))
	assert.Nil(t, err)
	assert.EqualValues(t, 0000, stat.Mode())

	stat, err = os.Stat(filepath.Join(cwd, "group.cache"))
	assert.Nil(t, err)
	assert.EqualValues(t, 0644, stat.Mode())
}
