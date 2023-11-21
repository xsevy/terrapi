package templates

import (
	"fmt"
	"os"
	"testing"
	"testing/fstest"

	"github.com/xsevy/terrapi/messages"
)

func TestRenameFile(t *testing.T) {
	replacements := &messages.CreateResourceMsg{
		ProjectName: "project_name_test",
	}

	tcs := []struct {
		name     string
		data     string
		expected string
	}{
		{
			name:     "Basic replacement",
			data:     "{{ProjectName}}",
			expected: "project_name_test",
		},
		{
			name:     "Multiple replacements",
			data:     "{{ProjectName}}{{ProjectName}}",
			expected: "project_name_testproject_name_test",
		},
		{
			name:     "Replacement with text before",
			data:     "xyz{{ProjectName}}",
			expected: "xyzproject_name_test",
		},
		{
			name:     "Replacement with text after",
			data:     "{{ProjectName}}xyz",
			expected: "project_name_testxyz",
		},
		{
			name:     "Replacement between text",
			data:     "abc{{ProjectName}}def",
			expected: "abcproject_name_testdef",
		},
		{
			name:     "Replacement with separator",
			data:     "{{ProjectName}}_xyz",
			expected: "project_name_test_xyz",
		},
		{
			name:     "No replacement",
			data:     "xyz",
			expected: "xyz",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := renameFile(tc.data, replacements)
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
				os.Mkdir(destDir, os.ModePerm)
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
