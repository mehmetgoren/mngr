package utils

import (
	"log"
	"mngr/models"
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

func GetStreamFolderPath() (string, error) {
	config, err := ConfigRep.GetConfig()
	if err != nil {
		log.Println("GetRecordFolderPath:", err)
		return "", err
	}
	dir := config.Path.Stream
	CreateDicIfNotExist(dir)

	return dir, nil
}

func GetRecordFolderPath() (string, error) {
	config, err := ConfigRep.GetConfig()
	if err != nil {
		log.Println("GetRecordFolderPath:", err)
		return "", err
	}
	// create if it doesn't exist
	dir := config.Path.Record
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

func SetHlsPath(config *models.Config, s *models.StreamModel) {
	s.HlsOutputPath = strings.Replace(s.HlsOutputPath, config.Path.Stream, "", -1)
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
