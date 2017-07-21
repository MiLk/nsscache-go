package vault

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateVaultClient(t *testing.T) {
	os.Setenv("VAULT_TOKEN", "test")
	client, err := CreateVaultClient()
	assert.NotNil(t, client)
	assert.Nil(t, err)
}

func TestCreateVaultSource(t *testing.T) {
	os.Setenv("VAULT_TOKEN", "test")
	source, err := CreateVaultSource("nsscache")
	assert.NotNil(t, source)
	assert.Nil(t, err)
}

func TestCreateVaultSource2(t *testing.T) {
	os.Setenv("VAULT_SKIP_VERIFY", "qwerty")
	source, err := CreateVaultSource("nsscache")
	assert.Nil(t, source)
	assert.NotNil(t, err)
	os.Unsetenv("VAULT_SKIP_VERIFY")
}

func TestCreateVaultSource3(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	os.Setenv("VAULT_AUTH_GITHUB_TOKEN", "qweryty")
	source, err := CreateVaultSource("nsscache")
	assert.Nil(t, source)
	assert.NotNil(t, err)
	assert.Equal(t, true, strings.Contains(err.Error(), "Put https://127.0.0.1:8200/v1/auth/github/login:"))
}
