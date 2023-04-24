package core_action

type ActionInterface interface {
	Run(string) (string,error)
}
