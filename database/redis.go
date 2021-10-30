package database

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
	"time"
)

var client *redis.Client
var timeClient *redis.Client
var submitClient *redis.Client
var roundsClient *redis.Client

func InitRedis() {
	redisAddr := os.Getenv("REDIS")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	adminPass := os.Getenv("ADMIN_PASS")
	if adminPass == "" {
		adminPass = "admin"
	}
	client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: adminPass,
		DB:       0,
	})
	submitClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: adminPass,
		DB:       1,
	})
	timeClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: adminPass,
		DB:       2,
	})
	roundsClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: adminPass,
		DB:       3,
	})
}

type FlagStruct struct {
	Flag    string
	ID      string
	Team    string
	Service string
}

func (f *FlagStruct) PutFlag() error {
	status := client.HMSet(f.Flag, map[string]interface{}{
		"team":    f.Team,
		"service": f.Service,
	})
	log.Println(status)
	return nil
}

func GetInfo(flag string) ([]interface{}, error) {
	result, err := client.HMGet(flag, "team", "service").Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *FlagStruct) Put() error {
	index := fmt.Sprintf("%s_%s", f.Team, f.Service)
	status := client.HMSet(index, map[string]interface{}{
		f.ID: f.Flag,
	})
	log.Println(status)
	return nil
}

func (f *FlagStruct) GetKeys() (result []string, err error) {
	index := fmt.Sprintf("%s_%s", f.Team, f.Service)
	result, err = client.HKeys(index).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (f *FlagStruct) GetFlag() (value string, err error) {
	index := fmt.Sprintf("%s_%s", f.Team, f.Service)
	value, err = client.HGet(index, f.ID).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func RemoveAllFlags() {
	client.FlushDB()
}

func WriteTime() {
	timeClient.RPush("time", time.Now().Format(time.RFC3339))
}
func GetTime(index int64) (string, error) {
	result, err := timeClient.LIndex("time", index).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
func GetStartTimeStamp() (string, error) {
	return GetTime(0)
}

func GetLastTimeStamp() (string, error) {
	return GetTime(-1)
}

func AddSubmitFlag(flagStruct *FlagStruct) {
	status := submitClient.HMSet(flagStruct.Flag, map[string]interface{}{
		"team":    flagStruct.Team,
		"service": flagStruct.Service,
	})
	log.Println(status)
}

func GetSubmitFlags(flag string) ([]interface{}, error) {
	result, err := submitClient.HMGet(flag, "team", "service").Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetRound() (int, error) {
	round, err := roundsClient.Get("round").Result()
	if err != nil {
		return 0, err
	}
	intRound, err := strconv.Atoi(round)
	if err != nil {
		return 0, err
	}
	return intRound, nil
}
func IncrRound() {
	roundsClient.Incr("round")
}
