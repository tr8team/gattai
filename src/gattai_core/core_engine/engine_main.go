package core_engine

import (
	"log"
	"strings"
	"github.com/tr8team/gattai/src/gattai_core/core_action"
)

type CommandFunc func(*core_action.Action) (string,error)

type Engine struct {
	fetchlookUp FetchLookUp
	commandFunc CommandFunc
}

type FetchFunc func(string,*Engine)(*core_action.Action,error)

func MakeEngine(cmdFunc CommandFunc) *Engine {
	return &Engine {
		fetchlookUp: MakeFetchLookUp(),
		commandFunc: cmdFunc,
	}
}

func GoroutineFetch(targetKey string, engine *Engine, output chan string) {
	out_result,err := engine.fetchlookUp.Get(targetKey,engine)
	if err != nil {
		log.Fatalf("GoroutineFetch actionlookUp Get error: %v", err)
	}
	result := strings.TrimSpace(out_result)
	output <- result
}

func (engine *Engine) Store(fetchTarget string, fetchFn FetchFunc) {
	engine.fetchlookUp.Set(fetchTarget,fetchFn)
}

func (engine *Engine) Fetch(fetchTarget string) string{
	// check if result for target already exist
	result := make(chan string)
	go GoroutineFetch(fetchTarget,engine, result)
	return <- result
}
