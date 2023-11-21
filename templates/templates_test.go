package templates

import (
	"testing"

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
