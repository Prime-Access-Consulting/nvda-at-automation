package server

import (
	"Server/client"
	"Server/command"
	"Server/response"
	"Server/session"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct {
	connection *websocket.Conn
	sessionID  *string
	client     *client.NVDA
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	c, err := upgrader.Upgrade(w, r, nil)

	s.connection = c

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer func(c *websocket.Conn) {
		var err = c.Close()
		if err != nil {
			panic(err)
		}
	}(c)

	for {
		mt, message, err := c.ReadMessage()

		if mt != websocket.TextMessage {
			log.Println("recv: [Binary], skipping")
			break
		}

		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)

		res := s.handleMessage(message)

		err = c.WriteMessage(websocket.TextMessage, res)
		if err != nil {
			log.Println("send:", err)
		}
	}
}

func (s *Server) handleMessage(message []byte) []byte {
	var c = command.AnyCommand{}
	err := json.Unmarshal(message, &c)

	if err != nil {
		log.Println(err)
		// Unable to parse this as a command
		return handleUnknownCommand(nil)
	}

	return s.handleAnyCommand(c, message)
}

func (s *Server) handleNewSessionCommand(message []byte) []byte {
	var c = command.NewSessionCommand{}
	err := json.Unmarshal(message, &c)

	if err != nil {
		return handleUnknownCommand(nil)
	}

	var e = func(message string) []byte {
		return response.ErrorResponseJSON("session not created", message, &c.ID)
	}

	if !s.client.MatchesCapabilities(c.Params.Capabilities.AlwaysMatch) {
		return e("requested capabilities could not be matched")
	}

	if s.sessionID != nil {
		return e("session already exists")
	}

	s.sessionID = session.NewSessionID()
	info := s.client.Capabilities

	log.Printf("started new session with %s %s on %s", info.Name, info.Version, info.Platform)

	return response.NewSessionResponseJSON(info, *s.sessionID)
}

func (s *Server) handleAnyCommand(c command.AnyCommand, message []byte) []byte {
	log.Printf("Received command %s\n", c.Method)

	if c.Method == command.NewSessionCommandMethod {
		return s.handleNewSessionCommand(message)
	}

	return handleUnknownCommand(&c.ID)
}

func handleUnknownCommand(ID *string) []byte {
	return response.ErrorResponseJSON("unknown command", "", ID)
}

func (s *Server) startWebsocketServer(host string) {
	log.Printf("Starting websocket server on %s", host)

	http.HandleFunc("/command", s.serve)
	err := http.ListenAndServe(host, nil)

	if err != nil {
		panic(err)
	}

	log.Printf("Started websocket on %s", host)

}

func New(client *client.NVDA, host string) (*Server, error) {
	server := new(Server)
	server.client = client
	server.startWebsocketServer(host)
	return server, nil
}
