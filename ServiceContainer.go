package framework

import (
	"github.com/kneu-messenger-pigeon/authorizer-client"
	"github.com/kneu-messenger-pigeon/client-framework/delayedDeleter"
	"github.com/kneu-messenger-pigeon/client-framework/delayedDeleter/contracts"
	"github.com/kneu-messenger-pigeon/events"
	"github.com/kneu-messenger-pigeon/score-client"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"io"
	"time"
)

const ScoreChangedEventProcessorCount = 6

type ServiceContainer struct {
	DebugLogger                    *DebugLogger
	UserRepository                 *UserRepository
	UserLogoutHandler              *UserLogoutHandler
	AuthorizerClient               *authorizer.Client
	ScoreClient                    *score.Client
	UserAuthorizedEventProcessor   KafkaConsumerProcessorInterface
	ScoreChangedEventProcessorPool [ScoreChangedEventProcessorCount]KafkaConsumerProcessorInterface
	WelcomeAnonymousDelayedDeleter contracts.DeleterInterface
	Executor                       *Executor
	ClientController               ClientControllerInterface
	UserCountMetricsSyncer         UserCountMetricsSyncerInterface
}

func NewServiceContainer(config BaseConfig, out io.Writer) *ServiceContainer {
	redisClient := redis.NewClient(config.redisOptions)

	kafkaDialer := &kafka.Dialer{
		Timeout:   config.kafkaTimeout,
		DualStack: kafka.DefaultDialer.DualStack,
	}

	container := &ServiceContainer{}

	container.DebugLogger = &DebugLogger{
		out:     out,
		enabled: config.debug,
	}

	container.UserRepository = &UserRepository{
		out:   out,
		redis: redisClient,
	}

	container.UserCountMetricsSyncer = &UserCountMetricsSyncer{
		out:            out,
		userRepository: container.UserRepository,
	}

	container.UserLogoutHandler = &UserLogoutHandler{
		out:    out,
		Client: config.clientName,
		writer: &kafka.Writer{
			Addr:     kafka.TCP(config.kafkaHost),
			Topic:    events.AuthorizedUsersTopic,
			Balancer: &kafka.LeastBytes{},
		},
	}

	container.AuthorizerClient = &authorizer.Client{
		Host:       config.authorizerHost,
		Secret:     config.appSecret,
		ClientName: config.clientName,
	}

	container.ScoreClient = &score.Client{
		Host: config.scoreStorageApiHost,
	}

	container.UserAuthorizedEventProcessor = &KafkaConsumerProcessor{
		out:             out,
		commitThreshold: config.commitThreshold,
		handler: &UserAuthorizedEventHandler{
			out:              out,
			clientName:       config.clientName,
			repository:       container.UserRepository,
			serviceContainer: container,
		},
		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers:     []string{config.kafkaHost},
				GroupID:     config.clientName,
				Topic:       events.AuthorizedUsersTopic,
				MinBytes:    10,
				MaxBytes:    10e3,
				MaxWait:     time.Second,
				MaxAttempts: config.kafkaAttempts,
				Dialer:      kafkaDialer,
			},
		),
	}

	scoreChangesReaderConfig := kafka.ReaderConfig{
		Brokers:     []string{config.kafkaHost},
		GroupID:     config.clientName,
		Topic:       events.ScoresChangesFeedTopic,
		MinBytes:    10,
		MaxBytes:    10e3,
		MaxWait:     time.Second,
		MaxAttempts: config.kafkaAttempts,
		Dialer:      kafkaDialer,
	}

	for i := 0; i < len(container.ScoreChangedEventProcessorPool); i++ {
		container.ScoreChangedEventProcessorPool[i] = &KafkaConsumerProcessor{
			out:             out,
			commitThreshold: config.commitThreshold,
			reader:          kafka.NewReader(scoreChangesReaderConfig),
			handler: &ScoreChangedEventHandler{
				out:                        out,
				serviceContainer:           container,
				debugLogger:                container.DebugLogger,
				repository:                 container.UserRepository,
				scoreClient:                container.ScoreClient,
				waitingForAnotherScoreTime: config.waitingForAnotherScoreTime,
				scoreChangedEventComposer: &ScoreChangeEventComposer{
					out:           out,
					redis:         redisClient,
					storageExpire: config.repeatScoreChangesTimeframe,
				},
				scoreChangedStateStorage: &ScoreChangedStateStorage{
					redis:         redisClient,
					storageExpire: config.repeatScoreChangesTimeframe,
				},
				scoreChangedMessageIdStorage: &ScoreChangedMessageIdStorage{
					out:           out,
					redis:         redisClient,
					storageExpire: config.repeatScoreChangesTimeframe,
				},
			},
		}
	}

	container.WelcomeAnonymousDelayedDeleter = delayedDeleter.NewWelcomeAnonymousMessageDelayedDeleter(
		redisClient, out, "welcome_anonymous_message",
	)

	container.Executor = &Executor{
		out:              out,
		serviceContainer: container,
	}

	return container
}

func (container *ServiceContainer) SetController(controller ClientControllerInterface) {
	container.ClientController = controller
	container.WelcomeAnonymousDelayedDeleter.SetHandler(controller)
}
