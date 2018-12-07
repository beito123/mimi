package mimi

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/gorilla/websocket"
	"github.com/orcaman/concurrent-map"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,
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
	HandlePacket(Packet)
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

func (sm *SessionManager) SendPacket(uid uuid.UUID, pk Packet) error {
	session, ok := sm.getSession(uid)
	if !ok {
		return errors.New("couldn't find a session")
	}

	data, err := EncodePacket(pk)
	if err != nil {
		return err
	}

	err = session.SendData(data)
	if err != nil {
		return err
	}

	return nil
}

func (sm *SessionManager) SendPacketAll(pk Packet) error {
	data, err := EncodePacket(pk)
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

type Session struct {
	Conn  *websocket.Conn
	State ConnectionState

	UUID uuid.UUID

	Error *ErrorHandler

	receivedData chan []byte
	sendData     chan []byte
	closeCh      chan bool
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
				session.Error.Handle(err)
				session.Close()
				break
			}

			switch typ {
			case websocket.BinaryMessage:
				session.receivedData <- data
			case websocket.CloseMessage:
				session.Close()
			default:
				session.Error.Handle(errors.New("unknown message type"))
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
				session.Error.Handle(err)
			}
		}
	}()
}

func (session *Session) Close() {
	// send close packet

	close(session.closeCh)

	session.Conn.Close()
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
		pk, err := DecodePacket(data)
		if err != nil {
			session.Error.Handle(err)
		}

		for _, hand := range handlers {
			hand.HandlePacket(pk)
		}
	}
}

func (session *Session) SendData(data []byte) error {
	//
	return nil
}
