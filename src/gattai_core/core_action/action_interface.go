package core_action

type ActionInterface interface {
	RunAction(string) (string,error)
}

type Action struct {
	Test ActionInterface
	Exec ActionInterface
}
