package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var StartupTime = time.Now()

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

var flagString *string = nil

func ParsePort() int {
	var err error
	defaultPort := 8072
	port := 0
	if flagString == nil {
		lock := make(chan bool, 1)
		lock <- true
		if flagString == nil {
			flagString = flag.String("port", "", "Web Server Port Number")
			flag.Parse()
		}
		<-lock
	}

	ep := os.Getenv("WEBSERVER_HOST")
	if len(ep) > 0 {
		port, err = strconv.Atoi(ep)
		if err != nil {
			port = defaultPort
			log.Println("An error occurred while converting Redis port value from environment variable: " + err.Error())
		}
	} else if len(*flagString) > 0 {
		port, err = strconv.Atoi(*flagString)
		if err != nil {
			port = defaultPort
			log.Println("An error occurred while converting Redis port value from arguments :" + err.Error())
		}
	} else {
		port = defaultPort
	}
	return port
}

func MinMax(array []int) (int, int) {
	if len(array) == 0 {
		return -1, -1
	}
	var max = array[0]
	var min = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func DeepCopy[T any](orig *T) (*T, error) {
	origJSON, err := json.Marshal(orig)
	if err != nil {
		return nil, err
	}

	clone := new(T)
	if err = json.Unmarshal(origJSON, &clone); err != nil {
		return nil, err
	}

	return clone, nil
}
