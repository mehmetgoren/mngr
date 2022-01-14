package utils

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const LiveFolderPath = "static/live"
const RelativeLiveFolderPath = "./" + LiveFolderPath
const PlaybackFolderPath = "static/playback"
const RelativePlaybackFolderPath = "./" + PlaybackFolderPath

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
