package functions

import (
	"os"
	"testing"
)

func TestAppendTextToFile(t *testing.T) {
	tcs := []struct {
		name                string
		textToAppend        string
		expectedFileContent string
		existingFileContent string
		fileExists          bool
		errExpected         bool
		fromNextLine        bool
	}{
		{
			name:                "Append to empty file successful",
			fileExists:          true,
			errExpected:         false,
			textToAppend:        "Lorem ipsum",
			expectedFileContent: "Lorem ipsum",
			existingFileContent: "",
			fromNextLine:        false,
		},
		{
			name:                "No text to append",
			fileExists:          true,
			errExpected:         true,
			textToAppend:        "",
			expectedFileContent: "",
			existingFileContent: "",
			fromNextLine:        false,
		},
		{
			name:                "Append to file that doesn't exist",
			fileExists:          false,
			errExpected:         true,
			textToAppend:        "Lorem ipsum",
			expectedFileContent: "",
			existingFileContent: "",
			fromNextLine:        false,
		},
		{
			name:                "Append to file that is not empty",
			fileExists:          true,
			errExpected:         false,
			existingFileContent: "lorem ipsum",
			textToAppend:        "dolor sit amet",
			expectedFileContent: "lorem ipsumdolor sit amet",
			fromNextLine:        false,
		},
		{
			name:                "Append to file that is not empty",
			fileExists:          true,
			errExpected:         false,
			existingFileContent: "lorem ipsum",
			textToAppend:        "dolor sit amet",
			expectedFileContent: `lorem ipsum
dolor sit amet`,
			fromNextLine: true,
		},
	}

	filePath := "/tmp/file.txt"

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.fileExists {
				_ = CreateEmptyFile(filePath)

				if tc.existingFileContent != "" {
					f, _ := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
					defer f.Close()

					_, _ = f.WriteString(tc.existingFileContent)
				}
			}
			defer os.RemoveAll(filePath)

			err := AppendTextToFile(filePath, tc.textToAppend, tc.fromNextLine)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			content, _ := os.ReadFile(filePath)
			if string(content) != tc.expectedFileContent {
				t.Errorf("expected %s, got %s", tc.expectedFileContent, string(content))
			}
		})
	}
}

// func TestCreateEmptyFile(t *testing.T) {
//   filePath := "tmp/"
//   err := CreateEmptyFile(filePath string)
// }
