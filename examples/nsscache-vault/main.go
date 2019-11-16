// example showing how to use vault with nsscache to write files readable by libnss-cache
package main

import (
	"io/ioutil"
	"os"

	nsscache "github.com/MiLk/nsscache-go"
	"github.com/MiLk/nsscache-go/source/vault"
)

func main() {
	if err := mainE(); err != nil {
		panic(err)
	}
}

func mainE() error {
	// Create a temporary file to read from ONLY for test purposes
	file, err := ioutil.TempFile("/tmp", "test-token-file")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	token := []byte(os.Getenv("VAULT_TOKEN"))

	if _, err := file.Write(token); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	src, err := vault.CreateSource("nsscache_test", file.Name())
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
