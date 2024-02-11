package utils

import (
	"os"
	"path/filepath"
)

func GetRootPath(path string) (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return filepath.Join(wd, "../../", path), nil
}
