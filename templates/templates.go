package templates

import (
	"embed"
	"fmt"
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

	switch id {
	case helpers.ResourceIDs.CreateAppSyncAPI:
		return createAppSyncApi(src, dest, replacements)
	case helpers.ResourceIDs.CreateAppSyncDataSource:
		return createAppSyncDataSource(src, dest, replacements)
	default:
		return fmt.Errorf("ID %s not found in ResourceIDs", id)
	}
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
		newPath := filepath.Join(dest, replaceValues(relativePath, replacements))

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

// createFile copies a file from source and replaces values
func createFile(src, dest string, replacements map[string]interface{}) error {
	data, err := sourceFiles.ReadFile(src)
	if err != nil {
		return err
	}

	newData := replaceValues(string(data), replacements)
	if err := os.WriteFile(dest, []byte(newData), 0644); err != nil {
		return fmt.Errorf("unable to create file %s: %v", dest, err)
	}
	return nil
}

// replaceValues replaces values using given map
func replaceValues(data string, replacements map[string]interface{}) string {
	for key, value := range replacements {
		data = strings.ReplaceAll(data, "{{"+key+"}}", fmt.Sprintf("%v", value))
	}
	return data
}
