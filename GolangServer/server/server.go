package server

import (
	"AT/GolangServer/client"
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"log"
	"net"
	"net/http"
	"nhooyr.io/websocket"
	"os"
	"os/signal"
	"time"
)

type AutomationServer struct {
	socket    *websocket.Conn
	server    *http.Server
	nvda      *client.NVDAClient
	port      int
	isStarted bool
}

func (s *AutomationServer) Start() error {
	if s.isStarted {
		return fmt.Errorf("server has already started")
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	log.Printf("listening on ws://%v", l.Addr())

	errc := make(chan error, 1)
	go func() {
		errc <- s.server.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s AutomationServer) New(port int, nvdaClient *client.NVDAClient) *AutomationServer {
	as := new(AutomationServer)

	as.port = port
	as.nvda = nvdaClient

	as.server = &http.Server{
		Handler:      as,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	return as
}

func (s *AutomationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:       []string{"aria-at-automation"},
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Printf("%v", err)
		return
	}

	defer c.Close(websocket.StatusInternalError, "boink")

	//if c.Subprotocol() != "aria-at-automation" {
	//	c.Close(websocket.StatusPolicyViolation, "client must speak the aria-at-automation subprotocol")
	//	return
	//}

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err = s.handleMessage(r.Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}

		if err != nil {
			log.Fatalf("failed to handle message from %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

func (s *AutomationServer) handleMessage(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	msg, _ := io.ReadAll(r)

	log.Printf("Received: %s.", string(msg))

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	// TODO abstract into switching logic

	if string(msg) == "startSession" {
		id, err := s.nvda.StartSession()

		if err != nil {
			_, err = io.WriteString(w, fmt.Sprintf("Failed to start session: %v", err))
		} else {
			_, err = io.WriteString(w, fmt.Sprintf("Session \"%s\"", *id))
		}

	} else {
		_, err = io.WriteString(w, fmt.Sprintf("Syntax Error: \"%s\"", msg))
	}

	if err != nil {
		return fmt.Errorf("failed to io.WriteString: %w", err)
	}

	return w.Close()
}
