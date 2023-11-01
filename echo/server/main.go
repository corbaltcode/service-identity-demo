package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %v server-address\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	err := run(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func run(serverAddr string) error {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		log.Printf("accepted: %v\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		shout, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("eof: %v\n", conn.RemoteAddr())
			} else {
				log.Println(err)
			}
			return
		}
		log.Printf("echoed: %v (%v bytes)\n", conn.RemoteAddr(), len(shout))
		fmt.Fprint(conn, shout)
	}
}
