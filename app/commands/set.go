package commands

import (
	"strconv"
	"strings"
	"time"
)

type Set struct {
	K    string
	Val  string
	Opts []string
}

func NewSet(k string, v string, opts []string) *Set {
	return &Set{k, v, opts}
}

func (s *Set) GetMetadata() map[string]any {
	meta := make(map[string]any)

	for i, v := range s.Opts {
		if strings.ToLower(v) == "px" {
			dur, err := strconv.Atoi(s.Opts[i+1])
			if err != nil {
				panic(err)
			}
			meta["expires"] = time.Now().Add(time.Duration(dur) * time.Millisecond)
		}
	}
	return meta
}
