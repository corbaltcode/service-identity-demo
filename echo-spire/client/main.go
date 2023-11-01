package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const (
	serverID             = "spiffe://corbalt.com/echo/server"
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
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	conn, err := spiffetls.DialWithMode(
		ctx,
		"tcp",
		serverAddr,
		spiffetls.MTLSClientWithSourceOptions(
			tlsconfig.AuthorizeID(spiffeid.RequireFromString(serverID)),
			workloadapi.WithClientOptions(workloadapi.WithAddr(spireAgentSocketAddr)),
		),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	for {
		shout, err := reader.ReadString('\n')
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
