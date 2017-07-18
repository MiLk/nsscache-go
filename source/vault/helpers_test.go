package vault

import (
	"os"
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
