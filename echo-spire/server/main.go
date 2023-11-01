package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const (
	clientID             = "spiffe://corbalt.com/echo/client"
	spireAgentSocketAddr = "unix:///tmp/spire-agent/public/api.sock"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %v server-address\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	err := run(context.Background(), os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, serverAddr string) error {
	listener, err := spiffetls.ListenWithMode(
		ctx,
		"tcp",
		serverAddr,
		spiffetls.MTLSServerWithSourceOptions(
			tlsconfig.AuthorizeID(spiffeid.RequireFromString(clientID)),
			workloadapi.WithClientOptions(workloadapi.WithAddr(spireAgentSocketAddr)),
		),
	)
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
