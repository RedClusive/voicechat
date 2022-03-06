package voiceio

import (
	"context"
	"io"
)

func ToFile(ctx context.Context, input <-chan *Chunk, output io.Writer) error {
	for {
		var chunk *Chunk
		select {
		case <-ctx.Done():
			return ctx.Err()
		case chunk = <-input:
		}

		_, err := output.Write(chunk.Payload)
		if err != nil {
			return err
		}
	}
}
