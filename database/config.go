package database

import (
	"context"
	"github.com/explabs/ad-ctf-paas-api/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
)

var defaultRoundInterval = 30

type PlatformConfig struct {
	Mode          string `bson:"mode"`
	RoundInterval int    `bson:"round_interval"`
}

func convertStrTimeToInt(strTime string) int {
	multiplication := map[string]int{
		"s": 1,
		"m": 60,
		"h": 3600,
	}
	timeValue := strTime[:len(strTime)-1]
	timeIntValue, err := strconv.Atoi(timeValue)
	if err != nil {
		return defaultRoundInterval
	}
	timeSuffix := strTime[len(strTime)-1:]
	m, ok := multiplication[timeSuffix]
	if ok {
		return timeIntValue * m
	}
	return defaultRoundInterval
}

func UploadConfig(configData config.Config) (*mongo.UpdateResult, error) {

	c := PlatformConfig{
		Mode:          configData.Mode,
		RoundInterval: convertStrTimeToInt(configData.RoundInterval),
	}

	filter := bson.M{"name": "config"}
	update := bson.M{
		"$set": c,
	}
	opts := options.Update().SetUpsert(true)
	return configurations.UpdateOne(context.Background(), filter, update, opts)
}

func GetMode() string {
	c := getConfig()
	return c.Mode
}
func GetRoundInterval() int {
	c := getConfig()
	return c.RoundInterval
}

func getConfig() *PlatformConfig {
	var c *PlatformConfig
	filter := bson.M{"name": "config"}
	if err := configurations.FindOne(context.Background(), filter).Decode(&c); err != nil {
		log.Println(err)
		return nil
	}
	return c
}
