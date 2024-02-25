package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/protocols/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
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

				err := resp.NewEncoder(c).WriteString([]byte("PONG")).Encode()
				if err != nil {
					return err
				}

			}
			if strings.ToLower(string(r.Elems[0].Parsed)) == "echo" {

				err := resp.NewEncoder(c).WriteBulkString(r.Elems[1].Parsed).Encode()
				if err != nil {
					return err
				}
				continue
			}
			if strings.ToLower(string(r.Elems[0].Parsed)) == "set" {

				opts := make([]string, 0)
				for _, v := range r.Elems {
					fmt.Println("Elem: " + string(v.Parsed))
					opts = append(opts, string(v.Parsed))
				}
				out := commands.NewSet(string(r.Elems[1].Parsed), string(r.Elems[2].Parsed), opts)
				storage.DefaultStore.Set(string(out.K), out.Val, out.GetMetadata())
				err := resp.NewEncoder(c).WriteString([]byte("OK")).Encode()
				if err != nil {
					return err
				}
				continue
			}

			if strings.ToLower(string(r.Elems[0].Parsed)) == "get" {
				cmd := commands.NewGet(string(r.Elems[1].Parsed), storage.DefaultStore)
				out := cmd.Execute()
				fmt.Printf("Get out: %s\n", out)
				c.Write([]byte(out))
				continue
			}
		}

		if err != nil {
			return fmt.Errorf("error writting: %v", err)
		}
	}

	return nil
}
