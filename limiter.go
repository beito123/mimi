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
	"net"
	"time"
)

type Addr struct {
	IP net.IP

	Counter  int
	LastTime time.Time

	Blocked     bool
	BlockedTime time.Time
}

type Limiter struct {
	Addrs       []*Addr
	BlockExpire time.Duration

	MaxCount int

	lastUpdateTime time.Time
}

func (hand *Limiter) Update() {
	now := time.Now()
	if now.Sub(hand.lastUpdateTime) > time.Second*1 {
		new := make([]*Addr, 0)
		for _, addr := range hand.Addrs {
			if addr.Counter > hand.MaxCount {
				hand.SetBlock(addr.IP, true)
			}

			addr.Counter = 0

			if addr.Blocked {
				if now.Sub(addr.BlockedTime) < hand.BlockExpire {
					new = append(new, addr)
				}
			}
		}

		hand.Addrs = new
		hand.lastUpdateTime = now
	}
}

func (hand *Limiter) AddAddr(ip net.IP) {
	hand.Addrs = append(hand.Addrs, &Addr{
		IP: ip,
	})
}

func (hand *Limiter) getAddr(ip net.IP) *Addr {
	for _, addr := range hand.Addrs {
		if addr.IP.Equal(ip) {
			return addr
		}
	}

	return nil
}

func (hand *Limiter) HasAddr(ip net.IP) bool {
	return hand.getAddr(ip) != nil
}

func (hand *Limiter) IsBlocked(ip net.IP) bool {
	addr := hand.getAddr(ip)
	return addr.Blocked
}

func (hand *Limiter) SetBlock(ip net.IP, b bool) {
	addr := hand.getAddr(ip)
	if addr != nil {
		addr.Blocked = b
		if b {
			addr.BlockedTime = time.Now()
		}
	}
}

func (hand *Limiter) Check(ip net.IP) (bool, error) {
	hand.Update()

	if !hand.HasAddr(ip) {
		hand.AddAddr(ip)
	}

	addr := hand.getAddr(ip)

	if addr.Blocked {
		return false, nil
	}

	addr.Counter++

	return true, nil
}
