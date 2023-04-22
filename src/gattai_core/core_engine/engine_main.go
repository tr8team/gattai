package core_engine

import (
	"log"
	"strings"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
)

type Engine struct {
	actionlookUp core_action.ActionLookUp
}

type FetchFunc func(*Engine)(string,error)

func MakeEngine() *Engine {
	return &Engine {
		actionlookUp: core_action.MakeActionLookUp(),
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
