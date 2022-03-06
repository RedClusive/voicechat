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
