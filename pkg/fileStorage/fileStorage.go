package fileStorage

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"github.com/artkescha/grader/online_checker/pkg/kit"
)

type FileStorage interface {
	UploadFile(file multipart.File) (string, error)
	DownloadFile(id uint64) ([]byte, error)
}

func New(rootPath string) (*fileStorage, error) {
	if exist := kit.ExistsFolder(rootPath); !exist {
		if err := kit.EnsureDir(rootPath); err != nil {
			fmt.Printf("create path %s failed %s", rootPath, err)
		}
	}
	return &fileStorage{rootPath: rootPath}, nil

}

type fileStorage struct {
	rootPath string
}

func (fs fileStorage) UploadFile(file multipart.File) (string, error) {
	// Create a temporary file within our temp-zip directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("./temp-zip", "upload-*.zip")
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

	path = filepath.Join(fs.rootPath, path)

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
