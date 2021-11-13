package temporary

import (
	"os"
	"testing"
)

func TestCreateTempDir(t *testing.T) {
	dirName, err := CreateTempDir()
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dirName)
	if _, statErr := os.Stat(dirName); statErr != nil {
		if os.IsNotExist(err) {
			t.Error("directory does not exist")
		} else {
			t.Error(err)
		}
	}
}

var testFileData = map[string]string{
	"test1": "FirstTestString",
	"test2": "secondTestString",
}

func TestWriteFileDataToDir(t *testing.T) {
	dirName, err := CreateTempDir()
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dirName)
	for filename, data := range testFileData {
		filepath, err := WriteFileDataToDir(dirName, filename, data)
		if err != nil {
			t.Error(err)
		}
		readData, err := os.ReadFile(filepath)
		if err != nil {
			t.Error(err)
		}
		if string(readData) != data {
			t.Errorf("Bad data in file %s != %s", data, string(readData))
		}
	}
}
