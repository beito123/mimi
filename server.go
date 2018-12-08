package mimi

import (
	"gitlab.com/beito123/mimi/config"
)

type Server struct {
}

func NewServer(conf config.Config) (*Server, error) {
	return &Server{}, nil
}
