//go:build mage
// +build mage

package main

import (
	"github.com/espitman/grpc-boilerplate/gutil"
	"github.com/magefile/mage/sh"
	"strings"
)

func appendToMainFile(m MainService) {
	gutil.AppendToFile(srcFolder+"/cmd/main-new-repository.tmpl", m.Dist+"/cmd/main.go", "NewRepository", m)
	gutil.AppendToFile(srcFolder+"/cmd/main-run-service.tmpl", m.Dist+"/cmd/main.go", "RunService", m)
	sh.RunV("gofmt", "-w", m.Dist+"/cmd/main.go")
}

func appendToAPiFile(m MainService) {
	gutil.AppendToFile(srcFolder+"/cmd/api/api-run-type.tmpl", m.Dist+"/cmd/api/api.go", "RunType", m)
	sh.RunV("gofmt", "-w", m.Dist+"/cmd/api/api.go")
}

func (Core) All() error {
	d, servicePath, domain := NewCoreDomain()
	d.create()

	p := CorePort(*d)
	p.Dist = strings.Replace(d.Dist, "core/domain/", "core/port/", 1)
	p.create()

	s := CoreService(*d)
	s.Dist = strings.Replace(d.Dist, "core/domain/", "core/service/", 1)
	s.create()

	//-------------------------------------------------
	m := MainService{
		Dist:   servicePath,
		Domain: domain,
	}

	var mainService MainService
	gutil.YamlReader(servicePath+"/.info/service.yaml", &mainService)
	m.Module = mainService.Module
	m.Name = mainService.Name

	if mainService.HTTP {
		appendToAPiFile(m)
	}

	appendToMainFile(m)

	return nil
}
