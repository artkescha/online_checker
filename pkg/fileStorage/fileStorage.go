package fileStorage

import (
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/kit"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileStorage interface {
	UploadFile(file multipart.File) (string, error)
	DownloadFile(id uint64) ([]byte, error)
}

func New(tempFolder, target string) (*fileStorage, error) {
	if exist := kit.ExistsFolder(target); !exist {
		if err := kit.EnsureDir(target); err != nil {
			fmt.Printf("create path %s failed %s", target, err)
		}
	}
	return &fileStorage{targetPath: target, tempFolder: tempFolder}, nil

}

type fileStorage struct {
	targetPath string
	tempFolder string
}

func (fs fileStorage) UploadFile(file multipart.File) (string, error) {
	// Create a temporary file within our temp-zip directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(fs.tempFolder, "upload-*.zip")
	if err != nil {
		return "", fmt.Errorf("creating temp file failed: %s", err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("reading temp file failed %s", err)
	}
	// write this byte array to our temporary file
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return "", fmt.Errorf("write this byte array to our temporary file fieled: %s", err)
	}
	// return that we have successfully uploaded our file!

	return tempFile.Name(), nil

}

func (fs fileStorage) DownloadFile(path string) ([]byte, error) {

	path = filepath.Join(fs.targetPath, path)

	fileBytes, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("read file with path %s failed %s", path, err)
	}
	//delete temp files
	if err := os.Remove(path); err != nil {
		fmt.Printf("remove zip arch with path %s failed %s", path, err)
	}
	return fileBytes, nil
}
