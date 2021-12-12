package kit

import (
	"os"
	"path/filepath"
)

// Exists reports whether the named file or directory existsFolder.
func ExistsFolder(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func EnsureDir(fileName string) error {
	if _, serr := os.Stat(fileName); serr != nil {
		merr := os.MkdirAll(fileName, os.ModePerm)
		if merr != nil {
			return merr
		}
	}
	return nil
}

func FilePathWalkDir(root string, expansion []string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && len(expansion) > 0 && IsAllowedExtension(filepath.Ext(info.Name()), expansion) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func IsAllowedExtension(extension string, allowedExtensiones []string) bool {
	for _, allowedExtension := range allowedExtensiones {
		if allowedExtension != extension {
			return true
		}
	}
	return false
}
