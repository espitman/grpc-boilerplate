//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/espitman/go-super-cli"
	gutil2 "github.com/espitman/grpc-boilerplate/gutil"
)

type CoreService struct {
	Service MainService
	Name    string
	Dist    string
	DB      DB
}

func NewCoreService() *CoreService {
	servicePath := cli.TextInput("Enter Main Service path:", "./service-name", false)
	var mainService MainService
	gutil2.YamlReader(servicePath+"/.info/service.yaml", &mainService)

	name := cli.TextInput("Enter Service name:", "", false)
	dist := "./" + mainService.Name + "/internal/core/service/"

	return &CoreService{
		Service: mainService,
		Name:    name,
		Dist:    dist,
	}

}

func (m *CoreService) create() {
	gutil2.Render(fs, srcFolder+"/internal/core/service/service.tmpl", m.Dist+m.Name+".go", m)
}

func (Core) Service() error {
	m := NewCoreService()
	fmt.Println(m.Service.Name)
	m.create()
	return nil
}
