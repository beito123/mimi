package mimi

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

type ConsoleManager struct {
	Consoles map[string]*Console
}

func StartConsole(loader Loader) *Console {
	con := &Console{}

	con.Cmder = &Cmder{
		WorkingDir: loader.Path(),
	}

	program, args := loader.Cmd()

	con.Cmder.Start(program, args...)

	return con
}

type Console struct {
	Cmder *Cmder

	closed bool
}

func (con *Console) Close() {
	if con.closed {
		return
	}

	con.closed = true

	con.Cmder.Close()
}

func (con *Console) Line() (string, bool) {
	return con.Cmder.Line()
}

func (con *Console) SendCommand(cmd string) {
	con.Cmder.Send(cmd)
}
