package templates

import (
	"fmt"
	"os"
	"testing"
	"testing/fstest"

	"github.com/xsevy/terrapi/messages"
)

func TestRenameFile(t *testing.T) {
	tcs := []struct {
		replacements messages.CreateResourceMsg
		name         string
		data         string
		expected     string
		expectError  bool
	}{
		{
			name:        "Basic replacement",
			data:        "{{ProjectName}}",
			expected:    "project_name_test",
			expectError: false,
			replacements: messages.CreateResourceMsg{
				ProjectName: "project_name_test",
			},
		},
		{
			name:        "Multiple replacements",
			data:        "{{ProjectName}}{{ProjectName}}",
			expected:    "project_name_testproject_name_test",
			expectError: false,
			replacements: messages.CreateResourceMsg{
				ProjectName: "project_name_test",
			},
		},
		{
			name:        "Replacement with text before",
			data:        "xyz{{ProjectName}}",
			expected:    "xyzproject_name_test",
			expectError: false,
			replacements: messages.CreateResourceMsg{
				ProjectName: "project_name_test",
			},
		},
		{
			name:        "Replacement with text after",
			data:        "{{ProjectName}}xyz",
			expected:    "project_name_testxyz",
			expectError: false,
			replacements: messages.CreateResourceMsg{
				ProjectName: "project_name_test",
			},
		},
		{
			name:        "Replacement between text",
			data:        "abc{{ProjectName}}def",
			expected:    "abcproject_name_testdef",
			expectError: false,
			replacements: messages.CreateResourceMsg{
				ProjectName: "project_name_test",
			},
		},
		{
			name:        "Replacement with separator",
			data:        "{{ProjectName}}_xyz",
			expected:    "project_name_test_xyz",
			expectError: false,
			replacements: messages.CreateResourceMsg{
				ProjectName: "project_name_test",
			},
		},
		{
			name:        "No replacement",
			data:        "xyz",
			expected:    "xyz",
			expectError: false,
			replacements: messages.CreateResourceMsg{
				ProjectName: "project_name_test",
			},
		},
		{
			name:         "Can't find replacement",
			data:         "{{ProjectName}}",
			expected:     "",
			expectError:  true,
			replacements: messages.CreateResourceMsg{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := renameFile(tc.data, &tc.replacements)
			if tc.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Test %s failed: Expected %s, got %s", tc.name, tc.expected, result)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	replacements := &messages.CreateResourceMsg{
		ProjectName: "project_name",
		AWSRegion:   "eu-central-1",
	}

	tcs := []struct {
		name            string
		templateContent string
		filename        string
		expected        string
		expectError     bool
	}{
		{
			name:            "Basic file",
			templateContent: "abc {{.ProjectName}}",
			filename:        "test_file.txt",
			expected:        "abc project_name",
			expectError:     false,
		},
		{
			name:            "Multiple variables",
			templateContent: "abc {{.ProjectName}} def {{.AWSRegion}}",
			filename:        "test_file.txt",
			expected:        "abc project_name def eu-central-1",
			expectError:     false,
		},
		{
			name: "Multiline file",
			templateContent: `First {{.ProjectName}} line
Second {{.AWSRegion}} line
Third line`,
			filename: "test_file.txt",
			expected: `First project_name line
Second eu-central-1 line
Third line`,
			expectError: false,
		},
		{
			name:            "File creation with missing template variable",
			templateContent: "{{.MissingVariable}}",
			filename:        "test_file.txt",
			expected:        "",
			expectError:     true,
		},
		{
			name:            "Empty template content",
			templateContent: "",
			filename:        "test_file.txt",
			expected:        "",
			expectError:     false,
		},
		{
			name:            "No variables",
			templateContent: "abc",
			filename:        "test_file.txt",
			expected:        "abc",
			expectError:     false,
		},
		{
			name:            "Not existing template",
			templateContent: "abc",
			filename:        "not_existing.txt",
			expected:        "",
			expectError:     true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			fs := fstest.MapFS{
				"test_file.txt": &fstest.MapFile{
					Data: []byte(tc.templateContent),
				},
			}

			destDir := "dest"
			if _, err := os.Stat(destDir); os.IsNotExist(err) {
				err := os.Mkdir(destDir, os.ModePerm)
				if err != nil {
					t.Error("failed to create directory")
				}
			}
			defer os.RemoveAll(destDir)

			destPath := fmt.Sprintf("%s/%s", destDir, tc.filename)

			err := createFile(
				fs,
				tc.filename,
				destPath,
				replacements,
			)

			if tc.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			content, err := os.ReadFile(destPath)
			if err != nil {
				t.Errorf("error reading result file %v", err)
			}

			if string(content) != tc.expected {
				t.Errorf("Test %s failed, expected %s, got %s", tc.name, tc.expected, string(content))
			}
		})
	}
}

func TestCreateDirectory(t *testing.T) {
	tcs := []struct {
		name        string
		path        string
		cleanUpPath string
		expectError bool
	}{
		{
			name:        "Create new directory",
			path:        "testDir1",
			cleanUpPath: "testDir1",
			expectError: false,
		},
		{
			name:        "Create nested directory",
			path:        "testDir2/nestedDir",
			cleanUpPath: "testDir2/",
			expectError: false,
		},
		{
			name:        "Create existing directory",
			path:        "testDir1",
			cleanUpPath: "testDir1",
			expectError: false,
		},
		{
			name:        "Invalid directory path",
			path:        "\x00",
			cleanUpPath: "",
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			defer os.RemoveAll(tc.cleanUpPath)

			err := createDirectory(tc.path)

			if tc.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if _, err := os.Stat(tc.path); os.IsNotExist(err) {
				t.Errorf("directory %s was not created", tc.path)
			}
		})
	}
}

func TestCopyFiles(t *testing.T) {
	tcs := []struct {
		name         string
		fsys         fstest.MapFS
		replacements messages.CreateResourceMsg
		destDir      string
		expected     []string
		expectError  bool
	}{
		{
			name: "Basic copy",
			fsys: fstest.MapFS{
				"src/{{ProjectName}}":                     {Data: []byte("project content")},
				"src/file1.txt":                           {Data: []byte("file 1 content")},
				"src/file2.txt":                           {Data: []byte("file 2 content")},
				"src/nested_src/file1.txt":                {Data: []byte("nested file 1 content")},
				"src/{{ProjectName}}/file2.txt":           {Data: []byte("project file 2 content")},
				"src/{{ProjectName}}/{{ProjectName}}.txt": {Data: []byte("project project content")},
			},
			replacements: messages.CreateResourceMsg{
				ProjectName: "test_project_name",
			},
			destDir: "dest",
			expected: []string{
				"dest/test_project_name",
				"dest/file1.txt",
				"dest/file2.txt",
				"dest/nested_src/file1.txt",
				"dest/test_project_name/file2.txt",
				"dest/test_project_name/test_project_name.txt",
			},
			expectError: false,
		},
		{
			name: "Copy no project name",
			fsys: fstest.MapFS{
				"src/{{ProjectName}}":                     {Data: []byte("project content")},
				"src/file1.txt":                           {Data: []byte("file 1 content")},
				"src/file2.txt":                           {Data: []byte("file 2 content")},
				"src/nested_src/file1.txt":                {Data: []byte("nested file 1 content")},
				"src/{{ProjectName}}/file2.txt":           {Data: []byte("project file 2 content")},
				"src/{{ProjectName}}/{{ProjectName}}.txt": {Data: []byte("project project content")},
			},
			replacements: messages.CreateResourceMsg{},
			destDir:      "dest",
			expectError:  true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := copyFiles(tc.fsys, "src", tc.destDir, &tc.replacements)

			if tc.expectError {
				if err == nil {
					t.Error("expected error got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			for _, path := range tc.expected {
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("expected file not found: %s", path)
				}
			}

			os.RemoveAll(tc.destDir)
		})
	}
}
