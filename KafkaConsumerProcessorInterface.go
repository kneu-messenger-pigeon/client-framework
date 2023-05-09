package framework

type KafkaConsumerProcessorInterface interface {
	ExecutableInterface
	Disable()
}
