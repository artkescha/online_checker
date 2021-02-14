package zipper

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Zipper interface {
	Add(filePaths []string, archName string) error
	Get() (string, error)
}

type zipper struct {
	archName string
	w        *zip.Writer
}

func New(targetPath,filePrefix string) (Zipper, error) {
	tmpFile, err := ioutil.TempFile(targetPath, "*" + filePrefix)
	if err != nil {
		return nil, fmt.Errorf("cannot create file for archive: %s", err)
	}
	return &zipper{archName: filepath.Base(tmpFile.Name()), w: zip.NewWriter(tmpFile)}, nil
}

func (z *zipper) Add(filePaths []string, archName string) error {
	for _, filePath := range filePaths {

		fileToZip, err := os.Open(filePath)
		if err != nil {
			return err
		}
		// Get the file information
		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.Join(archName, info.Name())

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := z.w.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, fileToZip)
		//because in loop
		fileToZip.Close()
	}
	return nil
}

func (z *zipper) Get() (string, error) {
	// Перед загрузкой файла следует закрыть io.Writer
	if z.w != nil {
		if err := z.w.Close(); err != nil {
			fmt.Printf("zip writer close error: %s", err)
		}
	}
	return z.archName, nil
}
