package pks

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

// Server -> Client
type ErrorMessage struct {
	Error int `json:"error"`
}
