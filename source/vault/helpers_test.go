package vault

import (
	"errors"
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
		t.Fatalf("unexpected error creating vault source: %s", err.Error())
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
