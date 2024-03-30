package gutil

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"Time":       Time,
	"Upper":      Upper,
	"Kebab":      KebabCase,
	"ModulePath": GetModulePath,
}

func Render(fs embed.FS, tmplFile string, outputFile string, data any) {

	tmplFile = strings.Replace(tmplFile, "./src", "src", 1)

	content, err := fs.ReadFile(tmplFile)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buffer bytes.Buffer
	tmpl, err := template.New(filepath.Base(tmplFile)).Funcs(funcMap).Parse(string(content))
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

	//oldImportPath := "github.com/espitman/grpc-boilerplate/" + serviceName
	moduleName, _ := GetModuleName()
	oldImportPath := moduleName + "/" + serviceName

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

func AppendToFile(fs embed.FS, templateFilePath string, goFilePath string, comment string, data any) error {
	comment = "// +salvation " + comment
	templateFilePath = strings.Replace(templateFilePath, "./src", "src", 1)

	content, err := fs.ReadFile(templateFilePath)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	goFileContent, err := ioutil.ReadFile(goFilePath)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	tmpl, err := template.New(filepath.Base(templateFilePath)).Funcs(funcMap).Parse(string(content))
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
	updatedContent := strings.Replace(string(goFileContent), comment, string(output), 1)
	err = ioutil.WriteFile(goFilePath, []byte(updatedContent), 0644)
	if err != nil {
		return err
	}
	fmt.Println("::Output written to", goFilePath)
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

func GetModuleName() (string, error) {
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'go list' command: %v", err)
	}

	moduleName := strings.TrimSpace(string(output))
	return moduleName, nil
}
