//go:build mage
// +build mage

package main

import (
	"github.com/espitman/grpc-boilerplate/gutil"
	"github.com/magefile/mage/sh"
	"strings"
)

func appendToMainFile(m MainService) {
	gutil.AppendToFile(fs, srcFolder+"/cmd/main-new-repository.tmpl", m.Dist+"/cmd/main.go", "NewRepository", m)
	sh.RunV("gofmt", "-w", m.Dist+"/cmd/main.go")
}

func appendToMainFileAPI(m MainService) {
	gutil.AppendToFile(fs, srcFolder+"/cmd/main-run-service-api.tmpl", m.Dist+"/cmd/main.go", "RunServiceAPI", m)
	sh.RunV("gofmt", "-w", m.Dist+"/cmd/main.go")
}

func appendToAPiFile(m MainService) {
	gutil.AppendToFile(fs, srcFolder+"/cmd/api/api-run-type.tmpl", m.Dist+"/cmd/api/api.go", "RunType", m)
	sh.RunV("gofmt", "-w", m.Dist+"/cmd/api/api.go")
}

func appendToMainFileGRPC(m MainService) {
	gutil.AppendToFile(fs, srcFolder+"/cmd/main-run-service-grpc.tmpl", m.Dist+"/cmd/main.go", "RunServiceGRPC", m)
	sh.RunV("gofmt", "-w", m.Dist+"/cmd/main.go")

}

func appendToGRPCFile(m MainService) {
	gutil.AppendToFile(fs, srcFolder+"/cmd/gRPC/gRPC-run.tmpl", m.Dist+"/cmd/gRPC/gRPC.go", "Run", m)
	gutil.AppendToFile(fs, srcFolder+"/cmd/gRPC/gRPC-run-service.tmpl", m.Dist+"/cmd/gRPC/gRPC.go", "RunService", m)
	sh.RunV("gofmt", "-w", m.Dist+"/cmd/gRPC/gRPC.go")

}

func appendToServerFileGRPC(m MainService) {
	gutil.AppendToFile(fs, srcFolder+"/internal/adapter/handler/gRPC/server-new-server.tmpl", m.Dist+"/internal/adapter/handler/gRPC/server.go", "NewServer", m)
	gutil.AppendToFile(fs, srcFolder+"/internal/adapter/handler/gRPC/server-new-handler.tmpl", m.Dist+"/internal/adapter/handler/gRPC/server.go", "NewHandler", m)
	sh.RunV("gofmt", "-w", m.Dist+"/internal/adapter/handler/gRPC/server.go")
}

func appendToHandlerFileGRPC(m MainService) {
	gutil.AppendToFile(fs, srcFolder+"/internal/adapter/handler/gRPC/handler-handler.tmpl", m.Dist+"/internal/adapter/handler/gRPC/handler.go", "Handler", m)
	gutil.AppendToFile(fs, srcFolder+"/internal/adapter/handler/gRPC/handler-new-handler-type.tmpl", m.Dist+"/internal/adapter/handler/gRPC/handler.go", "NewHandlerType", m)
	gutil.AppendToFile(fs, srcFolder+"/internal/adapter/handler/gRPC/handler-new-handler-service.tmpl", m.Dist+"/internal/adapter/handler/gRPC/handler.go", "NewHandlerService", m)
	sh.RunV("gofmt", "-w", m.Dist+"/internal/adapter/handler/gRPC/handler.go")

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
	m.DB = mainService.DB

	if mainService.HTTP {
		appendToMainFileAPI(m)
		appendToAPiFile(m)
	}

	if mainService.GRPC {
		appendToMainFileGRPC(m)
		appendToGRPCFile(m)
		appendToServerFileGRPC(m)
		appendToHandlerFileGRPC(m)
	}

	appendToMainFile(m)

	return nil
}
