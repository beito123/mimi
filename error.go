package mimi

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

const (
	ErrIDJaga = iota
	ErrIDInternalError
	ErrIDProgramNotFound
	ErrIDProgramAlreadyRunning
	ErrIDConsoleNotFound
	ErrIDConsoleAlreadyClosed
	ErrIDSessionNotJoinedConsole
)

type ErrorMessage struct {
	ID      int
	Message string
}

var (
	ErrJaga = &ErrorMessage{
		ID:      ErrIDJaga,
		Message: "Jagajaga",
	}
	ErrInternalError = &ErrorMessage{
		ID:      ErrIDInternalError,
		Message: "Internal error",
	}
	ErrProgramNotFound = &ErrorMessage{
		ID:      ErrIDProgramNotFound,
		Message: "A program is not found",
	}
	ErrProgramAlreadyRunning = &ErrorMessage{
		ID:      ErrIDProgramAlreadyRunning,
		Message: "A program is already running",
	}
	ErrConsoleNotFound = &ErrorMessage{
		ID:      ErrIDConsoleNotFound,
		Message: "A console is not found",
	}
	ErrConsoleAlreadyClosed = &ErrorMessage{
		ID:      ErrIDConsoleAlreadyClosed,
		Message: "A console is already closed",
	}
	ErrSessionNotJoined = &ErrorMessage{
		ID:      ErrIDSessionNotJoinedConsole,
		Message: "A session doesn't join a console",
	}
)
