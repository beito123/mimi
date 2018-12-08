package mimi

import "errors"

const (
	IDConnectionOne = iota
	IDConnectionRequest
	IDConnectionResponse
	IDIncompatibleProtocol
	IDBadRequest
	IDDisconnectionNotification
)

var Protocol = map[byte]Packet{
	IDConnectionOne:             &ConnectionOne{},
	IDConnectionRequest:         &ConnectionRequest{},
	IDConnectionResponse:        &ConnectionResponse{},
	IDIncompatibleProtocol:      &IncompatibleProtocol{},
	IDBadRequest:                &BadRequest{},
	IDDisconnectionNotification: &DisconnectionNotification{},
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

type Packet interface {
	ID() byte
	New() Packet
}

// ConnectionOne is a first packet from server
// It notifies that connected with server
type ConnectionOne struct {
	UUID string `json:"uuid"` // Management UUID in server side
	Time int64  `json:"time"` // Connected Time format: unix timestamp (second)
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
	ClientProtocol int    `json:"protocol"`
	ClientUUID     string `json:"cid"`
}

func (ConnectionRequest) ID() byte {
	return IDConnectionRequest
}

func (ConnectionRequest) New() Packet {
	return new(ConnectionRequest)
}

// ConnectionResponse is a packet
// Client -> Server
type ConnectionResponse struct {
	Time int64 `json:"time"`
}

func (ConnectionResponse) ID() byte {
	return IDConnectionResponse
}

func (ConnectionResponse) New() Packet {
	return new(ConnectionResponse)
}

// IncompatibleProtocol is a packet
// If a client is received, the connection is closed.
// Server -> Client
type IncompatibleProtocol struct {
	Protocol int `json:"protocol"`
}

func (IncompatibleProtocol) ID() byte {
	return IDIncompatibleProtocol
}

func (IncompatibleProtocol) New() Packet {
	return new(IncompatibleProtocol)
}

// BadRequest is a packet
// If a client is received, the connection is closed.
// Server -> Client
type BadRequest struct {
	Message string `json:"message"`
}

func (BadRequest) ID() byte {
	return IDBadRequest
}

func (BadRequest) New() Packet {
	return new(BadRequest)
}

// DisconnectionNotification is a packet
// Client -> Server or Server -> Client
type DisconnectionNotification struct {
}

func (DisconnectionNotification) ID() byte {
	return IDDisconnectionNotification
}

func (DisconnectionNotification) New() Packet {
	return new(DisconnectionNotification)
}
