package voiceio

const (
	ChunkPayloadSize = 1024
)

const (
	InitType = 0
	DataType = 1
	EndType  = 2
)

type ChunkType uint8

type Chunk struct {
	Type    ChunkType
	Payload []byte
}

var InitChunk = Chunk{
	Type:    InitType,
	Payload: nil,
}

var EndChunk = Chunk{
	Type:    EndType,
	Payload: nil,
}
