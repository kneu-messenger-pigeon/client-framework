package framework

import (
	"bytes"
	"github.com/kneu-messenger-pigeon/client-framework/delayedDeleter"
	"github.com/kneu-messenger-pigeon/client-framework/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewServiceContainer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		out := &bytes.Buffer{}
		config := BaseConfig{
			clientName:                  "test-client",
			appSecret:                   "test-secret",
			kafkaHost:                   "localhost",
			kafkaTimeout:                0,
			kafkaAttempts:               0,
			scoreStorageApiHost:         "localhost:8080",
			authorizerHost:              "localhost:8081",
			repeatScoreChangesTimeframe: time.Second * 1533,
			commitThreshold:             953,
			redisOptions: &redis.Options{
				Network:    "tcp",
				Addr:       "localhost:6379",
				ClientName: "test",
			},
			waitingForAnotherScoreTime: time.Second * 3,
		}

		serviceContainer := NewServiceContainer(config, out)

		assert.NotNil(t, serviceContainer.DebugLogger)
		assert.IsType(t, &DebugLogger{}, serviceContainer.DebugLogger)
		assert.Equal(t, out, serviceContainer.DebugLogger.out)
		assert.Equal(t, config.debug, serviceContainer.DebugLogger.enabled)

		assert.NotEmpty(t, serviceContainer.UserRepository)
		assert.NotEmpty(t, serviceContainer.UserRepository.redis)
		assert.Equal(t, out, serviceContainer.UserRepository.out)

		assert.NotEmpty(t, serviceContainer.UserCountMetricsSyncer)
		assert.IsType(t, &UserCountMetricsSyncer{}, serviceContainer.UserCountMetricsSyncer)

		userCountMetricsSyncer := serviceContainer.UserCountMetricsSyncer.(*UserCountMetricsSyncer)
		assert.Equal(t, serviceContainer.UserRepository, userCountMetricsSyncer.userRepository)
		assert.Equal(t, out, userCountMetricsSyncer.out)

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
		assert.Equal(t, config.commitThreshold, userAuthorizedEventProcessor.commitThreshold)
		assert.False(t, userAuthorizedEventProcessor.disabled)
		assert.NotNil(t, userAuthorizedEventProcessor.handler)
		assert.IsType(t, &UserAuthorizedEventHandler{}, userAuthorizedEventProcessor.handler)

		userAuthorizedEventHandler := userAuthorizedEventProcessor.handler.(*UserAuthorizedEventHandler)
		assert.Equal(t, serviceContainer.UserRepository, userAuthorizedEventHandler.repository)
		assert.Equal(t, serviceContainer, userAuthorizedEventHandler.serviceContainer)
		assert.Equal(t, out, userAuthorizedEventHandler.out)
		assert.Equal(t, config.clientName, userAuthorizedEventHandler.clientName)

		assert.NotNil(t, serviceContainer.ScoreChangedEventProcessorPool)
		for i := 0; i < ScoreChangedEventProcessorCount; i++ {
			assert.IsType(t, &KafkaConsumerProcessor{}, serviceContainer.ScoreChangedEventProcessorPool[i])
			scoreChangedEventProcessor := serviceContainer.ScoreChangedEventProcessorPool[i].(*KafkaConsumerProcessor)
			assert.NotNil(t, scoreChangedEventProcessor.reader)
			assert.Equal(t, out, scoreChangedEventProcessor.out)
			assert.Equal(t, config.commitThreshold, scoreChangedEventProcessor.commitThreshold)
			assert.False(t, scoreChangedEventProcessor.disabled)
			assert.NotNil(t, scoreChangedEventProcessor.handler)

			assert.IsType(t, &ScoreChangedEventHandler{}, scoreChangedEventProcessor.handler)
			scoreChangedEventHandler := scoreChangedEventProcessor.handler.(*ScoreChangedEventHandler)
			assert.Equal(t, serviceContainer.DebugLogger, scoreChangedEventHandler.debugLogger)
			assert.Equal(t, serviceContainer.UserRepository, scoreChangedEventHandler.repository)
			assert.Equal(t, out, scoreChangedEventHandler.out)
			assert.Equal(t, serviceContainer, scoreChangedEventHandler.serviceContainer)
			assert.Equal(t, config.waitingForAnotherScoreTime, scoreChangedEventHandler.waitingForAnotherScoreTime)

			assert.IsType(t, &ScoreChangeEventComposer{}, scoreChangedEventHandler.scoreChangedEventComposer)
			scoreChangedEventComposer := scoreChangedEventHandler.scoreChangedEventComposer.(*ScoreChangeEventComposer)
			assert.Equal(t, config.repeatScoreChangesTimeframe, scoreChangedEventComposer.storageExpire)
			assert.NotEmpty(t, scoreChangedEventComposer.redis)
			assert.Equal(t, out, scoreChangedEventComposer.out)

			assert.IsType(t, &ScoreChangedStateStorage{}, scoreChangedEventHandler.scoreChangedStateStorage)
			scoreChangedStateStorage := scoreChangedEventHandler.scoreChangedMessageIdStorage.(*ScoreChangedMessageIdStorage)
			assert.Equal(t, config.repeatScoreChangesTimeframe, scoreChangedStateStorage.storageExpire)
			assert.NotEmpty(t, scoreChangedStateStorage.redis)

			assert.IsType(t, &ScoreChangedMessageIdStorage{}, scoreChangedEventHandler.scoreChangedMessageIdStorage)
			scoreChangedMessageIdStorage := scoreChangedEventHandler.scoreChangedMessageIdStorage.(*ScoreChangedMessageIdStorage)
			assert.Equal(t, config.repeatScoreChangesTimeframe, scoreChangedMessageIdStorage.storageExpire)
			assert.NotEmpty(t, scoreChangedMessageIdStorage.redis)
			assert.Equal(t, out, scoreChangedMessageIdStorage.out)
		}

		assert.NotNil(t, serviceContainer.WelcomeAnonymousDelayedDeleter)
		assert.IsType(t, &delayedDeleter.Deleter{}, serviceContainer.WelcomeAnonymousDelayedDeleter)

		assert.NotNil(t, serviceContainer.Executor)
		assert.IsType(t, &Executor{}, serviceContainer.Executor)
		assert.Equal(t, serviceContainer, serviceContainer.Executor.serviceContainer)
		assert.Equal(t, out, serviceContainer.Executor.out)

		assert.Nil(t, serviceContainer.ClientController)
	})
}

func TestServiceContainer_SetController(t *testing.T) {
	controller := mocks.NewClientControllerInterface(t)

	deleter := mocks.NewDeleterInterface(t)
	deleter.On("SetHandler", controller).Return()

	serviceContainer := &ServiceContainer{
		WelcomeAnonymousDelayedDeleter: deleter,
	}

	serviceContainer.SetController(controller)
	assert.Equal(t, controller, serviceContainer.ClientController)
}
