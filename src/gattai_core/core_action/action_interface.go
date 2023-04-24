package core_action

type ActionInterface interface {
	RunAction(string) (string,error)
}
