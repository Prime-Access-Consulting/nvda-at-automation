package main

import (
	atClient "AT/client"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func isDefined(s string) bool {
	return len(s) > 0
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getHost() string {
	executable := os.Getenv("NVDA_EXECUTABLE")

	if !isDefined(executable) {
		log.Fatal("ENV.NVDA_EXECUTABLE is not set to any meaningful value")
	}

	nvdaPort := os.Getenv("NVDA_ADDON_WEBSOCKET_PORT")

	if !isDefined(nvdaPort) {
		log.Fatal("ENV.NVDA_ADDON_WEBSOCKET_PORT is not set to any meaningful value")
	}

	port, err := strconv.Atoi(nvdaPort)

	if err != nil {
		log.Fatal("NVDA addon port is not a valid integer")
	}

	return fmt.Sprintf("http://localhost:%d", port)
}

func main() {
	loadEnv()

	//serverPort, _ := strconv.Atoi(os.Getenv("HOST_WEBSOCKET_PORT"))

	//server := atServer.AutomationServer{}.New(serverPort)

	log.Printf("Connecting...\n")
	nvda := atClient.NVDAClient{}.New(getHost())

	log.Printf("Starting client session...\n")

	sessionId, err := nvda.StartSession()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Established client session %s", *sessionId)

	//err := server.StartSession(nvda)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Printf("Started server on port %d", serverPort)
}
