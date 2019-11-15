package vault

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateVaultClientAllOk(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "test-token-file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	token := []byte("token-test")

	if _, err := file.Write(token); err != nil {
		log.Fatal(err)
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	client, err := CreateVaultClient(file.Name())

	if err != nil {
		t.Fatalf("Unexpected error creating vault source: %s", err.Error())
	}

	assert.NotNil(t, client)
}

func TestCreateVaultClientTokenIsEmpty(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "test-token-file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	token := []byte("")

	if _, err := file.Write(token); err != nil {
		log.Fatal(err)
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	client, err := CreateVaultClient(file.Name())

	assert.EqualError(t, err, "token file is empty")
	assert.Nil(t, client)
}

func TestCreateVaultClientPathError(t *testing.T) {
	client, err := CreateVaultClient("")

	assert.EqualError(t, err, "open : no such file or directory")
	assert.Nil(t, client)
}

func TestCreateVaultClientWrappedToken(t *testing.T) {
	data := map[string]interface{}{
		"token": "another-test-token",
	}

	file, err := ioutil.TempFile("/tmp", "test-token-file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	teardownTest := setupTest(t)
	defer teardownTest(t)

	wrappedData, err := vaultClient.Logical().Write("sys/wrapping/wrap", data)
	if err != nil {
		t.Fatal(err)
	}

	wtBytes, _ := json.Marshal(wrappedData.WrapInfo)

	if _, err := file.Write(wtBytes); err != nil {
		log.Fatal(err)
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	client, err := CreateVaultClient(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, client)
}

func TestCreateVaultSourceFileInput(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "test-token-file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	token := []byte("token-test")

	if _, err := file.Write(token); err != nil {
		log.Fatal(err)
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	source, err := CreateVaultSource("prefix", file.Name())
	if err != nil {
		t.Fatalf("Unexpected error creating vault source: %s", err.Error())
	}

	assert.NotNil(t, source)
}
