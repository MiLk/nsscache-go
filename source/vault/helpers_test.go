package vault

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTokenReader struct {
	result []byte
	err    error
}

func (m *mockTokenReader) ReadToken() ([]byte, error) {
	return m.result, m.err
}

func TestCreateVaultSourceWithTokenReaderAllOk(t *testing.T) {
	token := []byte("my-token")
	tr := &mockTokenReader{token, nil}
	source, err := CreateVaultSourceWithTokenReader("prefix", tr)
	if err != nil {
		t.Fatalf("Unexpected error creating vault source: %s", err.Error())
	}

	assert.NotNil(t, source)
}

func TestCreateVaultSourceWithTokenReaderTokenIsEmpty(t *testing.T) {
	token := []byte("")
	tr := &mockTokenReader{token, nil}
	source, err := CreateVaultSourceWithTokenReader("prefix", tr)

	assert.EqualError(t, err, "Unable to fetch token from file")
	assert.Nil(t, source)
}

func TestCreateVaultSourceWithTokenReaderError(t *testing.T) {
	token := []byte("my-token")
	tr := &mockTokenReader{token, errors.New("error reading token")}
	source, err := CreateVaultSourceWithTokenReader("prefix", tr)

	assert.EqualError(t, err, "error reading token")
	assert.Nil(t, source)
}

func TestCreateVaultSourceWithTokenReaderWrappedToken(t *testing.T) {
	data := map[string]interface{}{
		"token": "my-token",
	}

	teardownTest := setupTest(t)
	defer teardownTest(t)

	wrappedData, err := vaultClient.Logical().Write("sys/wrapping/wrap", data)
	if err != nil {
		t.Fatal(err)
	}

	wtBytes, _ := json.Marshal(wrappedData.WrapInfo)
	tr := &mockTokenReader{wtBytes, nil}

	source, err := CreateVaultSourceWithTokenReader("prefix", tr)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, source)
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
