package utilities

import "os"

func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0755)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
