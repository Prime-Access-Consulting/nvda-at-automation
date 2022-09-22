package main

import (
	"Server/client"
	"Server/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func loadEnv() {
	log.Println("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}
}

func main() {
	loadEnv()

	nvda, err := client.New(os.Getenv("NVDA_ADDON_HOST"))

	if err != nil {
		log.Fatal(err)
	}

	c := nvda.Capabilities

	log.Printf("Connected to NVDA client: %s v%s on %s", c.Name, c.Version, c.Platform)

	_, err = server.New(nvda, os.Getenv("WEBSOCKET_HOST"))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server on %s Ready.\n", os.Getenv("WEBSOCKET_HOST"))
}
