package pks

import uuid "github.com/satori/go.uuid"

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

// RequestProgramList is a request packet
// If it send, it's sent a ResponseProgramList packet back
// Client -> Server
type RequestProgramList struct {
	BasePacket
}

func (RequestProgramList) ID() byte {
	return IDRequestProgramList
}

func (pk *RequestProgramList) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	return nil
}

func (pk *RequestProgramList) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	return nil
}

func (RequestProgramList) New() Packet {
	return new(RequestProgramList)
}

// ResponseProgramList is a response packet for RequestProgramList packet
// Server -> Client
type ResponseProgramList struct {
	BasePacket

	ProgramsLen byte
	Programs    []*Program
}

func (ResponseProgramList) ID() byte {
	return IDResponseProgramList
}

func (pk *ResponseProgramList) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutByte(byte(len(pk.Programs)))
	if err != nil {
		return err
	}

	for _, p := range pk.Programs {
		err = pk.PutProgram(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pk *ResponseProgramList) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.ProgramsLen, err = pk.Byte()
	if err != nil {
		return err
	}

	for i := 0; i < int(pk.ProgramsLen); i++ {
		p, err := pk.GetProgram()
		if err != nil {
			return err
		}

		pk.Programs = append(pk.Programs, p)
	}

	return nil
}

func (ResponseProgramList) New() Packet {
	return new(ResponseProgramList)
}

// StartProgram is a packet starting a program
// Client -> Server
type StartProgram struct {
	BasePacket

	ProgramName string
}

func (StartProgram) ID() byte {
	return IDStartProgram
}

func (pk *StartProgram) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutString(pk.ProgramName)
	if err != nil {
		return err
	}

	return nil
}

func (pk *StartProgram) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.ProgramName, err = pk.String()
	if err != nil {
		return err
	}

	return nil
}

func (StartProgram) New() Packet {
	return new(StartProgram)
}

// StopProgram is a packet stopping a program
// If Restart is true, it will restart a program
// Client -> Server
type StopProgram struct {
	BasePacket

	ProgramName string
	Restart     bool
}

func (StopProgram) ID() byte {
	return IDStopProgram
}

func (pk *StopProgram) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutString(pk.ProgramName)
	if err != nil {
		return err
	}

	err = pk.PutBool(pk.Restart)
	if err != nil {
		return err
	}

	return nil
}

func (pk *StopProgram) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.ProgramName, err = pk.String()
	if err != nil {
		return err
	}

	pk.Restart, err = pk.Bool()
	if err != nil {
		return err
	}

	return nil
}

func (StopProgram) New() Packet {
	return new(StopProgram)
}

// ProgramStatus is sended back when a server is received Program packets
// Server -> Client
type ProgramStatus struct {
	BasePacket

	ProgramName string
	ConsoleUUID uuid.UUID
	Running     bool
}

func (ProgramStatus) ID() byte {
	return IDProgramStatus
}

func (pk *ProgramStatus) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutString(pk.ProgramName)
	if err != nil {
		return err
	}

	err = pk.PutBool(pk.Running)
	if err != nil {
		return err
	}

	err = pk.PutUUID(pk.ConsoleUUID)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ProgramStatus) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.ProgramName, err = pk.String()
	if err != nil {
		return err
	}

	pk.Running, err = pk.Bool()
	if err != nil {
		return err
	}

	pk.ConsoleUUID, err = pk.GetUUID()
	if err != nil {
		return err
	}

	return nil
}

func (ProgramStatus) New() Packet {
	return new(ProgramStatus)
}

// RequestConsoleList is a request packet
// If it send, it will send ResponseConsoleList packet back by server
// Client -> Server
type RequestConsoleList struct {
	BasePacket
}

func (RequestConsoleList) ID() byte {
	return IDRequestConsoleList
}

func (pk *RequestConsoleList) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	return nil
}

func (pk *RequestConsoleList) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	return nil
}

func (RequestConsoleList) New() Packet {
	return new(RequestConsoleList)
}

// ResponseConsoleList is a response packet for consoles
// This is sent back for RequestConsoleList packet
// Server -> Client
type ResponseConsoleList struct {
	BasePacket

	ConsolesLen byte
	Consoles    []*Console
}

func (ResponseConsoleList) ID() byte {
	return IDResponseConsoleList
}

func (pk *ResponseConsoleList) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutByte(byte(len(pk.Consoles)))
	if err != nil {
		return err
	}

	for _, con := range pk.Consoles {
		err = pk.PutConsole(con)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pk *ResponseConsoleList) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.ConsolesLen, err = pk.Byte()
	if err != nil {
		return err
	}

	for i := 0; i < int(pk.ConsolesLen); i++ {
		con, err := pk.GetConsole()
		if err != nil {
			return err
		}

		pk.Consoles = append(pk.Consoles, con)
	}

	return nil
}

func (ResponseConsoleList) New() Packet {
	return new(ResponseConsoleList)
}

// JoinConsole is a packet joining a console
// Client -> Server
type JoinConsole struct {
	BasePacket

	ConsoleUUID uuid.UUID
}

func (JoinConsole) ID() byte {
	return IDJoinConsole
}

func (pk *JoinConsole) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutUUID(pk.ConsoleUUID)
	if err != nil {
		return err
	}

	return nil
}

func (pk *JoinConsole) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.ConsoleUUID, err = pk.GetUUID()
	if err != nil {
		return err
	}

	return nil
}

func (JoinConsole) New() Packet {
	return new(JoinConsole)
}

// QuitConsole is a packet quiting a console
// Client -> Server
type QuitConsole struct {
	BasePacket

	ConsoleUUID uuid.UUID
}

func (QuitConsole) ID() byte {
	return IDQuitConsole
}

func (pk *QuitConsole) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutUUID(pk.ConsoleUUID)
	if err != nil {
		return err
	}

	return nil
}

func (pk *QuitConsole) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.ConsoleUUID, err = pk.GetUUID()
	if err != nil {
		return err
	}

	return nil
}

func (QuitConsole) New() Packet {
	return new(QuitConsole)
}

// ConsoleMessages is a packet sent a console message by server to a client joing the console
// Server -> Client
type ConsoleMessages struct {
	BasePacket

	MessagesLen byte
	Messages    []string // older sorted
}

func (ConsoleMessages) ID() byte {
	return IDConsoleMessages
}

func (pk *ConsoleMessages) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutByte(byte(len(pk.Messages)))
	if err != nil {
		return err
	}

	for _, msg := range pk.Messages {
		err = pk.PutString(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pk *ConsoleMessages) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.MessagesLen, err = pk.Byte()
	if err != nil {
		return err
	}

	for i := 0; i < int(pk.MessagesLen); i++ {
		msg, err := pk.String()
		if err != nil {
			return err
		}

		pk.Messages = append(pk.Messages, msg)
	}

	return nil
}

func (ConsoleMessages) New() Packet {
	return new(ConsoleMessages)
}

// SendCommands is a packet for sending commands
// Client -> Server
type SendCommands struct {
	BasePacket

	CommandsLen byte
	Commands    []string
}

func (SendCommands) ID() byte {
	return IDSendCommands
}

func (pk *SendCommands) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutByte(byte(len(pk.Commands)))
	if err != nil {
		return err
	}

	for _, c := range pk.Commands {
		err = pk.PutString(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pk *SendCommands) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.CommandsLen, err = pk.Byte()
	if err != nil {
		return err
	}

	for i := 0; i < int(pk.CommandsLen); i++ {
		c, err := pk.String()
		if err != nil {
			return err
		}

		pk.Commands = append(pk.Commands, c)
	}

	return nil
}

func (SendCommands) New() Packet {
	return new(SendCommands)
}
