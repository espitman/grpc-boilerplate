//go:build mage
// +build mage

package main

import (
	"github.com/espitman/go-super-cli"
	"github.com/espitman/grpc-boilerplate/generator/gutil"
	"slices"
)

type MainService struct {
	Dist     string   `yaml:"-"`
	Module   string   `yaml:"module"`
	Name     string   `yaml:"name"`
	HTTP     bool     `yaml:"http"`
	GRPC     bool     `yaml:"grpc"`
	GRPCInfo GRPCInfo `yaml:"GRPCInfo"`
}

type GRPCInfo struct {
	PBModule string `yaml:"PBModule"`
}

func NewMainService() *MainService {

	name := cli.TextInput("Enter Service name:", "", false)
	module := cli.TextInput("Enter Module path:", "github.com/user/package", false)

	_, selectedItems := cli.Choices(
		"Please choose the communication protocol(s) you want to use:",
		[]string{"HTTP", "GRPC"},
		false)

	isHttp := slices.Contains(selectedItems, "HTTP")
	isGRPC := slices.Contains(selectedItems, "GRPC")
	dist := "../build/" + name

	return &MainService{
		Dist:   dist,
		Module: module,
		Name:   name,
		HTTP:   isHttp,
		GRPC:   isGRPC,
	}

}

func (m *MainService) create() {
	m.createDirs()
	m.generateMainFile()
	m.generateApi()
	m.generateGRPC()
	m.createYaml()
}

func (m *MainService) createDirs() {
	gutil.CreateDir(m.Dist)
	gutil.CreateDir(m.Dist + "/.info")
	gutil.CreateDir(m.Dist + "/cmd")
	gutil.CreateDir(m.Dist + "/internal")
	gutil.CreateDir(m.Dist + "/internal/adapter")
	gutil.CreateDir(m.Dist + "/internal/adapter/handler")
	gutil.CreateDir(m.Dist + "/internal/core")
	gutil.CreateDir(m.Dist + "/internal/core/domain")
	gutil.CreateDir(m.Dist + "/internal/core/port")
	gutil.CreateDir(m.Dist + "/internal/core/service")

}

func (m *MainService) generateMainFile() {
	gutil.Render("../src/cmd/main.tmpl", m.Dist+"/cmd/main.go", m)
}

func (m *MainService) generateApi() {
	if m.HTTP {
		gutil.CreateDir(m.Dist + "/cmd/api")
		gutil.Render("../src/cmd/api/api.tmpl", m.Dist+"/cmd/api/api.go", m)
		gutil.CreateDir(m.Dist + "/internal/adapter/handler/http")
	}

}

func (m *MainService) generateGRPC() {
	if m.GRPC {
		m.getGRPCInfo()

		gutil.CreateDir(m.Dist + "/cmd/gRPC")
		gutil.Render("../src/cmd/gRPC/gRPC.tmpl", m.Dist+"/cmd/gRPC/gRPC.go", m)

		gutil.CreateDir(m.Dist + "/internal/adapter/handler/gRPC")
		gutil.Render("../src/internal/adapter/handler/gRPC/interceptor.tmpl", m.Dist+"/internal/adapter/handler/gRPC/interceptor.go", m)
		gutil.Render("../src/internal/adapter/handler/gRPC/server.tmpl", m.Dist+"/internal/adapter/handler/gRPC/server.go", m)
		gutil.Render("../src/internal/adapter/handler/gRPC/handler.tmpl", m.Dist+"/internal/adapter/handler/gRPC/handler.go", m)
	}
}

func (m *MainService) getGRPCInfo() {
	path := cli.TextInput("Enter Proto path:", "git.alibaba.ir/taraaz/salvation2/monorepo/pkg/protos/protogen/price_service", false)
	m.GRPCInfo.PBModule = path
}

func (m *MainService) createYaml() {
	_ = gutil.YamlWriter(m.Dist+"/.info/service.yaml", m)
}

func (Build) Service() error {
	m := NewMainService()
	m.create()
	return nil
}
