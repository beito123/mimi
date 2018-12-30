package pks

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import uuid "github.com/satori/go.uuid"

// Program is data for program
type Program struct {
	Name       string
	LoaderName string
}

// Console is data for console
type Console struct {
	UUID    uuid.UUID
	Program *Program
}
