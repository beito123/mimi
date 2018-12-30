package pks

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import (
	"errors"

	uuid "github.com/satori/go.uuid"

	"github.com/beito123/binary"
)

// Packet is a simple packet interface
type Packet interface {

	// ID returns a packet ID
	ID() byte

	// Encode encodes a packet
	Encode() error

	// Decode decodes a packet
	Decode() error

	// Bytes returns encoded bytes
	Bytes() []byte

	// SetBytes sets bytes into buffer
	SetBytes([]byte)

	// New returns new instance of the packet
	New() Packet
}

// BasePacket is a basic implement for Packet
type BasePacket struct {
	Packet

	binary.Stream
}

// Encode encodes a packet
func (bpk *BasePacket) Encode(pk Packet) error {
	err := bpk.PutByte(pk.ID())
	if err != nil {
		return err
	}

	return nil
}

// Decode decodes a packet
func (bpk *BasePacket) Decode(pk Packet) error {
	bpk.Skip(1) // id

	return nil
}

// Bytes returns encoded bytes
func (bpk *BasePacket) Bytes() []byte {
	return bpk.AllBytes()
}

// SetBytes returns encoded bytes
func (bpk *BasePacket) SetBytes(b []byte) {
	bpk.SetBytes(b)
}

// String reads a string from bytes
// Format: 2bytes(bytes len) + xbytes(string)
func (bpk *BasePacket) String() (string, error) {
	ln, err := bpk.Short()
	if err != nil {
		return "", err
	}

	b := bpk.Get(int(ln))
	if len(b) < int(ln) {
		return "", errors.New("bytes isn't not enough")
	}

	return string(b), nil
}

// PutString writes a string to bytes
// Format: 2bytes(len) + xbytes(string)
func (bpk *BasePacket) PutString(str string) error {
	b := []byte(str)

	err := bpk.PutShort(uint16(len(b)))
	if err != nil {
		return err
	}

	err = bpk.Put(b)
	if err != nil {
		return err
	}

	return nil
}

// UUID reads a uuid from buffer's bytes
func (bpk *BasePacket) GetUUID() (uuid.UUID, error) {
	b := bpk.Get(uuid.Size)

	uid, err := uuid.FromBytes(b)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

// PutUUID writes a uuid to buffer
func (bpk *BasePacket) PutUUID(uid uuid.UUID) error {
	err := bpk.Put(uid.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// GetProgram reads a program from buffer
func (bpk *BasePacket) GetProgram() (p *Program, err error) {
	p = &Program{}

	p.Name, err = bpk.String()
	if err != nil {
		return nil, err
	}

	p.LoaderName, err = bpk.String()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// PutProgram writes a program to buffer
func (bpk *BasePacket) PutProgram(p *Program) error {
	err := bpk.PutString(p.Name)
	if err != nil {
		return err
	}

	err = bpk.PutString(p.LoaderName)
	if err != nil {
		return err
	}

	return nil
}

// GetConsole reads a console from buffer
func (bpk *BasePacket) GetConsole() (con *Console, err error) {
	con = &Console{}

	con.UUID, err = bpk.GetUUID()
	if err != nil {
		return nil, err
	}

	con.Program, err = bpk.GetProgram()
	if err != nil {
		return nil, err
	}

	return con, nil
}

// PutConsole writes a console to buffer
func (bpk *BasePacket) PutConsole(con *Console) error {
	err := bpk.PutUUID(con.UUID)
	if err != nil {
		return err
	}

	err = bpk.PutProgram(con.Program)
	if err != nil {
		return err
	}

	return nil
}
