package AT

import (
	"log"
	"os"
)

func LoadATs() *Clients {
	log.Println("Loading available ATs;")
	clients := make(Clients)

	nvda := NVDAClient{}.Client(os.Getenv("NVDA_ADDON_HOST"))

	capabilities, err := nvda.GetInfo()

	if err == nil {
		clients[capabilities] = nvda
		log.Printf("Connected to AT: %s v%s on %s", capabilities.Name, capabilities.Version, capabilities.Platform)
	} else {
		log.Println(err.Error())
	}

	return &clients
}

type Client interface {
	GetInfo() (*Capabilities, error)
	GetSettings() (*Settings, error)
}

type Clients map[*Capabilities]Client

type Capabilities struct {
	Name     string `json:"atName"`
	Version  string `json:"atVersion"`
	Platform string `json:"platformName"`
}

type Settings map[string]interface{}
