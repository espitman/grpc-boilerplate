//go:build mage
// +build mage

package main

import (
	"github.com/espitman/go-super-cli"
	gutil2 "github.com/espitman/grpc-boilerplate/gutil"
	"github.com/magefile/mage/sh"
	"slices"
	"strings"
)

type MainService struct {
	Dist       string     `yaml:"-"`
	Module     string     `yaml:"module"`
	Name       string     `yaml:"name"`
	HTTP       bool       `yaml:"http"`
	HTTPInfo   HTTPInfo   `yaml:"HTTPInfo"`
	GRPC       bool       `yaml:"grpc"`
	GRPCInfo   GRPCInfo   `yaml:"GRPCInfo"`
	Domain     string     `yaml:"domain"`
	DB         DB         `yaml:"DB"`
	Repository Repository `yaml:"Repository"`
}

type GRPCInfo struct {
	ProtoPath string          `yaml:"ProtoPath"`
	PBModule  string          `yaml:"PBModule"`
	Methods   []gutil2.Method `yaml:"-"`
}

type HTTPInfo struct {
	Name        string `yaml:"name"`
	ServiceName string `yaml:"serviceName"`
	DomainName  string `yaml:"domainName"`
}

type DB struct {
	PostgreSQL bool `yaml:"postgreSQL"`
	MongoDB    bool `yaml:"mongoDB"`
}

type Repository struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
}

var buildPath, srcFolder string

func init() {
	buildPath = "./"
	srcFolder = "./src"
}

func NewMainService() *MainService {

	name := cli.TextInput("Enter Service name:", "", false)
	module := cli.TextInput("Enter Module path:", "github.com/u/p", false)
	dist := buildPath + name

	m := MainService{
		Dist:   dist,
		Module: module,
		Name:   name,
	}

	m.getProtocolsInfo()
	m.getDomainInfo()
	m.getDBInfo()

	if m.HTTP {
		m.getHTTPInfo()
	}

	if m.GRPC {
		m.getGRPCInfo()
	}

	return &m

}

func (m *MainService) generate() {
	m.createDirs()

	m.generateApi()
	m.generateGRPC()

	m.generateCore()

	m.generateDB()

	m.generateMainFile()

	m.generateSwagger()
	m.generateConfig()
	m.createYaml()
}

func (m *MainService) createDirs() {
	gutil2.CreateDir(buildPath)
	gutil2.CreateDir(m.Dist)
	gutil2.CreateDir(m.Dist + "/.info")
	gutil2.CreateDir(m.Dist + "/cmd")
	gutil2.CreateDir(m.Dist + "/internal")
	gutil2.CreateDir(m.Dist + "/internal/adapter")
	gutil2.CreateDir(m.Dist + "/internal/adapter/handler")
	gutil2.CreateDir(m.Dist + "/internal/adapter/database")
	gutil2.CreateDir(m.Dist + "/internal/core")
	gutil2.CreateDir(m.Dist + "/internal/core/domain")
	gutil2.CreateDir(m.Dist + "/internal/core/port")
	gutil2.CreateDir(m.Dist + "/internal/core/service")

}

func (m *MainService) generateMainFile() {
	gutil2.Render(srcFolder+"/cmd/main.tmpl", m.Dist+"/cmd/main.go", m)
}

func (m *MainService) getProtocolsInfo() {
	_, selectedItems := cli.Choices(
		"Please choose the communication protocol(s) you want to use:",
		[]string{"HTTP", "GRPC"},
		false)

	m.HTTP = slices.Contains(selectedItems, "HTTP")
	m.GRPC = slices.Contains(selectedItems, "GRPC")
}

func (m *MainService) getDomainInfo() {
	m.Domain = cli.TextInput("Enter Domain name:", "", false)
}

func (m *MainService) generateCore() {

	dist := buildPath + "/" + m.Name + "/internal/core/domain/"

	d := CoreDomain{
		Service: *m,
		Name:    m.Domain,
		Dist:    dist,
	}
	d.create()

	p := CorePort(d)
	p.Dist = strings.Replace(d.Dist, "core/domain/", "core/port/", 1)
	p.create()

	s := CoreService(d)
	s.Dist = strings.Replace(d.Dist, "core/domain/", "core/service/", 1)
	s.create()

}

func (m *MainService) generateApi() {
	if m.HTTP {
		//m.getHTTPInfo()
		gutil2.CreateDir(m.Dist + "/cmd/api")
		gutil2.Render(srcFolder+"/cmd/api/api.tmpl", m.Dist+"/cmd/api/api.go", m)

		gutil2.CreateDir(m.Dist + "/internal/adapter/handler/http")

		gutil2.Render(srcFolder+"/internal/adapter/handler/http/dto_global.tmpl", m.Dist+"/internal/adapter/handler/http/dto_global.go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/http/middleware_authorize.tmpl", m.Dist+"/internal/adapter/handler/http/middleware_authorize.go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/http/router.tmpl", m.Dist+"/internal/adapter/handler/http/router.go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/http/server.tmpl", m.Dist+"/internal/adapter/handler/http/server.go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/http/swagger.tmpl", m.Dist+"/internal/adapter/handler/http/swagger.go", m)

		gutil2.Render(srcFolder+"/internal/adapter/handler/http/dto_name.tmpl", m.Dist+"/internal/adapter/handler/http/dto_"+m.HTTPInfo.Name+".go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/http/handler_name.tmpl", m.Dist+"/internal/adapter/handler/http/handler_"+m.HTTPInfo.Name+".go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/http/mapper_name.tmpl", m.Dist+"/internal/adapter/handler/http/mapper_"+m.HTTPInfo.Name+".go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/http/router_name.tmpl", m.Dist+"/internal/adapter/handler/http/router_"+m.HTTPInfo.Name+".go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/http/validator_name.tmpl", m.Dist+"/internal/adapter/handler/http/validator_"+m.HTTPInfo.Name+".go", m)
	}

}

func (m *MainService) getHTTPInfo() {
	//handlerName := cli.TextInput("Enter Http Handler name:", "handler", false)
	m.HTTPInfo = HTTPInfo{
		Name:        m.Domain,
		ServiceName: m.Domain,
		DomainName:  m.Domain,
	}
}

func (m *MainService) generateGRPC() {
	if m.GRPC {

		m.GRPCInfo.PBModule = gutil2.ExtractGoPackage(m.GRPCInfo.ProtoPath)
		methods, _ := gutil2.ExtractGRPCMethods(m.GRPCInfo.ProtoPath)
		m.GRPCInfo.Methods = methods

		gutil2.CreateDir(m.Dist + "/cmd/gRPC")
		gutil2.Render(srcFolder+"/cmd/gRPC/gRPC.tmpl", m.Dist+"/cmd/gRPC/gRPC.go", m)

		gutil2.CreateDir(m.Dist + "/internal/adapter/handler/gRPC")
		gutil2.Render(srcFolder+"/internal/adapter/handler/gRPC/interceptor.tmpl", m.Dist+"/internal/adapter/handler/gRPC/interceptor.go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/gRPC/server.tmpl", m.Dist+"/internal/adapter/handler/gRPC/server.go", m)
		gutil2.Render(srcFolder+"/internal/adapter/handler/gRPC/handler.tmpl", m.Dist+"/internal/adapter/handler/gRPC/handler.go", m)

		gutil2.Render(srcFolder+"/internal/adapter/handler/gRPC/handler_name.tmpl", m.Dist+"/internal/adapter/handler/gRPC/handler_"+m.Domain+".go", m)
	}
}

func (m *MainService) getGRPCInfo() {
	path := cli.TextInput("Enter Proto file path:", "*/*.proto", false)
	m.GRPCInfo.ProtoPath = path
}

func (m *MainService) getDBInfo() {
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
}

func (m *MainService) generateDB() {
	if m.DB.PostgreSQL {
		m.generatePostgreSQL()
	}
	if m.DB.MongoDB {
		m.generateMongoDB()
	}
}

func (m *MainService) generatePostgreSQL() {

	gutil2.CreateDir(m.Dist + "/internal/adapter/database/postgres")
	gutil2.CreateDir(m.Dist + "/internal/adapter/database/postgres/db")
	gutil2.CreateDir(m.Dist + "/internal/adapter/database/postgres/ent")
	gutil2.CreateDir(m.Dist + "/internal/adapter/database/postgres/ent/schema")

	gutil2.Render(srcFolder+"/internal/adapter/database/postgres/db/db.tmpl", m.Dist+"/internal/adapter/database/postgres/db/db.go", m)
	gutil2.Render(srcFolder+"/internal/adapter/database/postgres/ent/generate.tmpl", m.Dist+"/internal/adapter/database/postgres/ent/generate.go", m)

	m.generatePostgresEnt()

	gutil2.Render(srcFolder+"/internal/adapter/database/postgres/db/repository.tmpl", m.Dist+"/internal/adapter/database/postgres/db/repository_"+m.Domain+".go", m)
	gutil2.Render(srcFolder+"/internal/adapter/database/postgres/db/mapper.tmpl", m.Dist+"/internal/adapter/database/postgres/db/mapper_"+m.Domain+".go", m)

}

func (m *MainService) generatePostgresEnt() error {
	gutil2.Render(srcFolder+"/internal/adapter/database/postgres/ent/schema/schema.tmpl", m.Dist+"/internal/adapter/database/postgres/ent/schema/"+m.Domain+".go", m)

	cmd := "go generate " + m.Dist + "/internal/adapter/database/postgres/ent"
	sh.RunV("sh", "-c", cmd)

	gutil2.ReplaceImportPath(m.Dist+"/internal/adapter/database/postgres/ent", m.Name, m.Module)
	return nil
}

func (m *MainService) generateMongoDB() {
	// TODO: mongoDB
}

func (m *MainService) generateSwagger() {
	if m.HTTP {
		cmd := "swag init -g ./" + m.Name + "/cmd/api/api.go -o ./" + m.Name + "/cmd/api/docs --parseDependency"
		sh.RunV("sh", "-c", cmd)
	}
}

func (m *MainService) generateConfig() {
	gutil2.CreateDir(m.Dist + "/config")
	gutil2.Render(srcFolder+"/config/default.json.tmpl", m.Dist+"/config/default.json", m)
}

func (m *MainService) createYaml() {
	_ = gutil2.YamlWriter(m.Dist+"/.info/service.yaml", m)
}

func (Build) Service() error {
	m := NewMainService()
	m.generate()
	////fmt.Println(m)
	return nil
}