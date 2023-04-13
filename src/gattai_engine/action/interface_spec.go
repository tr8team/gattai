package action

type ActionSpec interface {
	TestAction(string,ActionArgs) (string,error)
	ExecAction(string,ActionArgs) (string,error)
}
