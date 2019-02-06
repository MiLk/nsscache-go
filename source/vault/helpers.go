package vault

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/hashicorp/vault/api"

	"github.com/MiLk/nsscache-go/source"
)

type wrappedData struct {
	Token           string `json:"token"`
	Accessor        string `json:"accessor"`
	TTL             int    `json:"ttl"`
	CreationTime    string `json:"creation_time"`
	CreationPath    string `json:"creation_path"`
	WrappedAccessor string `json:"wrapped_accessor"`
}

// CreateVaultSource returns a vault source with a client associated to work with
func CreateVaultSource(prefix string, fpath string) (source.Source, error) {
	client, err := CreateVaultClient(fpath)
	if err != nil {
		return nil, err
	}
	return NewSource(Client(client), Prefix(prefix))
}

// CreateVaultClient returns a Vault Client with a valid Token provided by the Vault Agent assigned to it.
//
// `fpath` indicates the path of the file to read from. This file is where the token provided by the agent is stored.
func CreateVaultClient(fpath string) (*api.Client, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	rawToken, err := ReadToken(fpath)
	if err != nil {
		return nil, err
	}

	var wrappedData wrappedData
	var token string

	// Check if the token has been stored in JSON format (wrapped token) or as a plain string
	if rawToken[0] == '{' {
		if err := json.Unmarshal(rawToken, &wrappedData); err == nil {
			unwrapToken := wrappedData.Token
			if unwrapToken == "" {
				return nil, errors.New("Unwrap token is empty")
			}

			secret, err := client.Logical().Unwrap(unwrapToken)
			if err != nil {
				return nil, err
			}

			if secret == nil {
				return nil, errors.New("Could not find wrapped response")
			}

			dataToken, ok := secret.Data["token"].(string)
			if !ok {
				return nil, errors.New("Key `token` was not found on the unwrapped data")
			}

			token = dataToken
		} else {
			return nil, err
		}
	} else {
		token = string(rawToken)
	}

	if token == "" {
		return nil, errors.New("Unable to fetch token from file")
	}

	client.SetToken(token)
	return client, nil
}

// ReadToken returns a byte array containing data from the designated file.
//
// `fpath` indicates the path where the file is located at.
func ReadToken(fpath string) ([]byte, error) {
	tokenFile, err := os.Open(fpath)
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
