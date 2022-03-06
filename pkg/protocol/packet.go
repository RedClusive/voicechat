package protocol

type PacketType uint8

type Header struct {
	Type       PacketType
	Length     uint32
	LengthUser uint32
}

type Packet struct {
	Header  Header
	Payload []byte
	User    string
}

const HeadeSize = 9

const (
	InitPacket          = 0
	DataPacket          = 1
	EndPacket           = 2
	HandshakePacket     = 3
	UserDisconectPacket = 4
)
