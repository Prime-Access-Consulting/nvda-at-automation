package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type NVDAClient struct {
	host    string
	handler AsyncEventHandler
	http    *http.Client
}

type Command struct {
	Name string `json:"name"`
}

type StartSessionCommandResponse struct {
	Success  bool                 `json:"success"`
	Command  Command              `json:"command"`
	Response StartSessionResponse `json:"response"`
}

type StartSessionResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type CommandErrorResponse struct {
	Success bool    `json:"success"`
	Command Command `json:"command"`
	Error   string  `json:"error"`
}

type AsyncEventHandler func(e string)

func (c NVDAClient) New(host string) *NVDAClient {
	client := new(NVDAClient)
	client.host = host
	client.http = &http.Client{
		Timeout: time.Second * 5,
	}

	info, err := client.getInfo()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to NVDA: %s\n", *info)

	//client.stream()

	return client
}

type event struct {
	Event string `json:"event"`
}

type events []event

func (c *NVDAClient) Stream() {
	servAddr := "localhost:5432"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		log.Print("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	log.Printf("Connected to %s", servAddr)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Print("Dial failed:", err.Error())
		os.Exit(1)
	}

	for {
		reply := make([]byte, 128*1000)
		var chunks []string
		var chunk []byte

		_, err := conn.Read(reply)
		if err != nil {
			//println("Read from server failed:", err.Error())
		}

		for _, r := range reply {
			if int(r) == 30 {
				chunks = append(chunks, string(chunk))
				chunk = nil
				continue
			}

			chunk = append(chunk, r)
		}

		for _, r := range chunks {
			el, _ := json.Marshal(r)
			c.handler(string(el))
		}
	}
}

func (c *NVDAClient) SetAsyncEventHandler(handler AsyncEventHandler) {
	c.handler = handler
}

func (c *NVDAClient) getInfo() (*string, error) {
	res, err := c.http.Get(c.host)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	info := string(body)

	return &info, nil
}

func (c *NVDAClient) StartSession() (*string, error) {
	command := Command{"startSession"}
	payload, err := json.Marshal(command)

	if err != nil {
		return nil, err
	}

	res, err := c.http.Post(fmt.Sprintf("%s/command", c.host), "application/json", bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	s := new(StartSessionCommandResponse)
	e := new(CommandErrorResponse)

	decoder := json.NewDecoder(res.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(s)

	if err == nil {
		return &s.Response.Id, nil
	}

	err = decoder.Decode(e)

	if err == nil {
		return nil, fmt.Errorf(e.Error)
	}

	return nil, fmt.Errorf("unexpected error")
}
