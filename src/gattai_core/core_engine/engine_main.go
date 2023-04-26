package core_engine

import (
	"log"
	"strings"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
)

type CommandFunc func(*core_action.Action) (string,error)

type Engine struct {
	actionlookUp core_action.ActionLookUp
	commandFunc CommandFunc
}

type FetchFunc func(*Engine)(*core_action.Action,error)

func MakeEngine(cmdFunc CommandFunc) *Engine {
	return &Engine {
		actionlookUp: core_action.MakeActionLookUp(),
		commandFunc: cmdFunc,
	}
}

func GoroutineFetch(targetKey string, engine *Engine,fetchFn FetchFunc, output chan string) {
	result, ok := engine.actionlookUp.Get(targetKey)
	if !ok {
		action, err := fetchFn(engine)
		if err != nil {
			log.Fatalf("GoroutineFetch fetchFn error: %v", err)
		}
		out_result, err := engine.commandFunc(action)
		if err != nil {
			log.Fatalf("GoroutineFetch commandFunc error: %v", err)
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
