package protocol

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/RedClusive/voicechat/pkg/voiceio"
	"golang.org/x/sync/errgroup"
)

func WritePacket(bw *bufio.Writer, packet *Packet) error {
	if err := binary.Write(bw, binary.LittleEndian, packet); err != nil {
		return err
	}
	if err := bw.Flush(); err != nil {
		return err
	}

	return nil
}

func (c *TCPConn) SendLoop(ctx context.Context, input <-chan *voiceio.Chunk, bw *bufio.Writer) error {
	for {
		var chunk *voiceio.Chunk

		select {
		case <-ctx.Done():
			return ctx.Err()
		case chunk = <-input:
		}

		packet := Packet{
			Header: Header{
				Type:       PacketType(chunk.Type),
				Length:     uint32(len(chunk.Payload)),
				LengthUser: uint32(len(c.User)),
			},
			Payload: chunk.Payload,
			User:    c.User,
		}

		if err := WritePacket(bw, &packet); err != nil {
			return err
		}
	}
}

func (c *TCPConn) ReceiveLoop(ctx context.Context, br *bufio.Reader, output chan<- *Packet) error {
	for {
		var headerBytes = make([]byte, HeadeSize)
		if _, err := io.ReadFull(br, headerBytes); err != nil {
			return err
		}

		var header Header
		if err := binary.Read(bytes.NewReader(headerBytes), binary.BigEndian, &header); err != nil {
			return err
		}

		var data = make([]byte, header.Length+header.LengthUser)
		if _, err := io.ReadFull(br, data); err != nil {
			return err
		}

		packet := Packet{
			Header:  header,
			Payload: data[:header.Length],
			User:    string(data[header.Length:]),
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case output <- &packet:
		}
	}
}

type TCPConn struct {
	Addr string
	User string
	Room string
}

func (c *TCPConn) Run(ctx context.Context, input <-chan *voiceio.Chunk, output chan<- *Packet) error {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}

	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(conn)

	if err := WritePacket(bw, &Packet{
		Header: Header{
			Type:       HandshakePacket,
			Length:     uint32(len(c.Room)),
			LengthUser: uint32(len(c.User)),
		},
		Payload: []byte(c.Room),
		User:    c.User,
	}); err != nil {
		return err
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return c.SendLoop(ctx, input, bw)
	})

	group.Go(func() error {
		return c.ReceiveLoop(ctx, br, output)
	})

	return group.Wait()
}
