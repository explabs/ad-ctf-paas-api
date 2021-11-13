package archive

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ArchiveData struct {
	Filename string
	Data string
}

func CompressStrings(data []*ArchiveData){

}

func Compress(inputFilePath, outputFilePath string) (err error) {
	inputFilePath = stripTrailingSlashes(inputFilePath)
	inputFilePath, outputFilePath, err = makeAbsolute(inputFilePath, outputFilePath)
	if err != nil {
		return err
	}
	undoDir, err := mkdirAll(filepath.Dir(outputFilePath), 0755)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			undoDir()
		}
	}()

	err = compress(inputFilePath, outputFilePath, filepath.Dir(inputFilePath))
	if err != nil {
		return err
	}

	return nil
}

// The main interaction with tar and gzip. Creates a archive and recursivly adds all files in the directory.
// The finished archive contains just the directory added, not any parents.
// This is possible by giving the whole path exept the final directory in subPath.
func compress(inPath, outFilePath, subPath string) (err error) {
	files, err := ioutil.ReadDir(inPath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("targz: input directory is empty")
	}

	file, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			os.Remove(outFilePath)
		}
	}()

	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)

	err = writeDirectory(inPath, tarWriter, subPath)
	if err != nil {
		return err
	}

	err = tarWriter.Close()
	if err != nil {
		return err
	}

	err = gzipWriter.Close()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}