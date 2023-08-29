package server

import (
	"Server/client"
	"Server/command"
	"Server/event"
	"Server/response"
	"Server/session"
	"Server/sse"
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

func (s *Server) sendMessage(message []byte) error {
	return s.connection.WriteMessage(websocket.TextMessage, message)
}

func (s *Server) handleMessage(message []byte) []byte {
	var c = command.AnyCommand{}
	err := json.Unmarshal(message, &c)

	if err != nil {
		log.Printf("Error unmarshaling: %s", err)
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

	captureOutput := func() {
		parser := sse.NewParser(func(sseEvent sse.Event) {
			message = event.InteractionCapturedOutputEventJSON(s.sessionID, sseEvent.Data)
			err := s.sendMessage(message)
			if err != nil {
				return
			}
			s.client.LastEventID = *sseEvent.ID
		})

		s.client.RegisterOnLineCallback(parser.Process)
	}

	go captureOutput()

	return response.NewSessionResponseJSON(info, *s.sessionID)
}

func getRequestedSettings(command command.GetSettingsCommand) []string {
	var settings []string

	for _, r := range command.Params.Settings {
		settings = append(settings, r.Name)
	}

	return settings
}

func mapSettingsToRetrievedSettings(settings *client.Settings) response.RetrievedSettings {
	s := response.RetrievedSettings{}

	for key, value := range *settings {
		s = append(s, response.RetrievedSetting{Name: key, Value: value})
	}

	return s
}

func (s *Server) getSettings(ID *string, requested []string) []byte {
	res := response.GetSettingsResponse{ID: *ID}

	settings, err := s.client.GetSettings(requested)

	if err != nil {
		log.Fatal(err)
	}

	res.Result = mapSettingsToRetrievedSettings(settings)

	r, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
	}

	return r
}

func (s *Server) handleGetSupportedSettingsCommand(message []byte) []byte {
	c := command.GetSupportedSettingsCommand{}
	err := json.Unmarshal(message, &c)

	if err != nil {
		return handleUnknownCommand(nil)
	}

	if s.client == nil {
		// no session available
		return handleUnknownCommand(&c.ID)
	}

	return s.getSettings(&c.ID, []string{})
}

func (s *Server) handleGetSettingsCommand(message []byte) []byte {
	c := command.GetSettingsCommand{}
	err := json.Unmarshal(message, &c)

	if err != nil {
		return handleUnknownCommand(nil)
	}

	if s.client == nil {
		// no session available
		return handleUnknownCommand(&c.ID)
	}

	if len(c.Params.Settings) == 0 {
		return response.ErrorResponseJSON("invalid argument", "No settings were requested", &c.ID)
	}

	requestedSettings := getRequestedSettings(c)

	return s.getSettings(&c.ID, requestedSettings)
}

func (s *Server) handleSetSettingsCommand(message []byte) []byte {
	c := command.SetSettingsCommand{}
	err := json.Unmarshal(message, &c)

	if err != nil {
		return handleUnknownCommand(nil)
	}

	if s.client == nil {
		// no session available
		return handleUnknownCommand(&c.ID)
	}

	res := response.SetSettingsResponse{ID: c.ID}

	err = s.client.SetSettings(c.Params.Settings)

	if err != nil {
		log.Fatal(err)
	}

	r, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
	}

	return r
}

func (s *Server) handlePressKeysCommand(message []byte) []byte {
	c := command.PressKeysCommand{}

	err := json.Unmarshal(message, &c)

	if err != nil {
		return handleUnknownCommand(nil)
	}

	if s.client == nil {
		// no session available
		return handleUnknownCommand(&c.ID)
	}

	res := response.PressKeysResponse{ID: c.ID}

	err = s.client.PressKeys(c.Params.Keys)

	if err != nil {
		log.Fatal(err)
	}

	r, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
	}

	return r
}

func (s *Server) handleAnyCommand(c command.AnyCommand, message []byte) []byte {
	log.Printf("Received command %s\n", c.Method)

	if c.Method == command.NewSessionCommandMethod {
		return s.handleNewSessionCommand(message)
	}

	if c.Method == command.GetSettingsCommandMethod {
		return s.handleGetSettingsCommand(message)
	}

	if c.Method == command.GetSupportedSettingsCommandMethod {
		return s.handleGetSupportedSettingsCommand(message)
	}

	if c.Method == command.SetSettingsCommandMethod {
		return s.handleSetSettingsCommand(message)
	}

	if c.Method == command.PressKeysCommandMethod {
		return s.handlePressKeysCommand(message)
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
