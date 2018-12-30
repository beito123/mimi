package pks

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

// DisconnectionNotification is a packet
// Client -> Server or Server -> Client
type DisconnectionNotification struct {
	BasePacket
}

func (DisconnectionNotification) ID() byte {
	return IDDisconnectionNotification
}

func (pk *DisconnectionNotification) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	return nil
}

func (pk *DisconnectionNotification) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	return nil
}

func (DisconnectionNotification) New() Packet {
	return new(DisconnectionNotification)
}

// ErrorMessage is a packet notifying a server is happened errors
// Server -> Client
type ErrorMessage struct {
	BasePacket

	Error int `json:"error"`
}

func (ErrorMessage) ID() byte {
	return IDErrorMessage
}

func (pk *ErrorMessage) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutInt(int32(pk.Error))
	if err != nil {
		return err
	}

	return nil
}

func (pk *ErrorMessage) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	e, err := pk.Int()
	if err != nil {
		return err
	}

	pk.Error = int(e)

	return nil
}

func (ErrorMessage) New() Packet {
	return new(ErrorMessage)
}
