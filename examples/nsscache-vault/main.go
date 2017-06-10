// example showing how to use vault with nsscache to write files readable by libnss-cache
package main

import (
	"os"

	"github.com/hashicorp/vault/api"

	nsscache "github.com/milk/nsscache-go"
	vaultsource "github.com/milk/nsscache-go/source/vault"
)

func main() {
	if err := mainE(); err != nil {
		panic(err)
	}
}

func mainE() error {
	client, err := api.NewClient(nil)
	if err != nil {
		return err
	}

	src, err := vaultsource.NewSource(vaultsource.Client(client))
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	cm := nsscache.NewCaches()

	if err := cm.FillCaches(src); err != nil {
		return err
	}

	return cm.WriteFiles(&nsscache.WriteOptions{
		Directory: cwd,
	})
}
