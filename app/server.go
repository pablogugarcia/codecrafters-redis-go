package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/protocols/resp"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleCoon(conn)
	}
}

func handleCoon(c net.Conn) error {
	defer c.Close()

	for {
		buf := make([]byte, 128)
		_, err := c.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}

		r, err := resp.New(buf)
		if err != nil {
			return fmt.Errorf("error writting: %v", err)
		}

		if r.Type == resp.Array {
			if strings.ToLower(string(r.Elems[0].Parsed)) == "ping" {

				err := resp.NewEncoder(c).Encode(resp.RESP{
					Type:   resp.String,
					Parsed: []byte("PONG"),
				})
				if err != nil {

					return err
				}
				continue

			}
			if strings.ToLower(string(r.Elems[0].Parsed)) == "echo" {

				err := resp.NewEncoder(c).Encode(resp.RESP{
					Type:   resp.BulkString,
					Parsed: []byte(r.Elems[1].Parsed),
				})
				if err != nil {
					return err
				}
				continue
			}
		}

		if err != nil {
			return fmt.Errorf("error writting: %v", err)
		}
	}

	return nil
}
