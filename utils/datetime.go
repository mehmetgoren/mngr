package utils

import (
	"path"
	"strconv"
	"strings"
	"time"
)

var sep = "_"

func StringToTime(dateString string) time.Time {
	splits := strings.Split(dateString, sep)
	l := len(splits)

	values := make([]int, 0)
	for _, split := range splits {
		value, _ := strconv.Atoi(split)
		values = append(values, value)
	}

	year := values[0]
	month := 1
	if l > 1 {
		month = values[1]
	}
	day := 1
	if l > 2 {
		day = values[2]
	}
	hour := 0
	if l > 3 {
		hour = values[3]
	}
	minute := 0
	if l > 4 {
		minute = values[4]
	}
	second := 0
	if l > 5 {
		second = values[5]
	}
	nanoSec := 0
	if l > 6 {
		nanoSec = values[6]
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, nanoSec, time.UTC)
}

func TimeToString(time time.Time, includeNanoSec bool) string {
	arr := make([]string, 0)
	arr = append(arr, strconv.Itoa(time.Year()))
	arr = append(arr, FixZero(int(time.Month())))
	arr = append(arr, FixZero(time.Day()))
	arr = append(arr, FixZero(time.Hour()))
	arr = append(arr, FixZero(time.Minute()))
	arr = append(arr, FixZero(time.Second()))
	if includeNanoSec {
		arr = append(arr, FixZero(time.Nanosecond()))
	}

	return strings.Join(arr, sep)
}

func FixZero(val int) string {
	if val < 10 {
		return "0" + strconv.Itoa(val)
	}
	return strconv.Itoa(val)
}

func FixZeroStr(val string) string {
	if len(val) == 1 {
		val = "0" + val
	}
	return val
}

type DateIndex struct {
	Year  string
	Month string
	Day   string
}

func (d *DateIndex) SetValuesFrom(dateStr string) *DateIndex {
	t := StringToTime(dateStr)
	d.Year = strconv.Itoa(t.Year())
	d.Month = FixZero(int(t.Month()))
	d.Day = FixZero(t.Day())
	return d
}

func (d *DateIndex) GetIndexedPath(rootPath string) string {
	arr := make([]string, 0)
	arr = append(arr, rootPath)
	arr = append(arr, d.Year)
	v, _ := strconv.Atoi(d.Month)
	if v > 0 {
		arr = append(arr, d.Month)
	}
	v, _ = strconv.Atoi(d.Day)
	if v > 0 {
		arr = append(arr, d.Day)
	}
	return path.Join(arr...)
}

func DatetimeNow() string {
	now := time.Now()
	return TimeToString(now, true)
}
