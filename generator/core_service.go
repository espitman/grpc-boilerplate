//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/espitman/go-super-cli"
	"github.com/espitman/grpc-boilerplate/generator/gutil"
)

type CoreService struct {
	Service MainService
	Name    string
	Dist    string
}

func NewCoreService() *CoreService {
	servicePath := cli.TextInput("Enter Main Service path:", "./build/x", false)
	var mainService MainService
	gutil.YamlReader(servicePath+"/.info/service.yaml", &mainService)

	name := cli.TextInput("Enter Service name:", "", false)
	dist := "../build/" + mainService.Name + "/internal/core/service/"

	return &CoreService{
		Service: mainService,
		Name:    name,
		Dist:    dist,
	}

}

func (m *CoreService) create() {
	gutil.Render(srcFolder+"/internal/core/service/service.tmpl", m.Dist+m.Name+".go", m)
}

func (Core) Service() error {
	m := NewCoreService()
	fmt.Println(m.Service.Name)
	m.create()
	return nil
}
