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
}

type SessionManager struct {
	Handlers []*PacketHandler

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
	for {
		select {
		case <-ctx.Done():
			sm.closeSessions()
			break
		default:
		}

		sm.rangeSessions(func(session *Session) bool {
			session.Update()

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

	receivedData chan []byte
	sendData     chan []byte
	errCh          chan error
	closed         chan bool
}

func (session *Session) Addr() net.Addr {
	return session.Conn.RemoteAddr()
}

func (session *Session) Start() {

	session.receivedData = make(chan []byte, 20)
	session.sendData = make(chan []byte, 20)
	session.errCh = make(chan error, 20)

	session.closed = make(chan bool)

	// Receive data
	go func() {
		for {
			select {
			case <-session.closed:
				return
			default:
			}

			typ, data, err := session.Conn.ReadMessage()
			if err != nil {
				session.errCh <- err
				session.Close()
				break
			}

			switch typ {
			case websocket.BinaryMessage:
				session.receivedData <- data
			case websocket.CloseMessage:
				session.Close()
			default:
				session.errCh <- errors.New("unknown message type")
			}
		}
	}()

	// Send data
	go func() {
		for {
			select {
			case <-session.closed:
				return
			default:
			}

		}
	}()
}

func (session *Session) Close() {
	// send close packet

	session.closed <- true

	session.Conn.Close()
}

func (session *Session) Update() {

}

func (session *Session) SendData(data []byte) error {
	//
	return nil
}
