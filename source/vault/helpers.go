package vault

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/MiLk/nsscache-go/source"
	"github.com/hashicorp/vault/api"
)

func CreateVaultSource(prefix string, fpath string) (source.Source, error) {
	client, err := CreateVaultClient(fpath)
	if err != nil {
		return nil, err
	}
	return NewSource(Client(client), Prefix(prefix))
}

// CreateVaultClient returns a Vault Client with a valid Token provided by the Vault Agent assigned to it.
//
// @fpath indicates the path of the file to read from. This file is where the token provided by the agent is supposed to be.
func CreateVaultClient(fpath string) (*api.Client, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	tokenFile, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer tokenFile.Close()

	rawToken, err := ioutil.ReadAll(tokenFile)
	if err != nil {
		return nil, err
	}

	var wrappedToken map[string]interface{}
	var token string

	// Check if the token has been stored in JSON format (wrapped token) or as a plain string
	if err := json.Unmarshal(rawToken, &wrappedToken); err == nil {
		secret, err := client.Logical().Unwrap(wrappedToken["token"].(string))
		if err != nil {
			return nil, err
		}

		if secret == nil {
			return nil, errors.New("Could not find wrapped response")
		}

		dataToken, ok := secret.Data["token"].(string)
		if !ok {
			return nil, errors.New("Key `token` was not found")
		}

		token = dataToken
	} else {
		token = string(rawToken)
	}

	if token == "" {
		return nil, errors.New("Unable to fetch token from file")
	}

	client.SetToken(token)
	return client, nil
}
