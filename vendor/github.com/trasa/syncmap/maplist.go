package syncmap

import "sync"

type MapList struct {
	sync.RWMutex
	innermap map[interface{}][]interface{}
}

func NewMapList() MapList {
	return MapList{
		innermap: make(map[interface{}][]interface{}),
	}
}

func (m *MapList) Add(k interface{}, v interface{}) {
	m.Lock()
	defer m.Unlock()
	m.innermap[k] = append(m.innermap[k], v)
}

func (m *MapList) Remove(k interface{}) {
	m.Lock()
	defer m.Unlock()
	delete(m.innermap, k)
}

func (m *MapList) RemoveItem(k interface{}, v interface{}) {
	m.Lock()
	defer m.Unlock()
	// find index of v in innermap[k]
	for index, candidate := range m.innermap[k] {
		if candidate == v {
			// found it
			m.innermap[k] = append(m.innermap[k][:index], m.innermap[k][index+1:]...)
			return
		}
	}
}

func (m *MapList) Get(k interface{}) []interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.innermap[k]
}
