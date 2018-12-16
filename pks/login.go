package pks

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
