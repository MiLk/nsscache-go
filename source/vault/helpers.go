package vault

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/MiLk/nsscache-go/source"
	"github.com/hashicorp/vault/api"
)

type tokenReader interface {
	ReadToken() ([]byte, error)
}

type fileTokenReader struct {
	path string
}

func (f *fileTokenReader) ReadToken() ([]byte, error) {
	tokenFile, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer tokenFile.Close()

	rawToken, err := ioutil.ReadAll(tokenFile)
	if err != nil {
		return nil, err
	}

	return rawToken, nil
}

type wrappedToken struct {
	Token           string `json:"token"`
	Accessor        string `json:"accessor"`
	TTL             int    `json:"ttl"`
	CreationTime    string `json:"creation_time"`
	CreationPath    string `json:"creation_path"`
	WrappedAccessor string `json:"wrapped_accessor"`
}

func CreateVaultSource(prefix string, fpath string) (source.Source, error) {
	return CreateVaultSourceWithTokenReader(prefix, &fileTokenReader{path: fpath})
}

/**
TODO:
1) mover el codigo que lee de un fichero a una funcion a parte.
2) leer sobre `interfaces` go
3) crear una interfaz que defina una funcion que haga lo de 1)
4) crear un `struct` que implemente esa funcion (la del paso 1)
5) dependency injection golang
*/
// CreateVaultClient returns a Vault Client with a valid Token provided by the Vault Agent assigned to it.
//
// @fpath indicates the path of the file to read from. This file is where the token provided by the agent is supposed to be.
func CreateVaultSourceWithTokenReader(prefix string, tr tokenReader) (source.Source, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	rawToken, err := tr.ReadToken()
	if err != nil {
		return nil, err
	}

	var wrappedToken wrappedToken
	var token string

	// Check if the token has been stored in JSON format (wrapped token) or as a plain string
	if err := json.Unmarshal(rawToken, &wrappedToken); err == nil {
		wt := wrappedToken.Token
		if wt == "" {
			return nil, errors.New("Key `token` is missing")
		}

		secret, err := client.Logical().Unwrap(wt)
		if err != nil {
			return nil, err
		}

		if secret == nil {
			return nil, errors.New("Could not find wrapped response")
		}

		dataToken, ok := secret.Data["token"].(string)
		if !ok {
			return nil, errors.New("Key `token` was not found on the response")
		}

		token = dataToken
	} else {
		token = string(rawToken)
	}

	if token == "" {
		return nil, errors.New("Unable to fetch token from file")
	}

	client.SetToken(token)
	return NewSource(Client(client), Prefix(prefix))
}
