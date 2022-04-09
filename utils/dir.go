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
	dir := path.Join(config.General.RootFolderPath, "stream")
	CreateDicIfNotExist(dir)

	return dir, nil
}

func GetRecordFolderPath(config *models.Config) (string, error) {
	dir := path.Join(config.General.RootFolderPath, "record")
	CreateDicIfNotExist(dir)

	return dir, nil
}

func GetOdFolder(config *models.Config) (string, error) {
	dir := path.Join(config.General.RootFolderPath, GetDetectedFolderName())
	CreateDicIfNotExist(dir)

	return dir, nil
}

func SetHlsPath(config *models.Config, s *models.StreamModel) {
	s.HlsOutputPath = strings.Replace(s.HlsOutputPath, path.Join(config.General.RootFolderPath, "stream"), "", -1)
}

func GetDetectedFolderName() string {
	return "od"
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
