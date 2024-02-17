package storage

import "sync"

var DefaultStore = &InMemory{
	data: make(map[string]string),
}

type InMemory struct {
	mut  sync.Mutex
	data map[string]string
}

func (m *InMemory) Set(k, v string) {
	m.mut.Lock()
	defer m.mut.Unlock()

	m.data[k] = v
}

func (m *InMemory) Get(k string) (string, bool) {
	m.mut.Lock()
	defer m.mut.Unlock()

	v, ok := m.data[k]

	return v, ok
}
