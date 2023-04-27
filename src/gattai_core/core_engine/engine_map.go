package core_engine

import (
	"fmt"
	"sync"
)

type FetchEntry struct {
	fetchFn FetchFunc
	result string
	done bool
}

type FetchLookUp struct {
	m sync.Map
}

func MakeFetchLookUp() FetchLookUp {
	return FetchLookUp {}
}

func (lu *FetchLookUp) Set(key string, val FetchFunc) {
    lu.m.Store(key, FetchEntry {
		fetchFn: val,
	})
}

func (lu *FetchLookUp) Get(key string, engine* Engine) (string, error) {
	if val, ok := lu.m.Load(key); ok {
		var result string
		if val.(FetchEntry).done {
			result = val.(FetchEntry).result
		} else {
			action,err := val.(FetchEntry).fetchFn(key,engine)
			if err != nil {
				return "", fmt.Errorf("FetchLookUp fetchFn error: %v", err)
			}
			out_result, err := engine.commandFunc(action)
			if err != nil {
				return "", fmt.Errorf("FetchLookUp commandFunc error: %v", err)
			}
			lu.m.Store(key, FetchEntry {
				fetchFn: val.(FetchEntry).fetchFn,
				result: out_result,
				done: true,
			})
			result = out_result
		}
		return result, nil
    }
    return "", fmt.Errorf("FetchLookUp Get cannot find key error: %v", key)
}
