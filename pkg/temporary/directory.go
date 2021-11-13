package temporary

import (
	"os"
	"path"
)

func CreateTempDir() (string, error) {
	dirName, err := os.MkdirTemp("", "sampledir")
	if err != nil{
		return "", err
	}
	return dirName, nil
}

func WriteFileDataToDir(dirName, filename, data string)(string, error) {
	filepath := path.Join(dirName, filename)
	err := os.WriteFile(filepath, []byte(data), 0644)
	if err != nil{
		return "", err
	}
	return filepath, nil
}