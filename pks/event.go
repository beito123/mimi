package pks

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

type Event struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}
