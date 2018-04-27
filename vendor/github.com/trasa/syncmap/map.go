package syncmap

import "sync"

type Map struct {
	sync.RWMutex
	innermap map[interface{}]interface{}
}

func New() Map {
	return Map{
		innermap: make(map[interface{}]interface{}),
	}
}

func (m *Map) Add(k interface{}, v interface{}) {
	m.Lock()
	defer m.Unlock()
	m.innermap[k] = v
}

func (m *Map) Remove(k interface{}) {
	m.Lock()
	defer m.Unlock()
	delete(m.innermap, k)
}

func (m *Map) Get(k interface{}) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.innermap[k]
}

func (m *Map) Iter(routine func(interface{})) {
	m.RLock()
	defer m.RUnlock()
	for _, v := range m.innermap {
		routine(v)
	}
}
