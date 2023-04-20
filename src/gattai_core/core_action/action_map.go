package core_action

import (
	"sync"
)

type ActionLookUp struct {
	m sync.Map
}

func MakeActionLookUp() ActionLookUp {
	return ActionLookUp {}
}

func (lu *ActionLookUp) Set(key string, val string) {
    lu.m.Store(key, val)
}

func (lu *ActionLookUp) Get(key string) (string,bool) {
	if val, ok := lu.m.Load(key); ok {
        return val.(string), true
    }
    return "", false
}
