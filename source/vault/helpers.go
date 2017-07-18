package vault

import (
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/milk/nsscache-go/source"
)

func CreateVaultSource(prefix string) (source.Source, error) {
	client, err := CreateVaultClient()
	if err != nil {
		return nil, err
	}
	return NewSource(Client(client), Prefix(prefix))
}

func CreateVaultClient() (*api.Client, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}
	// If there is no token read from the environment
	if client.Token() == "" {
		ghToken := os.Getenv("VAULT_AUTH_GITHUB_TOKEN")
		if ghToken != "" {
			// Try to authenticate with GitHub if there is a token
			secret, err := client.Logical().Write("auth/github/login", map[string]interface{}{
				"token": ghToken,
			})
			if err != nil {
				return nil, err
			}
			client.SetToken(secret.Auth.ClientToken)
		} else {
			// Try to authenticate with EC2
			pkcs7, err := getPkcs7()
			if err != nil {
				return nil, err
			}
			nonce, err := getNonce("/etc/vault-nonce")
			if err != nil {
				return nil, err
			}
			secret, err := client.Logical().Write("auth/aws-ec2/login", map[string]interface{}{
				"role":  "nsscache",
				"pkcs7": pkcs7,
				"nonce": nonce,
			})
			if err != nil {
				return nil, err
			}
			client.SetToken(secret.Auth.ClientToken)
		}
	}
	return client, nil
}
