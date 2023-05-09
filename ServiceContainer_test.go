package framework

import (
	"bytes"
	"github.com/kneu-messenger-pigeon/client-framework/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServiceContainer(t *testing.T) {
	t.Run("succeess", func(t *testing.T) {
		out := &bytes.Buffer{}
		config := BaseConfig{
			clientName:          "test-client",
			appSecret:           "test-secret",
			kafkaHost:           "localhost",
			kafkaTimeout:        0,
			kafkaAttempts:       0,
			scoreStorageApiHost: "localhost:8080",
			authorizerHost:      "localhost:8081",
			redisOptions: &redis.Options{
				Network:    "tcp",
				Addr:       "localhost:6379",
				ClientName: "test",
			},
		}

		serviceContainer := NewServiceContainer(config, out)

		assert.NotEmpty(t, serviceContainer.UserRepository)
		assert.NotEmpty(t, serviceContainer.UserRepository.redis)
		assert.Equal(t, out, serviceContainer.UserRepository.out)

		assert.NotEmpty(t, serviceContainer.UserLogoutHandler)
		assert.NotEmpty(t, serviceContainer.UserLogoutHandler.writer)
		assert.Equal(t, config.clientName, serviceContainer.UserLogoutHandler.Client)
		assert.Equal(t, out, serviceContainer.UserLogoutHandler.out)

		assert.NotNil(t, serviceContainer.AuthorizerClient)
		assert.Equal(t, config.authorizerHost, serviceContainer.AuthorizerClient.Host)
		assert.Equal(t, config.clientName, serviceContainer.AuthorizerClient.ClientName)
		assert.Equal(t, config.appSecret, serviceContainer.AuthorizerClient.Secret)

		assert.NotNil(t, serviceContainer.ScoreClient)
		assert.Equal(t, config.scoreStorageApiHost, serviceContainer.ScoreClient.Host)

		assert.NotNil(t, serviceContainer.UserAuthorizedEventProcessor)
		assert.IsType(t, &KafkaConsumerProcessor{}, serviceContainer.UserAuthorizedEventProcessor)
		userAuthorizedEventProcessor := serviceContainer.UserAuthorizedEventProcessor.(*KafkaConsumerProcessor)
		assert.NotNil(t, userAuthorizedEventProcessor.reader)
		assert.Equal(t, out, userAuthorizedEventProcessor.out)
		assert.False(t, userAuthorizedEventProcessor.disabled)
		assert.NotNil(t, userAuthorizedEventProcessor.handler)
		assert.IsType(t, &UserAuthorizedEventHandler{}, userAuthorizedEventProcessor.handler)

		userAuthorizedEventHandler := userAuthorizedEventProcessor.handler.(*UserAuthorizedEventHandler)
		assert.Equal(t, serviceContainer.UserRepository, userAuthorizedEventHandler.repository)
		assert.Equal(t, serviceContainer, userAuthorizedEventHandler.serviceContainer)
		assert.Equal(t, out, userAuthorizedEventHandler.out)
		assert.Equal(t, config.clientName, userAuthorizedEventHandler.clientName)

		assert.NotNil(t, serviceContainer.ScoreChangedEventProcessor)
		assert.IsType(t, &KafkaConsumerProcessor{}, serviceContainer.ScoreChangedEventProcessor)

		scoreChangedEventProcessor := serviceContainer.ScoreChangedEventProcessor.(*KafkaConsumerProcessor)
		assert.NotNil(t, scoreChangedEventProcessor.reader)
		assert.Equal(t, out, scoreChangedEventProcessor.out)
		assert.False(t, scoreChangedEventProcessor.disabled)
		assert.NotNil(t, scoreChangedEventProcessor.handler)

		assert.IsType(t, &ScoreChangedEventHandler{}, scoreChangedEventProcessor.handler)
		scoreChangedEventHandler := scoreChangedEventProcessor.handler.(*ScoreChangedEventHandler)
		assert.Equal(t, out, scoreChangedEventHandler.out)
		assert.Equal(t, serviceContainer, scoreChangedEventHandler.serviceContainer)

		assert.NotNil(t, serviceContainer.Executor)
		assert.IsType(t, &Executor{}, serviceContainer.Executor)
		assert.Equal(t, serviceContainer, serviceContainer.Executor.serviceContainer)
		assert.Equal(t, out, serviceContainer.Executor.out)

		assert.Nil(t, serviceContainer.ClientController)
	})
}

func TestServiceContainer_SetController(t *testing.T) {
	serviceContainer := &ServiceContainer{}
	controller := mocks.NewClientControllerInterface(t)
	serviceContainer.SetController(controller)
	assert.Equal(t, controller, serviceContainer.ClientController)
}
