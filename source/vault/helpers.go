package vault

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/MiLk/nsscache-go/source"
	"github.com/hashicorp/vault/api"
)

func CreateVaultSource(prefix string) (source.Source, error) {
	client, err := CreateVaultClient("/etc/token-via-agent")
	if err != nil {
		return nil, err
	}
	return NewSource(Client(client), Prefix(prefix))
}

// CreateVaultClient returns a Vault Client with a valid Token provided by the Vault Agent assigned to it.
// The token must be wrapped by the sink.
//
// @path indicates the path of the file to read from. This file is where the token provided by the agent is supposed to be.
func CreateVaultClient(path string) (*api.Client, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	tokenFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer tokenFile.Close()

	tokenRead, _ := ioutil.ReadAll(tokenFile)

	var wrappedToken map[string]interface{}
	json.Unmarshal(tokenRead, &wrappedToken)

	secret, err := client.Logical().Unwrap(wrappedToken["token"].(string))
	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, errors.New("Could not find wrapped response")
	}

	token, ok := secret.Data["token"]
	if !ok {
		return nil, errors.New("Key `token` was not found")
	}

	client.SetToken(token.(string))
	return client, nil
}
