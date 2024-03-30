//go:build mage
// +build mage

package main

import (
	//"fmt"
	"github.com/espitman/go-super-cli"
	gutil "github.com/espitman/grpc-boilerplate/gutil"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Api mg.Namespace

func appendToMainFile(m MainService) {
	gutil.AppendToFile(srcFolder+"/cmd/main-new-repository.tmpl", m.Dist+"/cmd/main.go", "NewRepository", m)
	gutil.AppendToFile(srcFolder+"/cmd/main-run-service.tmpl", m.Dist+"/cmd/main.go", "RunService", m)
	sh.RunV("gofmt", "-w", m.Dist+"/cmd/main.go")
}

func appendToAPiFile(m MainService) {
	gutil.AppendToFile(srcFolder+"/cmd/api/api-run-type.tmpl", m.Dist+"/cmd/api/api.go", "RunType", m)
	gutil.AppendToFile(srcFolder+"/cmd/api/api-new-handler.tmpl", m.Dist+"/cmd/api/api.go", "NewHandler", m)
	gutil.AppendToFile(srcFolder+"/cmd/api/api-new-server-handler.tmpl", m.Dist+"/cmd/api/api.go", "NewServerHandler", m)
	sh.RunV("gofmt", "-w", m.Dist+"/cmd/api/api.go")
}

func appendToServerFile(m MainService) {
	gutil.AppendToFile(srcFolder+"/internal/adapter/handler/http/server-server-type.tmpl", m.Dist+"/internal/adapter/handler/http/server.go", "ServerType", m)
	gutil.AppendToFile(srcFolder+"/internal/adapter/handler/http/server-new-server-type.tmpl", m.Dist+"/internal/adapter/handler/http/server.go", "NewServerType", m)
	gutil.AppendToFile(srcFolder+"/internal/adapter/handler/http/server-server-handler.tmpl", m.Dist+"/internal/adapter/handler/http/server.go", "ServerHandler", m)
	gutil.AppendToFile(srcFolder+"/internal/adapter/handler/http/server-router-run.tmpl", m.Dist+"/internal/adapter/handler/http/server.go", "RouterRun", m)
	sh.RunV("gofmt", "-w", m.Dist+"/internal/adapter/handler/http/server.go")
}

func appendToRouterFile(m MainService) {
	gutil.AppendToFile(srcFolder+"/internal/adapter/handler/http/router-router-type.tmpl", m.Dist+"/internal/adapter/handler/http/router.go", "RouterType", m)
	gutil.AppendToFile(srcFolder+"/internal/adapter/handler/http/router-new-router-type.tmpl", m.Dist+"/internal/adapter/handler/http/router.go", "NewRouterType", m)
	gutil.AppendToFile(srcFolder+"/internal/adapter/handler/http/router-router-handler.tmpl", m.Dist+"/internal/adapter/handler/http/router.go", "RouterHandler", m)
	gutil.AppendToFile(srcFolder+"/internal/adapter/handler/http/router-router-serve.tmpl", m.Dist+"/internal/adapter/handler/http/router.go", "RouterServe", m)
	sh.RunV("gofmt", "-w", m.Dist+"/internal/adapter/handler/http/router.go")
}

func (Api) Add() error {

	dist := cli.TextInput("Enter Service path:", "./service-name", false)
	handler := cli.TextInput("Enter Handler name:", "car", false)

	m := MainService{
		Dist:   dist,
		Domain: handler,
		HTTPInfo: HTTPInfo{
			Name:        handler,
			ServiceName: handler,
			DomainName:  handler,
		},
	}

	var mainService MainService
	gutil.YamlReader(dist+"/.info/service.yaml", &mainService)
	m.Module = mainService.Module
	m.Name = mainService.Name

	appendToRouterFile(m)
	appendToServerFile(m)
	appendToAPiFile(m)
	appendToMainFile(m)

	gutil.Render(srcFolder+"/internal/adapter/handler/http/dto_name.tmpl", m.Dist+"/internal/adapter/handler/http/dto_"+m.HTTPInfo.Name+".go", m)
	gutil.Render(srcFolder+"/internal/adapter/handler/http/handler_name.tmpl", m.Dist+"/internal/adapter/handler/http/handler_"+m.HTTPInfo.Name+".go", m)
	gutil.Render(srcFolder+"/internal/adapter/handler/http/mapper_name.tmpl", m.Dist+"/internal/adapter/handler/http/mapper_"+m.HTTPInfo.Name+".go", m)
	gutil.Render(srcFolder+"/internal/adapter/handler/http/router_name.tmpl", m.Dist+"/internal/adapter/handler/http/router_"+m.HTTPInfo.Name+".go", m)
	gutil.Render(srcFolder+"/internal/adapter/handler/http/validator_name.tmpl", m.Dist+"/internal/adapter/handler/http/validator_"+m.HTTPInfo.Name+".go", m)

	return nil
}
