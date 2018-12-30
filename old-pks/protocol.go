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
	IDRestartProgram
	IDProgramStatus
	IDRequestConsoleList
	IDResponseConsoleList
	IDJoinConsole
	IDQuitConsole
	IDConsoleMessage
	IDSendCommand
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
	IDRestartProgram:            &RestartProgram{},
	IDProgramStatus:             &ProgramStatus{},
	IDRequestConsoleList:        &RequestConsoleList{},
	IDResponseConsoleList:       &ResponseConsoleList{},
	IDJoinConsole:               &JoinConsole{},
	IDQuitConsole:               &QuitConsole{},
	IDConsoleMessage:            &ConsoleMessage{},
	IDSendCommand:               &SendCommand{},
}
