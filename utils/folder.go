package utils

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func GetStreamingFolderPath() (string, error) {
	folderPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(folderPath, "static/live"), nil
}

func CreateStreamingFolderIfNotExist(id string) (string, error) {
	folderPath, _ := GetStreamingFolderPath()
	dir := path.Join(folderPath, id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return dir, err
		}
	}

	return dir, nil
}

func GetRecordingFolderPath() (string, error) {
	config, err := ConfigRep.GetConfig()
	if err != nil {
		log.Println("GetRecordingFolderPath:", err)
		return "", err
	}

	return config.Recording.FolderPath, nil // "/mnt/sde1/playback", nil
}

func CreateRecordingFolderIfNotExist(id string) (string, error) {
	folderPath, _ := GetRecordingFolderPath()
	dir := path.Join(folderPath, id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return dir, err
		}
	}

	return dir, nil
}

func FromDateToString(t time.Time) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t.Year()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(int(t.Month())))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Day()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Hour()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Minute()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Second()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Nanosecond()))

	return sb.String()
}

//func RemovePrevStreamFolders() {
//	files, _ := ioutil.ReadDir(RelativeLiveFolderPath)
//	for _, file := range files {
//		if !file.IsDir() {
//			continue
//		}
//		folderPath := RelativeLiveFolderPath + "/" + file.Name()
//		os.RemoveAll(folderPath)
//	}
//}
//func ParseVideoFileName(fileName string) time.Time {
//	splits := strings.Split(fileName, ".")
//	splits = strings.Split(splits[0], "_")
//	year, _ := strconv.Atoi(splits[0])
//	month, _ := strconv.Atoi(splits[1])
//	day, _ := strconv.Atoi(splits[2])
//	hour, _ := strconv.Atoi(splits[3])
//	minute, _ := strconv.Atoi(splits[4])
//	second, _ := strconv.Atoi(splits[5])
//
//	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
//}
