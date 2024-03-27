package gutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func Render(tmplFile string, outputFile string, data any) {

	funcMap := template.FuncMap{
		"Time":       Time,
		"Upper":      CapitalizeFirstChar,
		"Kebab":      KebabCase,
		"ModulePath": GetModulePath,
	}

	var buffer bytes.Buffer
	tmpl, err := template.New(filepath.Base(tmplFile)).Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		os.Exit(1)
	}
	output := buffer.Bytes()
	err = ioutil.WriteFile(outputFile, output, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}

	fmt.Println("::Output written to", outputFile)
}

func ReplaceImportPath(dir string, serviceName string, newImportPath string) error {
	oldImportPath := "github.com/espitman/grpc-boilerplate/build/" + serviceName

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			err := replaceInFile(path, oldImportPath, newImportPath)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("Import path replacement completed successfully.")
	return nil
}

func replaceInFile(filePath, oldStr, newStr string) error {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	newContent := strings.ReplaceAll(string(fileBytes), oldStr, newStr)

	err = ioutil.WriteFile(filePath, []byte(newContent), 0)
	if err != nil {
		return err
	}

	return nil
}
