package event

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

const (
	IDServerStart    = "server_start"
	IDServerShutdown = "server_shutdown"
	IDPlayerJoin     = "player_join"
	IDPlayerQuit     = "player_quit"
)

type Event interface {
	Name() string
}

type ServerStart struct {
	ProgramName string
}

func (ServerStart) Name() string {
	return IDServerStart
}

type ServerShutdown struct {
	ProgramName string
}

func (ServerShutdown) Name() string {
	return IDServerShutdown
}
