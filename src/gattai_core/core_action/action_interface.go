package core_action

type ActionInterface interface {
	RunAction(string) (string,error)
}

type Action struct {
	Name string
	Test ActionInterface
	Exec ActionInterface
}
