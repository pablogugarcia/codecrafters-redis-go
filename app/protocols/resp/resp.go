package resp

import (
	"bytes"
	"errors"
	"fmt"
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
const NullBulkString = "$-1"

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
	result bytes.Buffer
	wr     io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{bytes.Buffer{}, w}
}

func (e *Encoder) WriteString(v []byte) *Encoder {
	toWrite := bytes.Join([][]byte{{String}, v, []byte(CRLF)}, []byte{})
	e.result.Write(toWrite)
	return e
}

func (e *Encoder) WriteBulkString(v []byte) *Encoder {
	if string(v) == "" {
		return e
	}
	length := fmt.Sprintf("%d", len(v))
	fmt.Printf("Len: %s", length)
	toWrite := bytes.Join([][]byte{{BulkString}, []byte(length), []byte(CRLF), v, []byte(CRLF)}, []byte{})
	fmt.Printf("To write: %s", toWrite)
	e.result.Write(toWrite)
	return e
}

func (e *Encoder) Encode() error {
	if e.result.Len() == 0 {
		e.wr.Write([]byte("$-1\r\n"))
		return nil
	}
	_, err := e.wr.Write(e.result.Bytes())
	if err != nil {
		fmt.Printf("Error writting the response: %v \n", err)
		return err
	}

	return nil
}
