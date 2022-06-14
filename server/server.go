package server

import (
	"AT/client"
	"fmt"
	"net/http"
	"nhooyr.io/websocket"
)

type AutomationServer struct {
	socket    *websocket.Conn
	server    *http.Server
	nvda      *client.NVDAClient
	isStarted bool
}

func (s AutomationServer) New(port int) *AutomationServer {
	as := new(AutomationServer)

	as.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      nil,
		ReadTimeout:  60000,
		WriteTimeout: 60000,
	}

	return as
}

func (s *AutomationServer) StartSession(nvda *client.NVDAClient) error {
	if s.isStarted {
		return fmt.Errorf("server has already started")
	}

	err := s.server.ListenAndServe()

	if err != nil {
		return err
	}

	s.isStarted = true
	s.nvda = nvda

	return nil
}
