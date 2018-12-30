package server

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import (
	"context"
	"errors"
	"time"

	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map"
	uuid "github.com/satori/go.uuid"
	"github.com/beito123/mimi/pks"

	"github.com/beito123/mimi"
)

type SessionManager struct {
	Handlers []mimi.PacketHandler

	sessions cmap.ConcurrentMap // map[uuid.UUID]Session
}

func (sm *SessionManager) getSession(uid uuid.UUID) (mimi.Session, bool) {
	val, ok := sm.sessions.Get(uid.String())
	if !ok {
		return nil, false
	}

	session, ok := val.(mimi.Session)
	if !ok {
		panic("couldn't convert to Session")
	}

	return session, true
}

func (sm *SessionManager) setSession(session mimi.Session) {
	sm.sessions.Set(session.UUID().String(), session)
}

func (sm *SessionManager) rangeSessions(f func(session mimi.Session) bool) {
	for item := range sm.sessions.IterBuffered() {
		session, ok := item.Val.(mimi.Session)
		if !ok {
			panic("couldn't convert to Session")
		}

		if !f(session) {
			break
		}
	}
}

func (sm *SessionManager) Start(ctx context.Context) {
	ticker := time.NewTicker(mimi.UpdateInterval)
	for _ = range ticker.C {
		select {
		case <-ctx.Done():
			sm.closeSessions()

			ticker.Stop()
			return
		default:
		}

		sm.rangeSessions(func(session mimi.Session) bool {
			session.Update(sm.Handlers)

			return true
		})
	}
}

func (sm *SessionManager) closeSessions() {
	sm.rangeSessions(func(session mimi.Session) bool {
		session.Close()

		return true
	})
}

func (sm *SessionManager) NewSession(conn *websocket.Conn) error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	session := &ServerSession{
		BaseSession: mimi.NewBaseSession(conn, uid, uuid.Nil),
	}

	sm.setSession(session)

	session.Start()

	return nil
}

func (sm *SessionManager) SendPacket(uid uuid.UUID, pk pks.Packet) error {
	session, ok := sm.getSession(uid)
	if !ok {
		return errors.New("couldn't find a session")
	}

	return session.SendPacket(pk)
}

func (sm *SessionManager) SendPacketAll(pk pks.Packet) (err error) {
	sm.rangeSessions(func(session mimi.Session) bool {
		err = session.SendPacket(pk)
		if err != nil {
			return false
		}

		return true
	})

	return err
}
