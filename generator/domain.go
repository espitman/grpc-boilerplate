//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/espitman/go-super-cli"
	"github.com/espitman/grpc-boilerplate/generator/gutil"
)

type Domain struct {
	Service MainService
	Name    string
	Dist    string
}

func NewDomain() *Domain {
	servicePath := cli.TextInput("Enter Service path:", "./build/x", false)
	var mainService MainService
	gutil.YamlReader(servicePath+"/.info/service.yaml", &mainService)

	name := cli.TextInput("Enter Domain name:", "", false)
	dist := "../build/" + mainService.Name + "/internal/core/domain/"

	return &Domain{
		Service: mainService,
		Name:    name,
		Dist:    dist,
	}

}

func (m *Domain) create() {
	gutil.Render("../src/internal/core/domain/domain.tmpl", m.Dist+m.Name+".go", m)
}

func (Build) Domain() error {
	m := NewDomain()
	fmt.Println(m.Service.Name)
	m.create()
	return nil
}
