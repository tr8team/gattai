package core_action

type ActionInterface interface {
	RunAcion(string) (string,error)
}
