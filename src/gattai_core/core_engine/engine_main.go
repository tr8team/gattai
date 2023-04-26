package core_engine

import (
	"log"
	"strings"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
)

type CommandFunc func(*core_action.Action,string) (string,error)

type Engine struct {
	actionlookUp core_action.ActionLookUp
	commandFunc CommandFunc
}

type FetchFunc func(*Engine)(string,error)

func MakeEngine(cmdFunc CommandFunc) *Engine {
	return &Engine {
		actionlookUp: core_action.MakeActionLookUp(),
		commandFunc: cmdFunc,
	}
}

func GoroutineFetch(targetKey string, engine *Engine,fetchFn FetchFunc, output chan string) {
	result, ok := engine.actionlookUp.Get(targetKey)
	if !ok {
		out_result, err := fetchFn(engine)
		if err != nil {
			log.Fatalf("GoroutineFetch fetchFn error: %v", err)
		}
		result = strings.TrimSpace(out_result)
		engine.actionlookUp.Set(string(targetKey), result)
	}
	output <- result
}

func (engine *Engine) Fetch(fetchTarget string,fetchFn FetchFunc) string{
	// check if result for target already exist
	result := make(chan string)
	go GoroutineFetch(fetchTarget,engine,fetchFn, result)
	return <- result
}

func (engine *Engine) TriggerCommand(action *core_action.Action,action_name string) (string,error) {
	return engine.commandFunc(action,action_name)
}
