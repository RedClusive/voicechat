package main

import (
	"context"
	"log"
	"os"

	"github.com/RedClusive/voicechat/pkg/protocol"
	"golang.org/x/sync/errgroup"
)

type clientConfig struct {
	User       string
	ServerAddr string
	Room       string
}

func run(ctx context.Context, config *clientConfig) error {
	conn := protocol.TCPConn{
		Addr: config.ServerAddr,
		User: config.User,
		Room: config.Room,
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {

	})
}

func main() {
	if len(os.Args) < 4 {
		log.Fatalln("not enough arguments")
		return
	}

	config := clientConfig{
		User:       os.Args[1],
		ServerAddr: os.Args[2],
		Room:       os.Args[3],
	}

	ctx := context.Background()
	err := run(ctx, &config)
	log.Println(err)
}
