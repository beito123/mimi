package mimi

import "errors"

const (
	IDConnectionOne = iota
	IDConnectionRequest
	IDConnectionResponse
)

var Packets = map[byte]Packet{
	IDConnectionOne:      &ConnectionOne{},
	IDConnectionRequest:  &ConnectionRequest{},
	IDConnectionResponse: &ConnectionResponse{},
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

	pk, ok := Packets[id]
	if !ok {
		return nil, errors.New("unknown packet")
	}

	err := json.Unmarshal(b[1:], pk)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

type Packet interface {
	ID() byte
	New() Packet
}

// ConnectionOne is a first packet from server
// It notifies that connected with server
type ConnectionOne struct {
	UUID string `json:"uuid"`
	Time int    `json:"time"`
}

func (ConnectionOne) ID() byte {
	return IDConnectionOne
}

func (ConnectionOne) New() Packet {
	return new(ConnectionOne)
}

// ConnectionRequest is a packet
// Client -> Server
type ConnectionRequest struct {
}

func (ConnectionRequest) ID() byte {
	return IDConnectionRequest
}

func (ConnectionRequest) New() Packet {
	return new(ConnectionRequest)
}

// ConnectionResponse is a packet
// Server -> Client
type ConnectionResponse struct {
}

func (ConnectionResponse) ID() byte {
	return IDConnectionResponse
}

func (ConnectionResponse) New() Packet {
	return new(ConnectionResponse)
}
