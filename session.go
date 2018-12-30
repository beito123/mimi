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

type Session interface {
	Addr() net.Addr
	UUID() uuid.UUID
	SetUUID(uuid.UUID)
	ClientUUID() uuid.UUID
	SetClientUUID(uuid.UUID)
	State() ConnectionState
	SetState(ConnectionState)
	Update([]PacketHandler)
	Close()
	SendPacket(pks.Packet) error
	SendBytes([]byte) error
}

type PacketHandler interface {
	HandlePacket(session Session, pk pks.Packet)
}

type SessionManager struct {
	Handlers []PacketHandler

	sessions cmap.ConcurrentMap // map[uuid.UUID]Session
}

func (sm *SessionManager) getSession(uid uuid.UUID) (Session, bool) {
	val, ok := sm.sessions.Get(uid.String())
	if !ok {
		return nil, false
	}

	session, ok := val.(Session)
	if !ok {
		panic("couldn't convert to Session")
	}

	return session, true
}

func (sm *SessionManager) setSession(session Session) {
	sm.sessions.Set(session.UUID().String(), session)
}

func (sm *SessionManager) rangeSessions(f func(session Session) bool) {
	for item := range sm.sessions.IterBuffered() {
		session, ok := item.Val.(Session)
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

		sm.rangeSessions(func(session Session) bool {
			session.Update(sm.Handlers)

			return true
		})
	}
}

func (sm *SessionManager) closeSessions() {
	sm.rangeSessions(func(session Session) bool {
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
		BaseSession: BaseSession{
			Conn:  conn,
			state: StateConnecting,
			uuid:  uid,
		},
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
	sm.rangeSessions(func(session Session) bool {
		err = session.SendPacket(pk)
		if err != nil {
			return false
		}

		return true
	})

	return err
}

var errClosed = errors.New("already closed")

type BaseSession struct {
	Conn *websocket.Conn

	state      ConnectionState
	uuid       uuid.UUID
	clientUUID uuid.UUID

	receivedData chan []byte
	sendData     chan []byte
	closeCh      chan bool
}

func (session *BaseSession) HandleError(err error) {
	Error("Session IP: %s Error: %s", session.Addr().String(), err.Error())
}

func (session *BaseSession) UUID() uuid.UUID {
	return session.uuid
}

func (session *BaseSession) SetUUID(uid uuid.UUID) {
	session.uuid = uid
}

func (session *BaseSession) ClientUUID() uuid.UUID {
	return session.clientUUID
}

func (session *BaseSession) SetClientUUID(uid uuid.UUID) {
	session.clientUUID = uid
}

func (session *BaseSession) Addr() net.Addr {
	return session.Conn.RemoteAddr()
}

func (session *BaseSession) State() ConnectionState {
	return session.state
}

func (session *BaseSession) SetState(state ConnectionState) {
	session.state = state
}

func (session *BaseSession) Start() {

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

func (session *BaseSession) Close() {
	if session.State() != StateDisconnected {
		return
	}

	session.sendData = make(chan []byte, MaxSendStack)
	session.SendPacket(&pks.DisconnectionNotification{})

	close(session.closeCh)

	session.Conn.Close()

	session.SetState(StateDisconnected)
}

func (session *BaseSession) Update(handlers []PacketHandler) {
	var received [][]byte
	for i := 0; i < MaxProcessData; i++ {
		select {
		case data := <-session.receivedData:
			received = append(received, data)
		default:
			break
		}
	}

	for _, b := range received {
		if len(b) == 0 {
			return
		}

		pk, ok := pks.GetPacket(b[0])
		if !ok {
			logger.Debugf("Received a unknown packet (0x%x)", b[0])

			// TODO: implements to unknown packet
			return
		}

		pk.SetBytes(b)

		err := pk.Decode()
		if err != nil {
			session.HandleError(err)
			return
		}

		for _, hand := range handlers {
			hand.HandlePacket(session, pk)
		}
	}
}

// SendPacket sends a encoded packet to session
func (session *BaseSession) SendPacket(pk pks.Packet) error {
	return session.SendBytes(pk.Bytes())
}

func (session *BaseSession) SendBytes(data []byte) error {
	if session.State() == StateDisconnected {
		return errClosed
	}

	session.sendData <- data // bad hack // TODO: fix blocking

	return nil
}

type ServerSession struct {
	BaseSession

	console *Console
	tracker *LogTracker
}

func (session *ServerSession) Update(handlers []PacketHandler) {
	session.Update(handlers)

	if session.State() != StateConnected {
		return
	}

	if session.console == nil {
		return
	}

	session.SendPacket(&pks.ConsoleMessages{
		Messages: session.console.Lines(session.tracker),
	})
}

func (session *ServerSession) HasJoined() bool {
	return session.console != nil
}

func (session *ServerSession) JoinConsole(con *Console) error {
	if con.Closed() {
		return errors.New("already closed the console")
	}

	session.console = con
	session.tracker = NewLogTracker()

	logger.Debugf("Session(%s) joins a console(%s)", session.Addr().String(), con.UUID.String())

	return nil
}

func (session *ServerSession) QuitConsole() error {
	if !session.HasJoined() {
		return errors.New("not joined a console")
	}

	session.console = nil
	session.tracker = nil

	return nil
}

type ServerSessionHandler struct {
	ProgramManager *ProgramManager
	ConsoleManager *ConsoleManager
	Console        *Console
	IngoreProtocol bool
}

func (sp *ServerSessionHandler) HandlePacket(session Session, pk pks.Packet) {
	switch npk := pk.(type) {
	case *pks.ConnectionRequest:
		logger.Debugf("Received a ConnectionRequest packet from %s\n", session.Addr().String())

		if session.State() != StateConnecting {
			logger.Debugf("Received a ConnectionRequest packet, but already connected\n")

			// TODO: implements to reconnect

			return
		}

		if npk.ClientProtocol != ProtocolVersion && !sp.IngoreProtocol {
			session.SendPacket(&pks.IncompatibleProtocol{
				Protocol: ProtocolVersion,
			})

			session.Close()

			return
		}

		session.SetClientUUID(npk.ClientUUID)

		session.SetState(StateConnected)

		session.SendPacket(&pks.ConnectionResponse{
			Time: time.Now().Unix(),
		})

		logger.Debugf("Established new connection IP: %s CID: %s\n", session.Addr().String(), session.ClientUUID().String())
	case *pks.RequestProgramList:
		logger.Debugf("Received a RequestProgramList packet\n")

		rpk := &pks.ResponseProgramList{}

		for _, p := range sp.ProgramManager.Programs {
			rpk.Programs = append(rpk.Programs, &pks.Program{
				Name:       p.Name,
				LoaderName: p.Loader.Name(),
			})
		}

		session.SendPacket(rpk)
	case *pks.StartProgram:
		logger.Debugf("Received a StartProgram packet\n")

		program, ok := sp.ProgramManager.Get(npk.ProgramName)
		if !ok {
			session.SendPacket(&pks.ErrorMessage{
				Error: ErrIDProgramNotFound,
			})

			return
		}

		con, err := sp.ConsoleManager.NewConsole(program.Loader)
		if err != nil {
			logger.Errorln(err)

			session.SendPacket(&pks.ErrorMessage{
				Error: ErrIDInternalError,
			})

			return
		}

		session.SendPacket(&pks.ProgramStatus{
			ProgramName: program.Name,
			ConsoleUUID: con.UUID,
			Running:     true,
		})
	case *pks.JoinConsole:
		logger.Debugf("Received a JoinConsole packet\n")

		con, ok := sp.ConsoleManager.Get(npk.ConsoleUUID)
		if !ok {
			session.SendPacket(&pks.ErrorMessage{
				Error: ErrIDConsoleNotFound,
			})

			return
		}

		serSession, ok := session.(*ServerSession)
		if !ok {
			Error("couldn't convert to *ServerSession")

			session.SendPacket(&pks.ErrorMessage{
				Error: ErrIDInternalError,
			})

			return
		}

		if con.Closed() {
			session.SendPacket(&pks.ErrorMessage{
				Error: ErrIDConsoleAlreadyClosed,
			})

			return
		}

		err := serSession.JoinConsole(con)
		if err != nil {
			Error("couldn't join a console error: %s", err.Error())

			session.SendPacket(&pks.ErrorMessage{
				Error: ErrIDInternalError,
			})
		}
	case *pks.QuitConsole:
		logger.Debugf("Received a QuitConsole packet\n")

		serSession, ok := session.(*ServerSession)
		if !ok {
			Error("couldn't convert to *ServerSession")

			session.SendPacket(&pks.ErrorMessage{
				Error: ErrIDInternalError,
			})

			return
		}

		if !serSession.HasJoined() {
			session.SendPacket(&pks.ErrorMessage{
				Error: ErrIDSessionNotJoinedConsole,
			})

			return
		}

		err := serSession.QuitConsole()
		if err != nil {
			Error("couldn't join a console error: %s", err.Error())

			session.SendPacket(&pks.ErrorMessage{
				Error: ErrIDInternalError,
			})
		}
	case *pks.DisconnectionNotification:
		logger.Debugf("Received disconnection packet IP: %s CID: %s\n", session.Addr().String(), session.ClientUUID().String())

		if session.State() != StateDisconnected {
			session.Close()
		}
	default:
		logger.Debugf("Received unknown packet ID:%d\n", npk.ID())
	}
}

type ClientSessionHandler struct {
	IngoreProtocol bool
}

func (sp *ClientSessionHandler) HandlePacket(session Session, pk pks.Packet) {
	switch npk := pk.(type) {
	case *pks.ConnectionOne:
		logger.Debugf("Received Connection One packet")

		if session.State() != StateConnecting {
			logger.Debugf("Received a ConnectionOne packet, but already connected\n")
			return
		}

		session.SetUUID(npk.UUID)

		session.SendPacket(&pks.ConnectionRequest{
			ClientProtocol: ProtocolVersion,
			ClientUUID:     session.ClientUUID(),
		})

		logger.Debugf("Send Connection Request packet")
	case *pks.ConnectionResponse:
		logger.Debugf("Received Connection Response packet")

		if session.State() != StateConnecting {
			//

			return
		}

		session.SetState(StateConnected)

		logger.Debugf("Established a connection for a server")
	case *pks.DisconnectionNotification:
		logger.Debugf("Received disconnection packet IP: %s CID: %s", session.Addr().String(), session.ClientUUID().String())

		if session.State() != StateDisconnected {
			session.Close()
		}
	default:
		logger.Debugf("Received unknown packet ID:%d", npk.ID())
	}
}
