package pks

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

const (
	IDConnectionOne = iota
	IDConnectionRequest
	IDConnectionResponse
	IDIncompatibleProtocol
	IDBadRequest
	IDDisconnectionNotification
	IDErrorMessage
	IDRequestProgramList
	IDResponseProgramList
	IDStartProgram
	IDStopProgram
	IDProgramStatus
	IDRequestConsoleList
	IDResponseConsoleList
	IDJoinConsole
	IDQuitConsole
	IDConsoleMessages
	IDSendCommands
)

var Protocol = map[byte]Packet{
	IDConnectionOne:             &ConnectionOne{},
	IDConnectionRequest:         &ConnectionRequest{},
	IDConnectionResponse:        &ConnectionResponse{},
	IDIncompatibleProtocol:      &IncompatibleProtocol{},
	IDBadRequest:                &BadRequest{},
	IDDisconnectionNotification: &DisconnectionNotification{},
	IDErrorMessage:              &ErrorMessage{},
	IDRequestProgramList:        &RequestProgramList{},
	IDResponseProgramList:       &ResponseProgramList{},
	IDStartProgram:              &StartProgram{},
	IDStopProgram:               &StopProgram{},
	IDProgramStatus:             &ProgramStatus{},
	IDRequestConsoleList:        &RequestConsoleList{},
	IDResponseConsoleList:       &ResponseConsoleList{},
	IDJoinConsole:               &JoinConsole{},
	IDQuitConsole:               &QuitConsole{},
	IDConsoleMessages:            &ConsoleMessages{},
	IDSendCommands:              &SendCommands{},
}

// GetPacket returns a packet registered by Protocol
func GetPacket(id byte) (Packet, bool) {
	pk, ok := Protocol[id]
	if !ok {
		return nil, false
	}

	return pk.New(), true
}
