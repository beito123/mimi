package mimi

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/beito123/mimi/config"
)

var (
	logger = logrus.New()
)

type Server struct {
	
}

func NewServer(conf config.Config) (*Server, error) {
	return &Server{}, nil
}
