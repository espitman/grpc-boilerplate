//go:build mage
// +build mage

package main

import (
	//"fmt"
	//"github.com/espitman/go-super-cli"
	cli "github.com/espitman/go-super-cli"
	gutil "github.com/espitman/grpc-boilerplate/gutil"
	"github.com/magefile/mage/mg"
)

type Postgre mg.Namespace

func (Postgre) Schema() error {

	dist := cli.TextInput("Enter Service path:", "./service-name", false)
	schema := cli.TextInput("Enter Schema name:", "car", false)
	m := MainService{
		Dist:   dist,
		Domain: schema,
	}

	var mainService MainService
	gutil.YamlReader(dist+"/.info/service.yaml", &mainService)
	m.Module = mainService.Module
	m.Name = mainService.Name

	m.generatePostgreSQL()
	return nil
}
