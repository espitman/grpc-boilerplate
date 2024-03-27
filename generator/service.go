//go:build mage
// +build mage

package main

import (
	"github.com/espitman/go-super-cli"
	"github.com/espitman/grpc-boilerplate/generator/gutil"
	"github.com/magefile/mage/sh"
	"slices"
)

type MainService struct {
	Dist     string   `yaml:"-"`
	Module   string   `yaml:"module"`
	Name     string   `yaml:"name"`
	HTTP     bool     `yaml:"http"`
	GRPC     bool     `yaml:"grpc"`
	GRPCInfo GRPCInfo `yaml:"GRPCInfo"`
	DB       DB       `yaml:"DB"`
}

type GRPCInfo struct {
	PBModule string `yaml:"PBModule"`
}

type DB struct {
	PostgreSQL bool `yaml:"postgreSQL"`
	MongoDB    bool `yaml:"mongoDB"`
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

func (m *MainService) generate() {
	m.createDirs()
	m.generateApi()
	m.generateGRPC()
	m.generateDB()

	m.generateMainFile()
	m.createYaml()
}

func (m *MainService) createDirs() {
	gutil.CreateDir(m.Dist)
	gutil.CreateDir(m.Dist + "/.info")
	gutil.CreateDir(m.Dist + "/cmd")
	gutil.CreateDir(m.Dist + "/internal")
	gutil.CreateDir(m.Dist + "/internal/adapter")
	gutil.CreateDir(m.Dist + "/internal/adapter/handler")
	gutil.CreateDir(m.Dist + "/internal/adapter/database")
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

func (m *MainService) generateDB() {
	_, selectedItems := cli.Choices(
		"Please choose the database(s) you want to use:",
		[]string{"PostgreSQL", "MongoDB"},
		false)
	isPostgreSQL := slices.Contains(selectedItems, "PostgreSQL")
	isMongoDB := slices.Contains(selectedItems, "MongoDB")
	m.DB = DB{
		PostgreSQL: isPostgreSQL,
		MongoDB:    isMongoDB,
	}
	if isPostgreSQL {
		m.generatePostgreSQL()
	}
	if isMongoDB {
		m.generateMongoDB()
	}
}

func (m *MainService) generatePostgreSQL() {

	gutil.CreateDir(m.Dist + "/internal/adapter/database/postgres")
	gutil.CreateDir(m.Dist + "/internal/adapter/database/postgres/db")
	gutil.CreateDir(m.Dist + "/internal/adapter/database/postgres/ent")
	gutil.CreateDir(m.Dist + "/internal/adapter/database/postgres/ent/schema")

	gutil.Render("../src/internal/adapter/database/postgres/db/db.tmpl", m.Dist+"/internal/adapter/database/postgres/db/db.go", m)
	gutil.Render("../src/internal/adapter/database/postgres/ent/generate.tmpl", m.Dist+"/internal/adapter/database/postgres/ent/generate.go", m)

	m.generatePostgresEnt()
}

func (m *MainService) generatePostgresEnt() error {
	cmd := "cd " + m.Dist + "/internal/adapter/database/postgres/ && go run -mod=mod entgo.io/ent/cmd/ent new User"
	sh.RunV("sh", "-c", cmd)

	cmd2 := "go generate " + m.Dist + "/internal/adapter/database/postgres/ent"
	sh.RunV("sh", "-c", cmd2)

	gutil.ReplaceImportPath(m.Dist+"/internal/adapter/database/postgres/ent", m.Name, m.Module)
	return nil
}

func (m *MainService) generateMongoDB() {
	// TODO: mongoDB
}

func (m *MainService) createYaml() {
	_ = gutil.YamlWriter(m.Dist+"/.info/service.yaml", m)
}

func (Build) Service() error {
	m := NewMainService()
	m.generate()
	return nil
}
