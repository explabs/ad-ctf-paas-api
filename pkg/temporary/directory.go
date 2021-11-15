package temporary

import (
	"os"
	"path"
)

func CreateTempDir(dirName string) (string, error) {
	dirName = "/tmp/" + dirName
	err := os.Mkdir(dirName, 0644)
	if err != nil {
		return "", err
	}
	return dirName, nil
}

func WriteFileDataToDir(dirName, filename, data string) (string, error) {
	filepath := path.Join(dirName, filename)
	err := os.WriteFile(filepath, []byte(data), 0644)
	if err != nil {
		return "", err
	}
	return filepath, nil
}
