package util

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import (
	"errors"
	"math"
	"net"
	"os"
	"path/filepath"
	"runtime"
)

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

// File

func GetDir(path string) string {
	return filepath.Dir(filepath.Clean(path))
}

func ExistFile(file string) bool {
	f, err := os.Stat(file)
	return err == nil && !f.IsDir()
}

func ExistDir(dir string) bool {
	f, err := os.Stat(dir)

	return err == nil && f.IsDir()
}

// Runtime

func IsWin() bool {
	return runtime.GOOS == "windows"
}

// Math

func MinInt(a int, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func MaxInt(a int, b int) int {
	return int(math.Max(float64(a), float64(b)))
}
