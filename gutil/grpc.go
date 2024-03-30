package gutil

import (
	"io/ioutil"
	"regexp"
)

type Method struct {
	Name       string
	InputType  string
	OutputType string
}

func ReadProtoFile(protoFilePath string) string {
	protoBytes, err := ioutil.ReadFile(protoFilePath)
	if err != nil {
		return ""
	}

	protoContent := string(protoBytes)
	return protoContent
}

func ExtractGRPCMethods(protoFilePath string) ([]Method, error) {
	protoContent := ReadProtoFile(protoFilePath)
	methodRegex := regexp.MustCompile(`rpc\s+(\w+)\(([^)]+)\)\s+returns\s+\(([^)]+)\)`)
	matches := methodRegex.FindAllStringSubmatch(protoContent, -1)

	var extractedMethods []Method
	for _, match := range matches {
		method := Method{
			Name:       match[1],
			InputType:  match[2],
			OutputType: match[3],
		}
		extractedMethods = append(extractedMethods, method)
	}

	return extractedMethods, nil
}

func ExtractGoPackage(protoFilePath string) string {
	protoContent := ReadProtoFile(protoFilePath)
	goPackageRegex := regexp.MustCompile(`option\s+go_package\s+=\s+"([^"]+)";`)
	match := goPackageRegex.FindStringSubmatch(protoContent)
	if len(match) != 2 {
		return ""
	}

	return match[1]
}

func ExtractServiceName(protoFilePath string) string {
	protoContent, err := ioutil.ReadFile(protoFilePath)
	if err != nil {
		return ""
	}

	regex := regexp.MustCompile(`service\s+(\w+)\s*{`)
	match := regex.FindStringSubmatch(string(protoContent))
	if len(match) < 2 {
		return ""
	}

	serviceName := match[1]
	return serviceName
}
