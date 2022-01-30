package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

func GetStreamingFolderPath() (string, error) {
	config, err := ConfigRep.GetConfig()
	if err != nil {
		log.Println("GetRecordingFolderPath:", err)
		return "", err
	}
	dir := config.Path.Streaming
	CreateDicIfNotExist(dir)

	return dir, nil
}

func GetRecordingFolderPath() (string, error) {
	config, err := ConfigRep.GetConfig()
	if err != nil {
		log.Println("GetRecordingFolderPath:", err)
		return "", err
	}
	// create if it doesn't exist
	dir := config.Path.Recording
	CreateDicIfNotExist(dir)

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
