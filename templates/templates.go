package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/xsevy/terrapi/helpers"
	"github.com/xsevy/terrapi/helpers/functions"
	"github.com/xsevy/terrapi/messages"
)

//go:embed source/*/*/.gitignore source/*/*/.terrapi source/*/*/resolvers/.gitkeep source/*
var sourceFiles embed.FS

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

type creationRecord struct {
	createdPaths []string
}

func (r *creationRecord) add(path string) {
	r.createdPaths = append(r.createdPaths, path)
}

func (r *creationRecord) rollback() {
	for i := len(r.createdPaths) - 1; i >= 0; i-- {
		path := r.createdPaths[i]
		os.RemoveAll(path)
	}
}

// CreateResources creates resources based on the specified ID
func CreateResources(id, dest string, replacements *messages.CreateResourceMsg) error {
	src := filepath.Join("source", id)
	var f func(src, dest string, replacements *messages.CreateResourceMsg) error

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
func createAppSyncDataSource(src, dest string, replacements *messages.CreateResourceMsg) error {
	if err := checkRequiredFields(replacements.ProjectName); err != nil {
		return err
	}

	if err := checkConfigFileExists(); err != nil {
		return err
	}

	if err := copyFiles(sourceFiles, src, dest, replacements); err != nil {
		return err
	}

	// adding new module to main file
	newModuleContent := fmt.Sprintf(newModuleContent, replacements.ProjectName, replacements.ProjectName)
	if err := functions.AppendTextToFile(terraformApiMainFileName, newModuleContent, true); err != nil {
		return err
	}

	// adding resolver and datasource to datasources file
	newDataSourceContent := fmt.Sprintf(
		newDataSourceContent,
		replacements.ProjectName,
		replacements.ProjectName,
		replacements.ProjectName,
	)
	if err := functions.AppendTextToFile(terraformDataSourcesFileName, newDataSourceContent, true); err != nil {
		return err
	}

	return nil
}

// createAppSyncApi creates resources for AppSync API
func createAppSyncApi(src, dest string, replacements *messages.CreateResourceMsg) error {
	if err := checkRequiredFields(replacements.ProjectName); err != nil {
		return err
	}

	if err := copyFiles(sourceFiles, src, dest, replacements); err != nil {
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
func checkRequiredFields(fields ...interface{}) error {
	for _, field := range fields {
		if field == nil {
			return fmt.Errorf("missing required field %s", field)
		}
	}

	return nil
}

// copyFiles copies files from src to dest
func copyFiles(fsys fs.FS, src, dest string, replacements *messages.CreateResourceMsg) error {
	record := creationRecord{}

	err := fs.WalkDir(fsys, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, _ := filepath.Rel(src, path)
		newPath, err := renameFile(relativePath, replacements)
		if err != nil {
			return err
		}
		newPath = filepath.Join(dest, newPath)

		if d.IsDir() {
			if err := createDirectory(newPath); err != nil {
				return err
			}
		} else {
			if err := createFile(fsys, path, newPath, replacements); err != nil {
				return err
			}
		}
		record.add(newPath)

		return nil
	})
	if err != nil {
		record.rollback()
		return err
	}

	return nil
}

// createDirectory creates a directory
func createDirectory(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("unable to create directory %s: %v", path, err)
	}
	return nil
}

// createFile copies a file from source and replaces values using text/template
func createFile(fsys fs.FS, src, dest string, replacements *messages.CreateResourceMsg) error {
	data, err := fs.ReadFile(fsys, src)
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

// renameFile replaces placeholders in the data string with values from the replacements struct.
func renameFile(data string, replacements *messages.CreateResourceMsg) (string, error) {
	v := reflect.ValueOf(replacements).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		placeholder := fmt.Sprintf("{{%s}}", field.Name)
		if strings.Contains(data, placeholder) {
			if value.IsValid() && !value.IsZero() {
				data = strings.ReplaceAll(data, placeholder, fmt.Sprintf("%v", value.Interface()))
			} else {
				return "", fmt.Errorf("placeholder can not be replaced %s", data)
			}
		}
	}
	return data, nil
}
