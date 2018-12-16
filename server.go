package mimi

import (
	"gitlab.com/beito123/mimi/config"
)

func StartServer(conf config.Config) (*Server, error) {
	ser := &Server{}
	
	return ser, nil
}

type Server struct {
	Config config.Config
}
