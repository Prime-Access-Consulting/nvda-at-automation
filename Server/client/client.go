package client

import (
	"Server/command"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Settings map[string]interface{}

type setSettingsPayload map[string]interface{}

type NVDA struct {
	host         string
	speechPort   string
	LastEventID  int
	http         *http.Client
	Capabilities *Capabilities
}

type Capabilities struct {
	Name     string `json:"atName"`
	Version  string `json:"atVersion"`
	Platform string `json:"platformName"`
}

type onLineCallback func(string)

func New(host string, speechPort string) (*NVDA, error) {
	nvda := &NVDA{
		host:        host,
		speechPort:  speechPort,
		LastEventID: 0,
		http: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	capabilities, err := nvda.getInfo()

	if err != nil {
		return nil, err
	}

	nvda.Capabilities = capabilities

	return nvda, nil
}

func (c *NVDA) getInfo() (*Capabilities, error) {
	res, err := c.http.Get(fmt.Sprintf("%s/%s", c.host, "info"))

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

	capabilities := new(Capabilities)

	if err := json.Unmarshal(body, capabilities); err != nil {
		return nil, err
	}

	return capabilities, nil
}

func (c *NVDA) GetSettings(requestedSettings []string) (*Settings, error) {
	qs := strings.Join(requestedSettings, ",")
	res, err := c.http.Get(fmt.Sprintf("%s/%s?q=%s", c.host, "settings", url.QueryEscape(qs)))

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

	settings := new(Settings)

	if err := json.Unmarshal(body, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func (c *NVDA) SetSettings(settings []command.VendorSettingsSetSettingsParameter) error {
	p := make(setSettingsPayload)
	for _, s := range settings {
		p[s.Name] = s.Value
	}

	payload, err := json.Marshal(p)

	if err != nil {
		return err
	}

	_, err = c.http.Post(fmt.Sprintf("%s/%s", c.host, "settings"), "application/json", bytes.NewReader(payload))

	return err
}

func semverMatches(provided string, requested string) bool {
	if requested == provided {
		return true
	}

	return VersionRequestMatches(provided, requested)
}

func (c *NVDA) MatchesCapabilities(capabilities *command.NewSessionCommandCapabilitiesRequest) bool {
	if capabilities == nil {
		return true
	}

	score := 0
	minimum := 0

	if capabilities.AtName != nil {
		minimum += 1
		if c.Capabilities.Name == *capabilities.AtName {
			score += 1
		}
	}

	if capabilities.AtVersion != nil {
		minimum += 1
		if semverMatches(c.Capabilities.Version, *capabilities.AtVersion) {
			score += 1
		}
	}

	if capabilities.PlatformName != nil {
		minimum += 1
		if c.Capabilities.Platform == *capabilities.PlatformName {
			score += 1
		}
	}

	return score == minimum
}

func (c *NVDA) RegisterOnLineCallback(callback onLineCallback) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%s", c.speechPort))

	if err != nil {
		fmt.Println("dial error:", err)
		return
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	_, writeErr := conn.Write([]byte(fmt.Sprintf("Last-Event-ID:%d", c.LastEventID)))
	if writeErr != nil {
		fmt.Println("ERROR", writeErr)
		return
	}

	r := bufio.NewReader(conn)

	fmt.Println("Reading from event socket")

	for {
		line, err := r.ReadString('\n')

		if err != nil {
			switch err {
			case io.EOF:
				time.Sleep(100 * time.Millisecond)
				continue
			default:
				fmt.Println("ERROR", err)
				continue
			}
		}

		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		callback(line)
	}
}
