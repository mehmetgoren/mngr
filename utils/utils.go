package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
	"log"
	"mngr/models"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func HandlePanic() {
	if r := recover(); r != nil {
		fmt.Println("RECOVER", r)
		debug.PrintStack()
	}
}

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

func GetFileNameWithoutExtension(fileName string) string {
	fileName = filepath.Base(fileName)
	extension := filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(extension)]
	return fileName
}

var re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

func ParseIp(address string) string {
	return re.FindString(address)
}

func FixZero(val int) string {
	if val < 9 {
		return "0" + strconv.Itoa(val)
	}
	return strconv.Itoa(val)
}

type TimeIndex struct {
	Year  string
	Month string
	Day   string
	Hour  string
}

func (i *TimeIndex) FixZeros() {
	if len(i.Month) == 1 {
		i.Month = "0" + i.Month
	}
	if len(i.Day) == 1 {
		i.Day = "0" + i.Day
	}
	if len(i.Hour) == 1 {
		i.Hour = "0" + i.Hour
	}
}

func (i *TimeIndex) SetValuesFrom(t *time.Time) *TimeIndex {
	i.Year = strconv.Itoa(t.Year())
	i.Month = FixZero(int(t.Month()))
	i.Day = FixZero(t.Day())
	i.Hour = FixZero(t.Hour())
	return i
}

func (i *TimeIndex) GetIndexedPath(rootPath string) string {
	arr := make([]string, 0)
	arr = append(arr, rootPath)
	arr = append(arr, i.Year)
	v, _ := strconv.Atoi(i.Month)
	if v > 0 {
		arr = append(arr, i.Month)
	}
	v, _ = strconv.Atoi(i.Day)
	if v > 0 {
		arr = append(arr, i.Day)
	}
	v, _ = strconv.Atoi(i.Hour)
	arr = append(arr, i.Hour)
	return path.Join(arr...)
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

func SetVideoFilePath(v *models.VideoFile) {
	ti := TimeIndex{Year: v.Year, Month: v.Month, Day: v.Day, Hour: v.Hour}
	ti.FixZeros()
	v.Path = path.Join("/playback", v.SourceId, v.Year, ti.Month, ti.Day, ti.Hour, v.Name)
}

func GetVideoFileAbsolutePath(v *models.VideoFile, root string) string {
	ti := TimeIndex{Year: v.Year, Month: v.Month, Day: v.Day, Hour: v.Hour}
	ti.FixZeros()
	return path.Join(root, v.SourceId, v.Year, ti.Month, ti.Day, ti.Hour, v.Name)
}

func NewId() string {
	return strings.ToLower(shortuuid.New()[:11])
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func IsDirNameValid(fileName string) bool {
	regExpString := "\\/?%*:|\"<>"
	reg, err := regexp.Compile(regExpString)
	if err != nil {
		return false
	}
	return !reg.MatchString(fileName)
}

func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func DecodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func GetQsValue(ctx *gin.Context, key string) string {
	qs := ctx.Request.URL.Query()
	if _, ok := qs[key]; !ok {
		log.Println("invalid qs")
		return ""
	}
	value := qs[key][0]
	if len(value) == 0 {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		log.Println("invalid qs")
		return ""
	}
	return value
}
