package mimi

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import (
	"time"
)

const (
	Version         = "1.0.0"
	ProtocolVersion = 1
)

const (
	QueryToken = "token" // ?token=jagajaga

	HandshakeTimeout = 10 * time.Second
)

const (
	MaxReceiveStack = 20
	MaxSendStack    = MaxReceiveStack

	MaxProcessData = 20

	UpdateInterval = 100 * time.Millisecond // 0.1
)
