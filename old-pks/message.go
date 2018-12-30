package pks

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

// Client -> Server
type RequestProgramList struct {
}

func (RequestProgramList) ID() byte {
	return IDRequestProgramList
}

func (RequestProgramList) New() Packet {
	return new(RequestProgramList)
}

type Program struct {
	Name       string `json:"name"`
	LoaderName string `json:"loader_name"`
}

// Server -> Client
type ResponseProgramList struct {
	Programs []Program `json:"programs"`
}

func (ResponseProgramList) ID() byte {
	return IDResponseProgramList
}

func (ResponseProgramList) New() Packet {
	return new(ResponseProgramList)
}

// Client -> Server
type StartProgram struct {
	ProgramName string `json:"program"`
}

func (StartProgram) ID() byte {
	return IDStartProgram
}

func (StartProgram) New() Packet {
	return new(StartProgram)
}

// Client -> Server
type StopProgram struct {
	ProgramName string `json:"program"`
}

func (StopProgram) ID() byte {
	return IDStopProgram
}

func (StopProgram) New() Packet {
	return new(StopProgram)
}

// Client -> Server
type RestartProgram struct {
	ProgramName string `json:"program"`
}

func (RestartProgram) ID() byte {
	return IDRestartProgram
}

func (RestartProgram) New() Packet {
	return new(RestartProgram)
}

// ProgramStatus is sended back when a server is received Program packets
// Server -> Client
type ProgramStatus struct {
	ProgramName string `json:"program"`
	ConsoleUUID string `json:"console_uuid"`
	Running     bool   `json:"running"`
}

func (ProgramStatus) ID() byte {
	return IDProgramStatus
}

func (ProgramStatus) New() Packet {
	return new(ProgramStatus)
}

type Console struct {
	UUID    string  `json:"uuid"`
	Program Program `json:"program"`
}

// Client -> Server
type RequestConsoleList struct {
}

func (RequestConsoleList) ID() byte {
	return IDRequestConsoleList
}

func (RequestConsoleList) New() Packet {
	return new(RequestConsoleList)
}

// Server -> Client
type ResponseConsoleList struct {
	Consoles []Console `json:"consoles"`
}

func (ResponseConsoleList) ID() byte {
	return IDResponseConsoleList
}

func (ResponseConsoleList) New() Packet {
	return new(ResponseConsoleList)
}

// Client -> Server
type JoinConsole struct {
	UUID string `json:"uuid"`
}

func (JoinConsole) ID() byte {
	return IDJoinConsole
}

func (JoinConsole) New() Packet {
	return new(JoinConsole)
}

// Client -> Server
type QuitConsole struct {
	UUID string `json:"uuid"`
}

func (QuitConsole) ID() byte {
	return IDQuitConsole
}

func (QuitConsole) New() Packet {
	return new(QuitConsole)
}

// Server -> Client
type ConsoleMessage struct {
	Messages []string `json:"messages"` // older sorted
}

func (ConsoleMessage) ID() byte {
	return IDConsoleMessage
}

func (ConsoleMessage) New() Packet {
	return new(ConsoleMessage)
}

// Client -> Server
type SendCommand struct {
}

func (SendCommand) ID() byte {
	return IDSendCommand
}

func (SendCommand) New() Packet {
	return new(SendCommand)
}
