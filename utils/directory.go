package utils

import "os"

const LiveFolderPath = "static/live"
const RelativeLiveFolderPath = "./" + LiveFolderPath

func CreateDirIfNotExist(dir string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return dir, err
		}
	}
	return dir, nil
}

func GetExecutablePath() (string, error) {
	return "/mnt/super/ionix/node/mngr", nil
	// todo: fix this
	//path, err := os.Executable()
	//if err != nil {
	//	return "", err
	//}
	//return path, nil
}
