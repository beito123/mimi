package pks

import (
	"encoding/json"
	"errors"
)

// Packet is a simple packet interface
type Packet interface {
	ID() byte
	New() Packet
}

func GetPacket(id byte) (Packet, bool) {
	pk, ok := Protocol[id]
	if !ok {
		return nil, false
	}

	return pk.New(), true
}

// Packet Format:
// 1byte: Packet ID
// ?bytes: Packet Data (json format)

// EncodePacket encodes a packet
func EncodePacket(pk Packet) ([]byte, error) {
	data, err := json.Marshal(pk)
	if err != nil {
		return nil, err
	}

	b := []byte{byte(pk.ID())}

	b = append(b, data...)

	return b, nil
}

func DecodePacket(b []byte) (Packet, error) {
	if len(b) == 0 {
		return nil, errors.New("empty data")
	}

	id := b[0]

	pk, ok := GetPacket(id)
	if !ok {
		return nil, errors.New("unknown packet")
	}

	err := json.Unmarshal(b[1:], pk)
	if err != nil {
		return nil, err
	}

	return pk, nil
}
