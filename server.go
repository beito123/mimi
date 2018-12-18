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
	"gitlab.com/beito123/mimi/config"

	"github.com/sirupsen/logrus"
)

// Logger

var logger = logrus.New()

// Server

func StartServer(conf config.Config) (*Server, error) {
	ser := &Server{}

	//

	return ser, nil
}

type Server struct {
	Config config.Config
}
