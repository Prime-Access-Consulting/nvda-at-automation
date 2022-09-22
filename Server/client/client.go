package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type NVDA struct {
	host         string
	http         *http.Client
	Capabilities *Capabilities
}

type Capabilities struct {
	Name     string `json:"atName"`
	Version  string `json:"atVersion"`
	Platform string `json:"platformName"`
}

func New(host string) (*NVDA, error) {
	nvda := &NVDA{
		host: host,
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
