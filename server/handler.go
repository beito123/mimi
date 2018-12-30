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
	"github.com/gorilla/websocket"
	"net/http"

	"github.com/beito123/mimi"
	"github.com/beito123/mimi/util"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: mimi.HandshakeTimeout,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		mimi.Error("happened errors while it's connecting error: %s", reason)
	},
}

type AuthHandler struct {
	Token   string
	Limiter *Limiter
}

func (hand *AuthHandler) Auth(rw mimi.Render, req *http.Request) (bool, error) {
	ip, err := util.IP(req.RemoteAddr)
	if err != nil {
		return false, err
	}

	ok, err := hand.Limiter.Check(ip)
	if err != nil {
		return false, err
	}

	if !ok {
		rw.Write(&mimi.Base{
			Status: http.StatusForbidden,
			Error:  "You are blocked",
		})
	}

	token := req.URL.Query().Get("token")
	if len(token) < 0 {
		return false, nil
	}

	if hand.Token != token {
		rw.Write(&mimi.Base{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})

		return false, nil
	}

	return true, nil
}

type StreamHandler struct {
	Auth    *AuthHandler
	Manager *SessionManager
}

func (hand *StreamHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	render := mimi.NewJSONRender(rw)

	ok, err := hand.Auth.Auth(render, req)
	if err != nil {
		mimi.Dump(err)
		return
	}

	if !ok {
		return
	}

	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		mimi.Dump(err)
		return
	}

	err = hand.Manager.NewSession(conn)
	if err != nil {
		mimi.Dump(err)
	}
}
