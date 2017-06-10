package main

import (
	"os"

	"github.com/hashicorp/vault/api"

	"github.com/milk/nsscache-go/cache"
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
	dirOption := cache.Dir(cwd)
	passwd := cache.NewCache("passwd", dirOption)
	shadow := cache.NewCache("shadow", dirOption, cache.Mode(0000))
	group := cache.NewCache("group", dirOption)

	if err := src.FillPasswdCache(passwd); err != nil {
		return err
	}
	if err := src.FillShadowCache(shadow); err != nil {
		return err
	}
	if err := src.FillGroupCache(group); err != nil {
		return err
	}

	if err := passwd.WriteFile(); err != nil {
		return err
	}
	if err := shadow.WriteFile(); err != nil {
		return err
	}
	if err := group.WriteFile(); err != nil {
		return err
	}

	return nil
}
