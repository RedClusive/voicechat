package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/RedClusive/voicechat/pkg/voiceio"

	"github.com/RedClusive/voicechat/pkg/protocol"
	"golang.org/x/sync/errgroup"
)

type clientConfig struct {
	User       string
	ServerAddr string
	Room       string
}

func runSendPipeline(ctx context.Context, connInput chan<- *voiceio.Chunk) error {
	for {
		var command = ""
		if _, err := fmt.Scanf("%s", &command); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		switch command {
		case "voice":
			var path = ""
			if _, err := fmt.Scanf("%s", &path); err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				log.Println(err)
				continue
			}

			if err = voiceio.FromFile(ctx, file, connInput); err != nil {
				return err
			}
		default:
			log.Println("unknown command")
		}
	}
}

func printUsers(users map[string]bool) {
	var outstr = ""
	for user, ok := range users {
		if outstr != "" {
			outstr += ", "
		}
		if ok {
			outstr += fmt.Sprintf("((( %s )))", user)
		} else {
			outstr += fmt.Sprintf("%s", user)
		}
	}

	fmt.Printf("%s\n", outstr)
}

func runReceivePipeline(ctx context.Context, connOutput <-chan *protocol.Packet) error {
	group, ctx := errgroup.WithContext(ctx)

	var ch = make(chan *voiceio.Chunk)

	group.Go(func() error {
		file, err := os.OpenFile("output", os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		return voiceio.ToFile(ctx, ch, bufio.NewWriter(file))
	})

	group.Go(func() error {
		var users = make(map[string]bool)

		for {
			var packet *protocol.Packet

			select {
			case <-ctx.Done():
				return ctx.Err()
			case packet = <-connOutput:
			}

			switch packet.Header.Type {
			case protocol.InitPacket:
				users[packet.User] = true
				printUsers(users)
			case protocol.DataPacket:
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ch <- &voiceio.Chunk{
					Type:    voiceio.DataType,
					Payload: packet.Payload,
				}:
				}
			case protocol.EndPacket:
				users[packet.User] = false
				printUsers(users)
			case protocol.UserDisconectPacket:
				delete(users, packet.User)
				printUsers(users)
			}
		}
	})

	return group.Wait()
}

func run(ctx context.Context, config *clientConfig) error {
	conn := protocol.TCPConn{
		Addr: config.ServerAddr,
		User: config.User,
		Room: config.Room,
	}

	group, ctx := errgroup.WithContext(ctx)

	connInput := make(chan *voiceio.Chunk, 10000)
	connOutput := make(chan *protocol.Packet, 10000)

	group.Go(func() error {
		err := conn.Run(ctx, connInput, connOutput)
		if err != nil {
			log.Fatalln(err) // have to kill blocked on fmt.Scanf goroutine
		}

		return nil
	})

	group.Go(func() error {
		return runSendPipeline(ctx, connInput)
	})

	group.Go(func() error {
		return runReceivePipeline(ctx, connOutput)
	})

	return group.Wait()
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
