package commands

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/protocols/resp"
)

type Get struct {
	K        string
	value    string
	ok       bool
	metadata map[string]any
}

type Getter interface {
	Get(k string) (string, map[string]any, bool)
}

func NewGet(k string, getter Getter) *Get {
	v, metadata, ok := getter.Get(k)

	return &Get{k, v, ok, metadata}
}

func (g *Get) Execute() string {
	if !g.ok {
		return resp.NullBulkString
	}
	for k := range g.metadata {
		if k == "expires" {
			dur, ok := g.metadata[k].(time.Time)
			if !ok {
				panic("metadata expires must be time.Time")
			}
			if time.Now().After(dur) {
				return resp.NullBulkString
			}
		}
	}

	return g.value
}
