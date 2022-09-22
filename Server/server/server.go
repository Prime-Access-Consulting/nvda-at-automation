package server

import (
	"Server/client"
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

		response := handleMessage(message)

		err = c.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			log.Println("send:", err)
		}
	}
}

func handleMessage(message []byte) []byte {
	return message
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
