// example showing how to use vault with nsscache to write files readable by libnss-cache
package main

import (
	"os"

	nsscache "github.com/milk/nsscache-go"
	vaultsource "github.com/milk/nsscache-go/source/vault"
)

func main() {
	if err := mainE(); err != nil {
		panic(err)
	}
}

func mainE() error {
	src, err := vaultsource.CreateVaultSource("nsscache_test")
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
