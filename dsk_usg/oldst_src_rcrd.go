package dsk_usg

import (
	"io/ioutil"
	"log"
	"mngr/utils"
	"os"
	"path"
	"strconv"
	"time"
)

type OldestSourceRecord struct {
	Found    bool
	Year     int
	YearStr  string
	Month    int
	MonthStr string
	Day      int
	DayStr   string
	Path     string
}

func (o *OldestSourceRecord) Init(fullPath string) {
	findMin := func(p string) int {
		fileNumbers := make([]int, 0)
		allFiles, _ := ioutil.ReadDir(p)
		for _, file := range allFiles {
			if !file.IsDir() {
				continue
			}
			dirNameAsInt, err := strconv.Atoi(file.Name())
			if err == nil {
				fileNumbers = append(fileNumbers, dirNameAsInt)
			}
		}
		if len(fileNumbers) == 0 {
			return -1
		}
		min, _ := utils.MinMax(fileNumbers)
		return min
	}
	o.Found = false
	o.Year = findMin(fullPath)
	if o.Year >= 0 {
		o.YearStr = utils.FixZero(o.Year)
		o.Month = findMin(path.Join(fullPath, o.YearStr))
		if o.Month >= 0 {
			o.MonthStr = utils.FixZero(o.Month)
			o.Day = findMin(path.Join(fullPath, o.YearStr, o.MonthStr))
			if o.Day >= 0 {
				o.DayStr = utils.FixZero(o.Day)
				o.Path = path.Join(fullPath, o.YearStr, o.MonthStr, o.DayStr)
				o.Found = true
			}
		}
	}
}

func (o *OldestSourceRecord) CreateMinTime() time.Time {
	return time.Date(o.Year, time.Month(o.Month), o.Day, 0, 0, 0, 0, time.UTC)
}
func (o *OldestSourceRecord) CreateMaxTime() time.Time {
	return time.Date(o.Year, time.Month(o.Month), o.Day, 23, 59, 59, 0, time.UTC)
}
func (o *OldestSourceRecord) CreateTmpFolderPathName() string {
	return path.Join(o.YearStr, o.MonthStr, o.DayStr)
}
func (o *OldestSourceRecord) DeleteParentDirectoryIfEmpty(fullPath string) {
	monthDirPath := path.Join(fullPath, o.YearStr, o.MonthStr)
	files, err := ioutil.ReadDir(monthDirPath)
	if err != nil {
		log.Println("an error occurred while reading a monthly parent directory, err: " + err.Error())
		return
	}
	if len(files) == 0 {
		err = os.Remove(path.Join(monthDirPath))
		if err != nil {
			log.Println("an error occurred while reading a yearly parent directory, err: " + err.Error())
			return
		}

		yearDirPath := path.Join(fullPath, o.YearStr)
		files, err = ioutil.ReadDir(yearDirPath)
		if err != nil {
			log.Println("an error occurred while reading a yearly parent directory, err: " + err.Error())
			return
		}
		if len(files) == 0 {
			err = os.Remove(path.Join(yearDirPath))
			if err != nil {
				log.Println("an error occurred while reading a yearly parent directory, err: " + err.Error())
				return
			}
		}
	}
}
