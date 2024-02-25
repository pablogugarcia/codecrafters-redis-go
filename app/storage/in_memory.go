package storage

import "sync"

var DefaultStore = &InMemory{
	data: make(map[string]stored),
}

type stored struct {
	value    string
	metadata map[string]any
}

type InMemory struct {
	mut  sync.Mutex
	data map[string]stored
}

func (m *InMemory) Set(k string, v string, metadata map[string]any) {
	m.mut.Lock()
	defer m.mut.Unlock()

	m.data[k] = stored{value: v, metadata: metadata}
}

func (m *InMemory) Get(k string) (string, map[string]any, bool) {
	m.mut.Lock()
	defer m.mut.Unlock()

	v, ok := m.data[k]

	return v.value, v.metadata, ok
}
