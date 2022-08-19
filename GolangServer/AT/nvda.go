package AT

import (
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

func (c NVDAClient) Client(host string) *NVDAClient {
	client := new(NVDAClient)
	client.host = host
	client.http = &http.Client{
		Timeout: time.Second * 5,
	}

	return client
}

func (c *NVDAClient) GetSettings() (*Settings, error) {
	// TODO accept requested []string as a parameter

	res, err := c.http.Get(fmt.Sprintf("%s/%s", c.host, "settings"))

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

func (c *NVDAClient) GetInfo() (*Capabilities, error) {
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
