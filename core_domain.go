//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/espitman/go-super-cli"
	gutil2 "github.com/espitman/grpc-boilerplate/gutil"
	"github.com/magefile/mage/mg"
)

type Core mg.Namespace

type CoreDomain struct {
	Service MainService
	Name    string
	Dist    string
}

func NewCoreDomain() *CoreDomain {
	servicePath := cli.TextInput("Enter Main Service path:", "./build/x", false)
	var mainService MainService
	gutil2.YamlReader(servicePath+"/.info/service.yaml", &mainService)

	name := cli.TextInput("Enter Domain name:", "", false)
	dist := "../build/" + mainService.Name + "/internal/core/domain/"

	return &CoreDomain{
		Service: mainService,
		Name:    name,
		Dist:    dist,
	}

}

func (m *CoreDomain) create() {
	gutil2.Render(srcFolder+"/internal/core/domain/domain.tmpl", m.Dist+m.Name+".go", m)
}

func (Core) Domain() error {
	m := NewCoreDomain()
	fmt.Println(m.Service.Name)
	m.create()
	return nil
}
