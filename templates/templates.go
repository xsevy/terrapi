package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/xsevy/terrapi/helpers"
	"github.com/xsevy/terrapi/helpers/functions"
)

var (
	//go:embed source/*/*/.gitignore source/*/*/.terrapi source/*/*/resolvers/.gitkeep source/*
	sourceFiles embed.FS

	commonRequiredFields = []string{"project_name"}
)

const (
	configFileName               = ".terrapi"
	terraformApiMainFileName     = "main.tf"
	terraformDataSourcesFileName = "datasources.tf"

	newModuleContent = `
module "%s_data_source" {
  source = "./%s"
}`
	newDataSourceContent = `
resource "aws_appsync_datasource" "%s_data_source" {
  name = "${local.project_name}_%s_data_source"
  api_id = aws_appsync_graphql_api.appsync.id
  service_role_arn = aws_iam_role.iam_appsync_role.arn
  type = "AWS_LAMBDA"
  lambda_config {
    function_arn = module.%s_data_source.lambda_function_arn
  }
}`
)

// CreateResources creates resources based on the specified ID
func CreateResources(id, dest string, replacements map[string]interface{}) error {
	src := fmt.Sprintf("source/%s", id)
	var f func(src, dest string, replacements map[string]interface{}) error

	switch id {
	case helpers.ResourceIDs.CreateAppSyncAPI:
		f = createAppSyncApi
	case helpers.ResourceIDs.CreateAppSyncDataSource:
		f = createAppSyncDataSource
	default:
		return fmt.Errorf("ID %s not found in ResourceIDs", id)
	}

	return f(src, dest, replacements)
}

// createAppSyncDataSource creates resources for AppSync API
func createAppSyncDataSource(src, dest string, replacements map[string]interface{}) error {
	if err := checkRequiredFields(commonRequiredFields, replacements); err != nil {
		return err
	}

	if err := checkConfigFileExists(); err != nil {
		return err
	}

	if err := copyFiles(src, dest, replacements); err != nil {
		return err
	}

	// adding new module to main file
	newModuleContent := fmt.Sprintf(newModuleContent, replacements["project_name"], replacements["project_name"])
	if err := functions.AppendTextToFile(terraformApiMainFileName, newModuleContent); err != nil {
		return err
	}

	// adding resolver and datasource to datasources file
	newDataSourceContent := fmt.Sprintf(
		newDataSourceContent,
		replacements["project_name"],
		replacements["project_name"],
		replacements["project_name"],
	)
	if err := functions.AppendTextToFile(terraformDataSourcesFileName, newDataSourceContent); err != nil {
		return err
	}

	return nil
}

// createAppSyncApi creates resources for AppSync API
func createAppSyncApi(src, dest string, replacements map[string]interface{}) error {
	if err := checkRequiredFields(commonRequiredFields, replacements); err != nil {
		return err
	}

	if err := copyFiles(src, dest, replacements); err != nil {
		return err
	}

	return nil
}

func checkConfigFileExists() error {
	_, err := os.Stat(configFileName)
	if os.IsNotExist(err) {
		return fmt.Errorf("missing %s config file", configFileName)
	}

	return err
}

// checkRequiredFields checks if required fields are present
func checkRequiredFields(requiredFields []string, replacements map[string]interface{}) error {
	for _, field := range requiredFields {
		if _, ok := replacements[field]; !ok {
			return fmt.Errorf("missing required field %s", field)
		}
	}

	return nil
}

// copyFiles copies files from src to dest
func copyFiles(src, dest string, replacements map[string]interface{}) error {
	return fs.WalkDir(sourceFiles, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, _ := filepath.Rel(src, path)
		newPath := filepath.Join(dest, renameFile(relativePath, replacements))

		if d.IsDir() {
			return createDirectory(newPath)
		}
		return createFile(path, newPath, replacements)
	})
}

// createDirectory creates a directory
func createDirectory(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("unable to create directory %s: %v", path, err)
	}
	return nil
}

// createFile copies a file from source and replaces values using text/template
func createFile(src, dest string, replacements map[string]interface{}) error {
	data, err := sourceFiles.ReadFile(src)
	if err != nil {
		return err
	}

	tmpl, err := template.New("file").Parse(string(data))
	if err != nil {
		return fmt.Errorf("error creating template: %s %v", data, err)
	}

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", dest, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, replacements); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}

// renameFile replaces values using given map
func renameFile(data string, replacements map[string]interface{}) string {
	for key, value := range replacements {
		data = strings.ReplaceAll(data, "{{"+key+"}}", fmt.Sprintf("%v", value))
	}
	return data
}
