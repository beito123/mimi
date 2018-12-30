package pks

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import uuid "github.com/satori/go.uuid"

// IncompatibleProtocol is a error packet
// It notifies that client's protocol is incompatible with server's protocol.
// If a client is received, the connection is closed.
// Server -> Client
type IncompatibleProtocol struct {
	BasePacket

	Protocol byte
}

func (IncompatibleProtocol) ID() byte {
	return IDIncompatibleProtocol
}

func (pk *IncompatibleProtocol) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutByte(pk.Protocol)
	if err != nil {
		return err
	}

	return nil
}

func (pk *IncompatibleProtocol) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.Protocol, err = pk.Byte()
	if err != nil {
		return err
	}

	return nil
}

func (IncompatibleProtocol) New() Packet {
	return new(IncompatibleProtocol)
}

// BadRequest is a error packet
// It notifies that a request received client had problems for some reasons
// If a client is received, the connection is closed.
// Server -> Client
type BadRequest struct {
	BasePacket

	Message string
}

func (BadRequest) ID() byte {
	return IDBadRequest
}

func (pk *BadRequest) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutString(pk.Message)
	if err != nil {
		return err
	}

	return nil
}

func (pk *BadRequest) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.Message, err = pk.String()
	if err != nil {
		return err
	}

	return nil
}

func (BadRequest) New() Packet {
	return new(BadRequest)
}

// ConnectionOne is a first packet sent from server
// It notifies that connected with server
type ConnectionOne struct {
	BasePacket

	UUID uuid.UUID // Server's UUID
	Time int64     // Connected Time format: unix timestamp (second)
}

func (ConnectionOne) ID() byte {
	return IDConnectionOne
}

func (pk *ConnectionOne) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutUUID(pk.UUID)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Time)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ConnectionOne) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.UUID, err = pk.GetUUID()
	if err != nil {
		return err
	}

	pk.Time, err = pk.Long()
	if err != nil {
		return err
	}

	return nil
}

func (ConnectionOne) New() Packet {
	return new(ConnectionOne)
}

// ConnectionRequest is a connection packet sent by client
// Client -> Server
type ConnectionRequest struct {
	BasePacket

	ClientProtocol byte
	ClientUUID     uuid.UUID
}

func (ConnectionRequest) ID() byte {
	return IDConnectionRequest
}

func (pk *ConnectionRequest) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutByte(pk.ClientProtocol)
	if err != nil {
		return err
	}

	err = pk.PutUUID(pk.ClientUUID)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ConnectionRequest) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.ClientProtocol, err = pk.Byte()
	if err != nil {
		return err
	}

	pk.ClientUUID, err = pk.GetUUID()
	if err != nil {
		return err
	}

	return nil
}

func (ConnectionRequest) New() Packet {
	return new(ConnectionRequest)
}

// ConnectionResponse is a packet notifying it is established connection with server
// Client -> Server
type ConnectionResponse struct {
	BasePacket

	Time int64
}

func (ConnectionResponse) ID() byte {
	return IDConnectionResponse
}

func (pk *ConnectionResponse) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Time)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ConnectionResponse) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.Time, err = pk.Long()
	if err != nil {
		return err
	}

	return nil
}

func (ConnectionResponse) New() Packet {
	return new(ConnectionResponse)
}
