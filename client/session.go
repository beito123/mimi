package client

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import (
	"gitlab.com/beito123/mimi"
	"gitlab.com/beito123/mimi/pks"
)

type ClientSessionHandler struct {
	//
}

func (sp *ClientSessionHandler) HandlePacket(session mimi.Session, pk pks.Packet) {
	switch npk := pk.(type) {
	case *pks.ConnectionOne:
		logger.Debugf("Received Connection One packet")

		if session.State() != mimi.StateConnecting {
			logger.Debugf("Received a ConnectionOne packet, but already connected\n")
			return
		}

		session.SetUUID(npk.UUID)

		session.SendPacket(&pks.ConnectionRequest{
			ClientProtocol: mimi.ProtocolVersion,
			ClientUUID:     session.ClientUUID(),
		})

		logger.Debugf("Send Connection Request packet")
	case *pks.ConnectionResponse:
		logger.Debugf("Received Connection Response packet")

		if session.State() != mimi.StateConnecting {
			//

			return
		}

		session.SetState(mimi.StateConnected)

		logger.Debugf("Established a connection for a server")
	case *pks.DisconnectionNotification:
		logger.Debugf("Received disconnection packet IP: %s CID: %s", session.Addr().String(), session.ClientUUID().String())

		if session.State() != mimi.StateDisconnected {
			session.Close()
		}
	default:
		logger.Debugf("Received unknown packet ID:%d", npk.ID())
	}
}
