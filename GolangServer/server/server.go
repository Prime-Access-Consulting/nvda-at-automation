package server

import (
	"AT/GolangServer/AT"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

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

		response := s.HandleMessage(message)

		err = c.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			log.Println("send:", err)
		}
	}
}

func (s *Server) HandleMessage(message []byte) []byte {
	var c = AT.AnyCommand{}
	err := json.Unmarshal(message, &c)

	if err == nil {
		return s.HandleAnyCommand(c, message)
	}

	return s.HandleUnknownCommand()
}

func (s *Server) HandleAnyCommand(command AT.AnyCommand, message []byte) []byte {
	log.Printf("Received command %s\n", command.Method)

	if command.Method == "session.new" {
		var c = AT.SessionNewCommand{}
		err := json.Unmarshal(message, &c)
		if err == nil {
			return s.HandleNewSessionCommand(c)
		}
	}

	if command.Method == "nvda:settings.getSettings" {
		var c = AT.GetSettingsCommand{}
		err := json.Unmarshal(message, &c)
		if err == nil {
			return s.HandleGetSettingsCommand()
		}
	}

	return s.HandleUnknownCommand()
}

func (s *Server) HandleGetSettingsCommand() []byte {
	var c = AT.GetSettingsResponse{}

	for _, v := range *s.clients {
		settings, err := v.GetSettings()
		if err == nil {
			c.Settings = *settings

			response, err := json.Marshal(c)

			if err != nil {
				panic(err)
			}

			return response
		} else {
			panic(err)
		}
	}

	return s.HandleUnknownCommand()
}

func (s *Server) HandleUnknownCommand() []byte {
	var c = AT.ErrorResponse{
		ID:      nil,
		Error:   "unknown command",
		Message: "",
	}

	response, err := json.Marshal(c)

	if err != nil {
		panic(err)
	}

	return response
}

func (s *Server) HandleNewSessionCommand(command AT.SessionNewCommand) []byte {
	//var id = uuid.New()
	//*s.sessions{id.String()} = s.clients[]
	return []byte("started new session")
}

func (s *Server) Start(clients *AT.Clients) {
	log.Println("Starting websocket server")

	s.clients = clients
	s.sessions = new(AT.Sessions)

	http.HandleFunc("/command", s.serve)
	err := http.ListenAndServe(os.Getenv("WEBSOCKET_HOST"), nil)

	if err != nil {
		panic(err)
	}

	log.Printf("Started websocket on %s", os.Getenv("WEBSOCKET_HOST"))
}

type Server struct {
	clients    *AT.Clients
	connection *websocket.Conn
	sessions   *AT.Sessions
}
