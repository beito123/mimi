package mimi

import (
	"errors"
	"net"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func IP(addr string) (net.IP, error) {
	ipstr, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, errors.New("couldn't split address")
	}

	ip := net.ParseIP(ipstr)
	if ip == nil {
		return nil, errors.New("couldn't parse ip")
	}

	return ip, nil
}
