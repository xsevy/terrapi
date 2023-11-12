package functions

import (
	"fmt"
	"os"
	"path/filepath"
)

// AppendTextToFile appends a string to a file
func AppendTextToFile(filename string, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)

	return err
}

// CreateEmptyFile creates an empty file at the specified path
func CreateEmptyFile(filePath string) error {
	dir := filepath.Dir(filePath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("unable to create directory %s: %v", dir, err)
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %v", filePath, err)
	}
	defer file.Close()

	return nil
}
