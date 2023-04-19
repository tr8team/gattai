package action

type ActionInterface interface {
	TestAction(string) (string,error)
	ExecAction(string) (string,error)
}
