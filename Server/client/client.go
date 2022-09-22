package client

import (
	"Server/command"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Settings map[string]interface{}

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
