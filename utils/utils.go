package utils

import "os"

// DirectoryExists check whether a directory exists.
func DirectoryExists(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
