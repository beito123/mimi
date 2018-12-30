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
	"errors"
	"time"

	"github.com/beito123/mimi"
	"github.com/beito123/mimi/pks"
)

type ServerSession struct {
	mimi.BaseSession

	console *Console
	tracker *LogTracker
}

func (session *ServerSession) Update(handlers []mimi.PacketHandler) {
	session.Update(handlers)

	if session.State() != mimi.StateConnected {
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

func (sp *ServerSessionHandler) HandlePacket(session mimi.Session, pk pks.Packet) {
	switch npk := pk.(type) {
	case *pks.ConnectionRequest:
		logger.Debugf("Received a ConnectionRequest packet from %s\n", session.Addr().String())

		if session.State() != mimi.StateConnecting {
			logger.Debugf("Received a ConnectionRequest packet, but already connected\n")

			// TODO: implements to reconnect

			return
		}

		if npk.ClientProtocol != mimi.ProtocolVersion && !sp.IngoreProtocol {
			session.SendPacket(&pks.IncompatibleProtocol{
				Protocol: mimi.ProtocolVersion,
			})

			session.Close()

			return
		}

		session.SetClientUUID(npk.ClientUUID)

		session.SetState(mimi.StateConnected)

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
				Error: mimi.ErrIDProgramNotFound,
			})

			return
		}

		con, err := sp.ConsoleManager.NewConsole(program.Loader)
		if err != nil {
			logger.Errorln(err)

			session.SendPacket(&pks.ErrorMessage{
				Error: mimi.ErrIDInternalError,
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
				Error: mimi.ErrIDConsoleNotFound,
			})

			return
		}

		serSession, ok := session.(*ServerSession)
		if !ok {
			mimi.Error("couldn't convert to *ServerSession")

			session.SendPacket(&pks.ErrorMessage{
				Error: mimi.ErrIDInternalError,
			})

			return
		}

		if con.Closed() {
			session.SendPacket(&pks.ErrorMessage{
				Error: mimi.ErrIDConsoleAlreadyClosed,
			})

			return
		}

		err := serSession.JoinConsole(con)
		if err != nil {
			mimi.Error("couldn't join a console error: %s", err.Error())

			session.SendPacket(&pks.ErrorMessage{
				Error: mimi.ErrIDInternalError,
			})
		}
	case *pks.QuitConsole:
		logger.Debugf("Received a QuitConsole packet\n")

		serSession, ok := session.(*ServerSession)
		if !ok {
			mimi.Error("couldn't convert to *ServerSession")

			session.SendPacket(&pks.ErrorMessage{
				Error: mimi.ErrIDInternalError,
			})

			return
		}

		if !serSession.HasJoined() {
			session.SendPacket(&pks.ErrorMessage{
				Error: mimi.ErrIDSessionNotJoinedConsole,
			})

			return
		}

		err := serSession.QuitConsole()
		if err != nil {
			mimi.Error("couldn't join a console error: %s", err.Error())

			session.SendPacket(&pks.ErrorMessage{
				Error: mimi.ErrIDInternalError,
			})
		}
	case *pks.DisconnectionNotification:
		logger.Debugf("Received disconnection packet IP: %s CID: %s\n", session.Addr().String(), session.ClientUUID().String())

		if session.State() != mimi.StateDisconnected {
			session.Close()
		}
	default:
		logger.Debugf("Received unknown packet ID:%d\n", npk.ID())
	}
}
