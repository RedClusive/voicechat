package voiceio

import (
	"context"
	"io"
)

func FromFile(ctx context.Context, input io.Reader, output chan<- *Chunk) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case output <- &InitChunk:
	}

	for {
		var data = make([]byte, ChunkPayloadSize)
		n, err := io.ReadFull(input, data)
		if err == io.EOF {
			break
		}

		chunk := Chunk{
			Type:    DataType,
			Payload: data[:n],
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case output <- &chunk:
		}

		if err == io.ErrUnexpectedEOF {
			break
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case output <- &EndChunk:
	}

	return nil
}
