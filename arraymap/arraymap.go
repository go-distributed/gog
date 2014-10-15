package arraymap

import (
	"sync"
)

type ArrayMap struct {
	positions map[interface{}]int
	keys      []interface{}
	values    []interface{}
	RWL       sync.RWMutex
}

func NewArrayMap() *ArrayMap {
	return &ArrayMap{
		positions: make(map[interface{}]int),
		keys:      make([]interface{}, 0),
		values:    make([]interface{}, 0),
	}
}

func (a *ArrayMap) Len() int {
	a.RWL.RLock()
	defer a.RWL.RUnlock()
	return len(a.keys)
}

func (a *ArrayMap) Append(key, value interface{}) {
	a.RWL.Lock()
	defer a.RWL.Unlock()
	if _, existed := a.positions[key]; existed {
		return
	}
	a.keys = append(a.keys, key)
	a.values = append(a.values, value)
	a.positions[key] = len(a.keys) - 1
}

func (a *ArrayMap) GetKeyAt(i int) interface{} {
	a.RWL.RLock()
	defer a.RWL.RUnlock()
	return a.keys[i]
}

func (a *ArrayMap) GetValueAt(i int) interface{} {
	a.RWL.RLock()
	defer a.RWL.RUnlock()
	return a.values[i]
}

func (a *ArrayMap) GetValueOf(key interface{}) interface{} {
	a.RWL.RLock()
	defer a.RWL.RUnlock()
	return a.values[a.positions[key]]
}

func (a *ArrayMap) Has(key interface{}) bool {
	a.RWL.RLock()
	defer a.RWL.RUnlock()
	_, existed := a.positions[key]
	return existed
}

func (a *ArrayMap) RemoveAt(i int) {
	a.RWL.Lock()
	defer a.RWL.Unlock()

	removingKey, lastKey := a.keys[i], a.keys[len(a.keys)-1]
	// Swap the removing item and the last.
	a.keys[i], a.keys[len(a.keys)-1] = a.keys[len(a.keys)-1], a.keys[i]
	a.values[i], a.values[len(a.values)-1] = a.values[len(a.values)-1], a.values[i]

	// Update the position.
	a.positions[lastKey] = i

	// Removing.
	a.keys = a.keys[:len(a.keys)-1]
	a.values = a.values[:len(a.values)-1]
	delete(a.positions, removingKey)
}

func (a *ArrayMap) Remove(key interface{}) {
	if _, exisited := a.positions[key]; exisited {
		a.RemoveAt(a.positions[key])
	}
}
