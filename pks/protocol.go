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
)

var Protocol = map[byte]Packet{
	IDConnectionOne:             &ConnectionOne{},
	IDConnectionRequest:         &ConnectionRequest{},
	IDConnectionResponse:        &ConnectionResponse{},
	IDIncompatibleProtocol:      &IncompatibleProtocol{},
	IDBadRequest:                &BadRequest{},
	IDDisconnectionNotification: &DisconnectionNotification{},
}
