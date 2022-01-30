package reps

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/lithammer/shortuuid/v3"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	MAIN     = 0
	SERVICE  = 1
	SOURCES  = 2
	EVENTBUS = 3
)

func CreateRedisConnection(db int) *redis.Client {
	host := os.Getenv("REDIS_HOST")
	fmt.Println("Redis host: ", host)
	if len(host) == 0 {
		host = "127.0.0.1"
	}
	portStr := os.Getenv("REDIS_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Println("An error occurred while converting Redis port value:" + err.Error())
		port = 6379
	}

	return redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: "", // no password set
		DB:       db, // use default DB
	})
}

func NewId() string {
	return strings.ToLower(shortuuid.New()[:11])
}

func Map(in interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}
