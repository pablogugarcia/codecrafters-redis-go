package resp

import (
	"errors"
	"io"
	"log"
	"strconv"
)

// RESP is a binary protocol that uses control sequences encoded in standard ASCII.
type RESP struct {
	Elems  []*RESP
	Type   Type
	Raw    []byte
	Parsed []byte
	Count  int
}

// Type is the first byte in an RESP-serialized payload.
type Type byte

const (
	// Simple
	Integer = ':'
	String  = '+'
	Error   = '-'
	// Aggregates or bulk
	Array      = '*'
	BulkString = '$'
)

const CRLF = "\r\n"

func New(b []byte) (*RESP, error) {
	return parse(b)
}

func parse(d []byte) (*RESP, error) {
	t := Type(d[0])

	if t == BulkString {
		log.Println("Parsing Bulkstring")
		i := 1

		for ; i < len(d); i++ {
			if d[i] == '\r' {
				break
			}
		}

		length, err := strconv.Atoi(string(d[1:i]))
		if err != nil {
			return nil, err
		}

		return &RESP{
			Type:   BulkString,
			Count:  length,
			Raw:    d[:i+len(CRLF)+length+len(CRLF)],
			Elems:  nil,
			Parsed: d[i+len(CRLF) : i+len(CRLF)+length],
		}, nil
	}

	if t == Array {
		log.Println("Parsing Array")

		i := 1

		for ; i < len(d); i++ {
			if d[i] == '\r' {
				break
			}
		}

		log.Printf("Converting : %s", string(d[1:i]))

		count, err := strconv.Atoi(string(d[1:i]))
		if err != nil {
			return nil, err
		}

		r := &RESP{
			Count: count,
			Type:  Array,
			Raw:   d[4:],
		}
		ii := 0
		j := 0

		for ; j < len(r.Raw); j++ {
			if r.Raw[j] == '\n' {
				continue
			}
			// end of part
			if r.Raw[j] == '\r' {
				resp, err := parse(r.Raw[ii:])
				if err != nil {
					return nil, err
				}
				r.Elems = append(r.Elems, resp)
				// Fix: this 2 must to change in a dynamic value for cases with two more digits. (10, 11, 12, ...)
				length, _ := strconv.Atoi(string(r.Raw[ii+1 : ii+2]))
				ii = j + len(CRLF) + length + len(CRLF)
				j = ii
			}
		}

		return r, nil
	}

	return nil, errors.New("TODO: missing data type parser")
}

type Encoder struct {
	wr io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(v RESP) error {
	r := make([]byte, v.Count)
	if v.Type == BulkString {
		length := len([]byte(v.Parsed))
		r = append(r, BulkString)
		r = append(r, byte(length))
		r = append(r, []byte(v.Parsed)...)
	}

	_, err := e.wr.Write(r)
	if err != nil {
		return err
	}

	return nil
}
