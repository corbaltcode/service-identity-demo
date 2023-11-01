package main

import (
	"bufio"
	"fmt"
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

func run(serverAddress string) error {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	stdinReader := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	for {
		shout, err := stdinReader.ReadString('\n')
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(conn, shout)
		if err != nil {
			return err
		}

		echo, err := connReader.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print(echo)
	}
}
