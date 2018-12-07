package mimi

import (
	"time"
)

const (
	Version    = "1.0.0"
	APIVersion = 1
)

const (
	QueryToken = "token" // ?token=jagajaga
)

const (
	MaxReceiveStack = 20
	MaxSendStack    = MaxReceiveStack

	MaxProcessData = 20

	UpdateInterval = 100 * time.Millisecond // 0.1
)
