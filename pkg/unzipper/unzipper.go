package unzipper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"github.com/artkescha/grader/online_checker/pkg/kit"
)

type UnZipper interface {
	Unzip(src string, dest string, fileExpansiones []string) error
}

type unzipper struct {
	files []string
}

func New() UnZipper {
	return &unzipper{}
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 100500) to an output directory (parameter 2).
func (uz unzipper) Unzip(src string, dest string, fileExpansiones []string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("open zip archive with path %s", src)
	}

	if err = os.Remove(src); err != nil {
		fmt.Printf("remove zip arch with path %s failed %s", src, err)
	}
	defer func() {
		r.Close()
	}()

	for _, f := range r.File {
		if f.FileInfo().IsDir() && kit.IsAllowedExtension(filepath.Ext(f.FileInfo().Name()), fileExpansiones) {
			//folder is continue
			continue
		}
		// Store filename/path for returning and using later on
		//only fileName
		fpath := filepath.Join(dest, f.FileInfo().Name())

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		uz.files = append(uz.files, fpath)

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
