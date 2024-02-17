package resp_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/protocols/resp"
)

func TestEcho(t *testing.T) {
	input := []byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")

	r, err := resp.New(input)
	if err != nil {
		t.Errorf("error creating resp: %v", err)
	}

	if r.Type != resp.Array {
		t.Errorf("expected array type, got %s", string(r.Type))
	}

	if string(r.Elems[0].Parsed) != "ECHO" {
		t.Errorf("expected array type, got %s", string(r.Type))
	}

	if string(r.Elems[1].Parsed) != "hey" {
		t.Errorf("expected array type, got %s", string(r.Type))
	}
}
