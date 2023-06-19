package framework

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

type BaseConfig struct {
	clientName                  string
	appSecret                   string
	kafkaHost                   string
	kafkaTimeout                time.Duration
	kafkaAttempts               int
	scoreStorageApiHost         string
	authorizerHost              string
	redisOptions                *redis.Options
	repeatScoreChangesTimeframe time.Duration
	commitThreshold             int
}

func LoadBaseConfig(envFilename string, clientName string) (BaseConfig, error) {
	if envFilename != "" {
		err := godotenv.Load(envFilename)
		if err != nil {
			return BaseConfig{}, errors.New(fmt.Sprintf("Error loading %s file: %s", envFilename, err))
		}
	}

	kafkaTimeout, err := strconv.Atoi(os.Getenv("KAFKA_TIMEOUT"))
	if kafkaTimeout == 0 || err != nil {
		kafkaTimeout = 10
	}

	kafkaAttempts, err := strconv.Atoi(os.Getenv("KAFKA_ATTEMPTS"))
	if kafkaAttempts <= 0 || err != nil {
		kafkaAttempts = 0
	}

	commitThreshold, err := strconv.Atoi(os.Getenv("COMMIT_THRESHOLD"))
	if commitThreshold <= 0 || err != nil {
		commitThreshold = 0
	}

	repeatScoreChangesTimeframeSeconds, err := strconv.Atoi(os.Getenv("TIMEFRAME_TO_COMBINE_REPEAT_SCORE_CHANGES"))
	if repeatScoreChangesTimeframeSeconds == 0 || err != nil {
		repeatScoreChangesTimeframeSeconds = 600
	}

	config := BaseConfig{
		clientName:                  clientName,
		appSecret:                   os.Getenv("APP_SECRET"),
		kafkaHost:                   os.Getenv("KAFKA_HOST"),
		kafkaTimeout:                time.Second * time.Duration(kafkaTimeout),
		kafkaAttempts:               kafkaAttempts,
		scoreStorageApiHost:         os.Getenv("SCORE_STORAGE_API_HOST"),
		authorizerHost:              os.Getenv("AUTHORIZER_HOST"),
		repeatScoreChangesTimeframe: time.Second * time.Duration(repeatScoreChangesTimeframeSeconds),
		commitThreshold:             commitThreshold,
	}

	if config.appSecret == "" {
		return BaseConfig{}, errors.New("empty APP_SECRET")
	}

	if config.kafkaHost == "" {
		return BaseConfig{}, errors.New("empty KAFKA_HOST")
	}

	if config.scoreStorageApiHost == "" {
		return BaseConfig{}, errors.New("empty SCORE_STORAGE_API_HOST")
	}

	if config.authorizerHost == "" {
		return BaseConfig{}, errors.New("empty AUTHORIZER_HOST")
	}

	config.redisOptions, err = redis.ParseURL(os.Getenv("REDIS_DSN"))

	if err != nil {
		return BaseConfig{}, err
	}

	return config, nil
}
