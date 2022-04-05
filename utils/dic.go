package utils

import (
	"io/ioutil"
	"mngr/models"
	"os"
	"path"
	"strings"
)

func CreateDicIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetStreamFolderPath(config *models.Config) (string, error) {
	dir := config.Path.Stream
	CreateDicIfNotExist(dir)

	return dir, nil
}

func GetRecordFolderPath(config *models.Config) (string, error) {
	// create if it doesn't exist
	dir := config.Path.Record
	CreateDicIfNotExist(dir)

	return dir, nil
}

func SetHlsPath(config *models.Config, s *models.StreamModel) {
	s.HlsOutputPath = strings.Replace(s.HlsOutputPath, config.Path.Stream, "", -1)
}

func GetDetectedFolderName() string {
	return "detected"
}

func GetDetectedFolderPath() string {
	return "./static/" + GetDetectedFolderName()
}

func CleanDetectedFolder() error {
	rootPath := GetDetectedFolderPath()
	files, _ := ioutil.ReadDir(rootPath)
	for _, file := range files {
		os.Remove(path.Join(rootPath, file.Name()))
	}

	return nil
}
