package framework

type UserLogoutHandlerInterface interface {
	Handle(clientUserId string) error
}
