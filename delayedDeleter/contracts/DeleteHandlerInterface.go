package contracts

type DeleteHandlerInterface interface {
	HandleDeleteTask(task *DeleteTask) error
}
