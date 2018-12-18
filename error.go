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
)

type Error struct {
	ID      int
	Message string
}

var (
	ErrJaga = &Error{
		ID:      ErrIDJaga,
		Message: "jagajaga",
	}
)
