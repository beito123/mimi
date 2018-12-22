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
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"gitlab.com/beito123/mimi/pks"

	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: HandshakeTimeout,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		Dump(reason) // TODO: change
	},
}

type ConnectionState int

const (
	StateConnecting ConnectionState = iota
	StateConnected
	StateDisconnected
)

type PacketHandler interface {
	HandlePacket(*Session, pks.Packet)
}

type SessionManager struct {
	Handlers []PacketHandler
	Error    *ErrorHandler

	sessions cmap.ConcurrentMap // map[uuid.UUID]*Session
}

func (sm *SessionManager) getSession(uid uuid.UUID) (*Session, bool) {
	val, ok := sm.sessions.Get(uid.String())
	if !ok {
		return nil, false
	}

	session, ok := val.(*Session)
	if !ok {
		panic("couldn't convert to Session")
	}

	return session, true
}

func (sm *SessionManager) setSession(session *Session) {
	sm.sessions.Set(session.UUID.String(), session)
}

func (sm *SessionManager) rangeSessions(f func(session *Session) bool) {
	for item := range sm.sessions.IterBuffered() {
		session, ok := item.Val.(*Session)
		if !ok {
			panic("couldn't convert to Session")
		}

		if !f(session) {
			break
		}
	}
}

func (sm *SessionManager) Start(ctx context.Context) {
	ticker := time.NewTicker(UpdateInterval)
	for _ = range ticker.C {
		select {
		case <-ctx.Done():
			sm.closeSessions()

			ticker.Stop()
			return
		default:
		}

		sm.rangeSessions(func(session *Session) bool {
			session.Update(sm.Handlers)

			return true
		})
	}
}

func (sm *SessionManager) closeSessions() {
	sm.rangeSessions(func(session *Session) bool {
		session.Close()

		return true
	})
}

func (sm *SessionManager) NewSession(conn *websocket.Conn) error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	session := &Session{
		Conn:  conn,
		State: StateConnecting,
		UUID:  uid,
		Error: sm.Error,
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

func (sm *SessionManager) SendPacketAll(pk pks.Packet) error {
	data, err := pks.EncodePacket(pk)
	if err != nil {
		return err
	}

	sm.rangeSessions(func(session *Session) bool {
		err = session.SendData(data)
		if err != nil {
			return false
		}

		return true
	})

	return err
}

var errClosed = errors.New("already closed")

type Session struct {
	Conn  *websocket.Conn
	State ConnectionState

	UUID       uuid.UUID
	ClientUUID uuid.UUID

	Error *ErrorHandler

	receivedData chan []byte
	sendData     chan []byte
	closeCh      chan bool
}

func (session *Session) HandleError(err error) {
	session.Error.Handle(errors.New("Session IP:" + session.Addr().String() + " " + err.Error()))
}

func (session *Session) Addr() net.Addr {
	return session.Conn.RemoteAddr()
}

func (session *Session) Start() {

	session.receivedData = make(chan []byte, MaxReceiveStack)
	session.sendData = make(chan []byte, MaxSendStack)

	session.closeCh = make(chan bool)

	// Receive data
	go func() {
		for {
			select {
			case <-session.closeCh:
				return
			default:
			}

			typ, data, err := session.Conn.ReadMessage()
			if err != nil {
				session.HandleError(err)
				session.Close()
				break
			}

			switch typ {
			case websocket.BinaryMessage:
				session.receivedData <- data
			case websocket.CloseMessage:
				session.Close()
			default:
				session.HandleError(errors.New("unknown message type"))
			}
		}
	}()

	// Send data
	go func() {
		for {
			var data []byte

			select {
			case <-session.closeCh:
				return
			case n := <-session.sendData:
				data = n
			}

			err := session.Conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				session.HandleError(err)
			}
		}
	}()
}

func (session *Session) Close() {
	if session.State != StateDisconnected {
		return
	}

	session.sendData = make(chan []byte, MaxSendStack)
	session.SendPacket(&pks.DisconnectionNotification{})

	close(session.closeCh)

	session.Conn.Close()

	session.State = StateDisconnected
}

func (session *Session) Update(handlers []PacketHandler) {
	var received [][]byte
	for i := 0; i < MaxProcessData; i++ {
		select {
		case data := <-session.receivedData:
			received = append(received, data)
		default:
			break
		}
	}

	for _, data := range received {
		pk, err := pks.DecodePacket(data)
		if err != nil {
			session.HandleError(err)
		}

		for _, hand := range handlers {
			hand.HandlePacket(session, pk)
		}
	}
}

func (session *Session) SendPacket(pk pks.Packet) error {
	data, err := pks.EncodePacket(pk)
	if err != nil {
		return err
	}

	return session.SendData(data)
}

func (session *Session) SendData(data []byte) error {
	if session.State == StateDisconnected {
		return errClosed
	}

	session.sendData <- data // bad hack // TODO: fix blocking

	return nil
}

type ServerSessionHandler struct {
	ProgramManager *ProgramManager
	ConsoleManager *ConsoleManager
	IngoreProtocol bool
}

func (sp *ServerSessionHandler) HandlePacket(session Session, pk pks.Packet) {
	switch npk := pk.(type) {
	case *pks.ConnectionRequest:
		if session.State != StateConnecting {
			//

			return
		}

		if npk.ClientProtocol != ProtocolVersion && !sp.IngoreProtocol {
			session.SendPacket(&pks.IncompatibleProtocol{
				Protocol: ProtocolVersion,
			})

			session.Close()

			return
		}

		uid, err := uuid.FromString(npk.ClientUUID)
		if err != nil {
			session.SendPacket(&pks.BadRequest{
				Message: "couldn't parse a uuid",
			})

			session.HandleError(err)

			session.Close()

			return
		}

		session.ClientUUID = uid

		session.State = StateConnected

		session.SendPacket(&pks.ConnectionResponse{
			Time: time.Now().Unix(),
		})

		logger.Debugf("Established new connection IP: %s CID: %s", session.Addr().String(), session.ClientUUID.String())
	case *pks.RequestProgramList:
		logger.Debugf("Received a RequestProgramList packet")

		rpk := &pks.ResponseProgramList{}

		for _, p := range sp.ProgramManager.Programs {
			rpk.Programs = append(rpk.Programs, pks.Program {
				Name: p.Name,
				LoaderName:  p.Loader.Name(),
			})
		}

		session.SendPacket(rpk)
	case *pks.DisconnectionNotification:
		logger.Debugf("Received disconnection packet IP: %s CID: %s", session.Addr().String(), session.ClientUUID.String())

		if session.State != StateDisconnected {
			session.Close()
		}
	default:
		logger.Debugf("Received unknown packet ID:%d", npk.ID())
	}
}

type ClientSessionHandler struct {
	IngoreProtocol bool
}

func (sp *ClientSessionHandler) HandlePacket(session Session, pk pks.Packet) {
	switch npk := pk.(type) {
	case *pks.ConnectionOne:
		logger.Debugf("Received Connection One packet")

		if session.State != StateConnecting {
			//

			return
		}

		uid, err := uuid.FromString(npk.UUID)
		if err != nil {
			session.HandleError(err)

			session.Close()
			return
		}

		session.UUID = uid

		session.SendPacket(&pks.ConnectionRequest{
			ClientProtocol: ProtocolVersion,
			ClientUUID:     session.ClientUUID.String(),
		})

		logger.Debugf("Send Connection Request packet")
	case *pks.ConnectionResponse:
		logger.Debugf("Received Connection Response packet")

		if session.State != StateConnecting {
			//

			return
		}

		session.State = StateConnected

		logger.Debugf("Established a connection for a server")
	case *pks.DisconnectionNotification:
		logger.Debugf("Received disconnection packet IP: %s CID: %s", session.Addr().String(), session.ClientUUID.String())

		if session.State != StateDisconnected {
			session.Close()
		}
	default:
		logger.Debugf("Received unknown packet ID:%d", npk.ID())
	}
}

