package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type NVDAClient struct {
	host string
	http *http.Client
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

func (c NVDAClient) New(host string) *NVDAClient {
	client := new(NVDAClient)

	client.host = host

	client.http = &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 5,
	}

	info, err := client.getInfo()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to NVDA: %s\n", *info)

	return client
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
