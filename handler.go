package mimi

import (
	"net/http"
)

type AuthHandler struct {
	Token   string
	Limiter *Limiter
}

func (hand *AuthHandler) Auth(rw Render, req *http.Request) (bool, error) {
	ip, err := IP(req.RemoteAddr)
	if err != nil {
		return false, err
	}

	ok, err := hand.Limiter.Check(ip)
	if err != nil {
		return false, err
	}

	if !ok {
		rw.Write(&Base{
			Status: http.StatusForbidden,
			Error:  "You are blocked",
		})
	}

	token := req.URL.Query().Get("token")
	if len(token) < 0 {
		return false, nil
	}

	if hand.Token != token {
		rw.Write(&Base{
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
	render := NewJSONRender(rw)

	ok, err := hand.Auth.Auth(render, req)
	if err != nil {
		Dump(err)
		return
	}

	if !ok {
		return
	}

	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		Dump(err)
		return
	}

	err = hand.Manager.NewSession(conn)
	if err != nil {
		Dump(err)
	}
}
