package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func GetFileName(Dst string) string {
	return path.Base(Dst)
}

func GetFileDir(Dst string) string {
	return filepath.Dir(Dst)
}

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println("Unable to find the directory, creating a new one")
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}
