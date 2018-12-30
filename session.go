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
	"errors"
	"net"

	"gitlab.com/beito123/mimi/pks"

	uuid "github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type ConnectionState int

const (
	StateConnecting ConnectionState = iota
	StateConnected
	StateDisconnected
)

type Session interface {

	// Addr returns client addr
	Addr() net.Addr

	// UUID returns uuid
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

var errClosed = errors.New("already closed")

func NewBaseSession(conn *websocket.Conn, uid uuid.UUID, cid uuid.UUID) BaseSession {
	return BaseSession{
		Conn:       conn,
		state:      StateConnecting,
		uuid:       uid,
		clientUUID: cid,
	}
}

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
	Error("session error (%s): %s", session.Addr().String(), err.Error())
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
