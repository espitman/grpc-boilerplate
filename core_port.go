//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/espitman/go-super-cli"
	gutil2 "github.com/espitman/grpc-boilerplate/gutil"
)

type CorePort struct {
	Service MainService
	Name    string
	Dist    string
}

func NewCorePort() *CorePort {
	servicePath := cli.TextInput("Enter Main Service path:", "./build/x", false)
	var mainService MainService
	gutil2.YamlReader(servicePath+"/.info/service.yaml", &mainService)

	name := cli.TextInput("Enter Port name:", "", false)
	dist := "../build/" + mainService.Name + "/internal/core/port/"

	return &CorePort{
		Service: mainService,
		Name:    name,
		Dist:    dist,
	}

}

func (m *CorePort) create() {
	gutil2.Render(srcFolder+"/internal/core/port/service.tmpl", m.Dist+"service_"+m.Name+".go", m)
	gutil2.Render(srcFolder+"/internal/core/port/repository.tmpl", m.Dist+"repository_"+m.Name+".go", m)
}

func (Core) Port() error {
	m := NewCorePort()
	fmt.Println(m.Service.Name)
	m.create()
	return nil
}